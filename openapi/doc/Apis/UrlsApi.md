# UrlsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUrl**](UrlsApi.md#createUrl) | **POST** /urls | Creates a new url for user |
| [**getAllUrls**](UrlsApi.md#getAllUrls) | **GET** /urls | Returns all urls of user |
| [**getDayStats**](UrlsApi.md#getDayStats) | **GET** /urls/{id}/stats | Returns url monitoring stats |


<a name="createUrl"></a>
# **createUrl**
> ModelURL createUrl(RequestURL)

Creates a new url for user

    Creates a new url for user

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **RequestURL** | [**RequestURL**](../Models/RequestURL.md)|  | [optional] |

### Return type

[**ModelURL**](../Models/ModelURL.md)

### Authorization

[jwtBearerAuth](../README.md#jwtBearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

<a name="getAllUrls"></a>
# **getAllUrls**
> List getAllUrls()

Returns all urls of user

    Returns all urls of user in a list

### Parameters
This endpoint does not need any parameter.

### Return type

[**List**](../Models/ModelURL.md)

### Authorization

[jwtBearerAuth](../README.md#jwtBearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

<a name="getDayStats"></a>
# **getDayStats**
> List getDayStats(id, day, month, year)

Returns url monitoring stats

    Returns monitoring stats for a specific url. Stats can be filtered using query parameters

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| url id | [default to null] |
| **day** | **Integer**| day of the month (1-31) | [optional] [default to null] |
| **month** | **Integer**| month number (1-12) | [optional] [default to null] |
| **year** | **Integer**|  | [optional] [default to null] |

### Return type

[**List**](../Models/ModelDayStat.md)

### Authorization

[jwtBearerAuth](../README.md#jwtBearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

