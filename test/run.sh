newman run http-monitoring.postman_collection.json \
-d api_test_data.json \
-e http-monitoring.postman_environment.json \
--verbose \
--delay-request 100 \
--timeout 60000