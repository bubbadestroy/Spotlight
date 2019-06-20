# Architecture

There are three components to this project:
- *The Scanning Engine:*  This engine, when executed, will run all the configured scans 
  against all of the configured domains and generate json output.
- *Scanning Engine Plugins:*  These will be plugins that will let you execute different
  types of scans.  They will be used by the scanning engine, or can be executed by hand.
  It is anticipated that over time, more plugins will be created to expand the functionality
  of the scanning engine.
- *Scan Result UI:*  This is a web frontend that can be used to download/search scan results.

The components are deployed using CircleCI into cloud.gov and once a day, the scanning engine is
executed using https://docs.cloudfoundry.org/devguide/using-tasks.html#run-tasks, generating
it's json files into S3.  The UI then can be used to find and download the scan results.
