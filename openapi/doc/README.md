# Documentation for http-monitoring

<a name="documentation-for-api-endpoints"></a>
## Documentation for API Endpoints

All URIs are relative to *http://localhost*

| Class | Method | HTTP request | Description |
|------------ | ------------- | ------------- | -------------|
| *AlertsApi* | [**getAlerts**](Apis/AlertsApi.md#getalerts) | **GET** /alerts/{id} | Gets all alerts |
| *UrlsApi* | [**createUrl**](Apis/UrlsApi.md#createurl) | **POST** /urls | Creates a new url for user |
*UrlsApi* | [**getAllUrls**](Apis/UrlsApi.md#getallurls) | **GET** /urls | Returns all urls of user |
*UrlsApi* | [**getDayStats**](Apis/UrlsApi.md#getdaystats) | **GET** /urls/{id}/stats | Returns url monitoring stats |
| *UsersApi* | [**createUser**](Apis/UsersApi.md#createuser) | **POST** /users | Creates a new user |
*UsersApi* | [**loginUser**](Apis/UsersApi.md#loginuser) | **POST** /users/login | Authenticates user and generates JWT token |


<a name="documentation-for-models"></a>
## Documentation for Models

 - [ModelAlert](./Models/ModelAlert.md)
 - [ModelDate](./Models/ModelDate.md)
 - [ModelDayStat](./Models/ModelDayStat.md)
 - [ModelURL](./Models/ModelURL.md)
 - [ModelUser](./Models/ModelUser.md)
 - [RequestURL](./Models/RequestURL.md)
 - [RequestUser](./Models/RequestUser.md)
 - [V4HTTPError](./Models/V4HTTPError.md)


<a name="documentation-for-authorization"></a>
## Documentation for Authorization

<a name="jwtBearerAuth"></a>
### jwtBearerAuth

- **Type**: HTTP basic authentication

