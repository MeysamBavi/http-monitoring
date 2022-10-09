# AlertsApi

All URIs are relative to *http://localhost*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getAlerts**](AlertsApi.md#getAlerts) | **GET** /alerts/{id} | Gets all alerts |


<a name="getAlerts"></a>
# **getAlerts**
> List getAlerts(id)

Gets all alerts

    Gets all alerts

### Parameters

|Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | **String**| url id | [default to null] |

### Return type

[**List**](../Models/ModelAlert.md)

### Authorization

[jwtBearerAuth](../README.md#jwtBearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

