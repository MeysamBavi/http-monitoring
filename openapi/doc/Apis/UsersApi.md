# UsersApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUser**](UsersApi.md#createUser) | **POST** /users | Creates a new user |
| [**loginUser**](UsersApi.md#loginUser) | **POST** /users/login | Authenticates user and generates JWT token |


<a name="createUser"></a>
# **createUser**
> ModelUser createUser(RequestUser)

Creates a new user

    Creates a new user with the given username and password

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **RequestUser** | [**RequestUser**](../Models/RequestUser.md)|  | [optional] |

### Return type

[**ModelUser**](../Models/ModelUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="loginUser"></a>
# **loginUser**
> String loginUser(RequestUser)

Authenticates user and generates JWT token

    Authenticates user and generates JWT token

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **RequestUser** | [**RequestUser**](../Models/RequestUser.md)|  | [optional] |

### Return type

**String**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: JWT token, application/json

