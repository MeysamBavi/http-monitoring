echo off
newman run http-monitoring.postman_collection.json -d api_test_data.json --verbose --delay-request 100 --timeout 20000