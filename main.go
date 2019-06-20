package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var bucket string

// "s3": [
//  {
//   "binding_name": null,
//   "credentials": {
//    "access_key_id": "XXX",
//    "additional_buckets": [],
//    "bucket": "bucketname",
//    "region": "us-gov-west-1",
//    "secret_access_key": "YYYYYYYYYYYYYY",
//    "uri": "s3://XXX:YYYYYYYYYYYYYY@s3-us-gov-west-1.amazonaws.com/bucketname"
//   },
//   "instance_name": "storage",
//   "label": "s3",
//   "name": "storage",
//   "plan": "basic-sandbox",
//   "provider": null,
//   "syslog_drain_url": null,
//   "tags": [
//    "AWS",
//    "S3",
//    "object-storage"
//   ],
//   "volume_mounts": []
//  }
// ],

func main() {
	http.HandleFunc("/runscan", runScanHandler)
	http.HandleFunc("/scans", scansHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Set up AWS credentials
	servicejson := os.Getenv("VCAP_SERVICES")
	var serviceinfo map[string]interface{}
	err := json.Unmarshal([]byte(servicejson), &serviceinfo)
	if err != nil {
		log.Printf("could not parse VCAP_SERVICES")
		return
	}
	s3info := serviceinfo["s3"].([]interface{})
	firsts3info := s3info[0].(map[string]interface{})
	credinfo := firsts3info["credentials"].(map[string]interface{})
	bucket = credinfo["bucket"].(string)
	os.Setenv("AWS_ACCESS_KEY_ID", credinfo["access_key_id"].(string))
	os.Setenv("AWS_SECRET_ACCESS_KEY", credinfo["secret_access_key"].(string))
	os.Setenv("AWS_DEFAULT_REGION", credinfo["region"].(string))

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func scansHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize a session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_DEFAULT_REGION")),
	})
	if err != nil {
		log.Printf("Unable to initiate session: %v", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	// Get the list of items
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		log.Printf("Unable to list items: %v", err)
	}

	fmt.Fprintf(w, "<html><body>\n")
	for _, item := range resp.Contents {
		fmt.Fprintln(w, "Name:         ", *item.Key)
		fmt.Fprintln(w, "Last modified:", *item.LastModified)
		fmt.Fprintln(w, "Size:         ", *item.Size)
		fmt.Fprintln(w, "<hr>")
	}

	fmt.Fprintln(w, "Found", len(resp.Contents), "items")
	fmt.Fprintf(w, "</body></html>\n")
}

func runScanHandler(w http.ResponseWriter, r *http.Request) {
	// XXX Make sure that this is a request from something we scheduled
	// XXX use the application_id in VCAP_APPLICATION?
	// auth_info := os.Getenv("VCAP_APPLICATION")
	// if auth_info == "" {
	// 	log.Printf("VCAP_APPLICATION not set in environment")
	// 	http.Error(w, "cannot authenticate client", 500)
	// 	return
	// }
	// if r.Header.Get("Scan-Authorization") != auth_info {
	// 	http.NotFound(w, r)
	// 	return
	// }

	// run gsutil rsync to get the logs where they need to go
	cmd := exec.Command("scan_engine.py", bucket)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("scan engine failure: %s %s\n", err.Error(), output)
		http.Error(w, "scan engine failed", 503)
		return
	}
	log.Printf("scan engine run succeded")
	fmt.Fprint(w, "OK")
}
