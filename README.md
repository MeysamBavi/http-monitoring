# http-monitoring
A simple HTTP monitoring service in Go.

## API Documentation
For reading full API documentation, visit [here](openapi/doc).

## HTTP Framework
This service uses [echo](https://echo.labstack.com/) for handling http requests.

## Database
This service uses MongoDB as database. There is also an in-memory datastore implementation for testing purposes, which can be configured in [config.json](config.json). 