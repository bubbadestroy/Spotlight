from django.test import TestCase

# Create your tests here.
from django.urls import reverse
from rest_framework.test import APITestCase, APIClient
from rest_framework.views import status
from .serializers import DomainsSerializer
from .views import getScansfromS3
from django.conf import settings
import boto3


# tests for views
class CheckBucketTest(APITestCase):
    """
    This test case checks that the bucketname is getting parsed properly
    """
    def test_bucket_parse(self):
        self.assertIsNotNone(settings.BUCKETNAME)

class GetAllScansTest(APITestCase):
    client = APIClient()

    def test_get_all_scans(self):
        """
        This test ensures that all scans
        exist when we make a GET request to the scans/ endpoint
        """
        # hit the API endpoint
        response = self.client.get(
            reverse("domain-list", kwargs={"version": "v1"})
        )
        # fetch the data
        expected = getScansfromS3()

        serialized = DomainsSerializer(expected, many=True)
        self.assertEqual(response.data, serialized.data)
        self.assertEqual(response.status_code, status.HTTP_200_OK)

