# {{classname}}

All URIs are relative to *http://localhost/api/*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ChartrepoChartsPost**](ChartRepositoryApi.md#ChartrepoChartsPost) | **Post** /chartrepo/charts | Upload a chart file to the defult &#x27;library&#x27; project.
[**ChartrepoHealthGet**](ChartRepositoryApi.md#ChartrepoHealthGet) | **Get** /chartrepo/health | Check the health of chart repository service.
[**ChartrepoRepoChartsGet**](ChartRepositoryApi.md#ChartrepoRepoChartsGet) | **Get** /chartrepo/{repo}/charts | Get all the charts under the specified project
[**ChartrepoRepoChartsNameDelete**](ChartRepositoryApi.md#ChartrepoRepoChartsNameDelete) | **Delete** /chartrepo/{repo}/charts/{name} | Delete all the versions of the specified chart
[**ChartrepoRepoChartsNameGet**](ChartRepositoryApi.md#ChartrepoRepoChartsNameGet) | **Get** /chartrepo/{repo}/charts/{name} | Get all the versions of the specified chart
[**ChartrepoRepoChartsNameVersionDelete**](ChartRepositoryApi.md#ChartrepoRepoChartsNameVersionDelete) | **Delete** /chartrepo/{repo}/charts/{name}/{version} | Delete the specified chart version
[**ChartrepoRepoChartsNameVersionGet**](ChartRepositoryApi.md#ChartrepoRepoChartsNameVersionGet) | **Get** /chartrepo/{repo}/charts/{name}/{version} | Get the specified chart version
[**ChartrepoRepoChartsPost**](ChartRepositoryApi.md#ChartrepoRepoChartsPost) | **Post** /chartrepo/{repo}/charts | Upload a chart file to the specified project.
[**ChartrepoRepoProvPost**](ChartRepositoryApi.md#ChartrepoRepoProvPost) | **Post** /chartrepo/{repo}/prov | Upload a provance file to the specified project.

# **ChartrepoChartsPost**
> ChartrepoChartsPost(ctx, chart)
Upload a chart file to the defult 'library' project.

Upload a chart file to the default 'library' project. Uploading together with the prov file at the same time is also supported.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **chart** | ***os.File*****os.File**|  | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoHealthGet**
> HarborInlineResponse2001 ChartrepoHealthGet(ctx, )
Check the health of chart repository service.

Check the health of chart repository service.

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**HarborInlineResponse2001**](inline_response_200_1.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsGet**
> []HarborChartInfoEntry ChartrepoRepoChartsGet(ctx, repo)
Get all the charts under the specified project

Get all the charts under the specified project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **repo** | **string**| The project name | 

### Return type

[**[]HarborChartInfoEntry**](ChartInfoEntry.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsNameDelete**
> ChartrepoRepoChartsNameDelete(ctx, repo, name)
Delete all the versions of the specified chart

Delete all the versions of the specified chart

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **repo** | **string**| The project name | 
  **name** | **string**| The chart name | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsNameGet**
> []HarborChartVersion ChartrepoRepoChartsNameGet(ctx, repo, name)
Get all the versions of the specified chart

Get all the versions of the specified chart

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **repo** | **string**| The project name | 
  **name** | **string**| The chart name | 

### Return type

[**[]HarborChartVersion**](array.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsNameVersionDelete**
> ChartrepoRepoChartsNameVersionDelete(ctx, repo, name, version)
Delete the specified chart version

Delete the specified chart version

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **repo** | **string**| The project name | 
  **name** | **string**| The chart name | 
  **version** | **string**| The chart version | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsNameVersionGet**
> HarborChartVersionDetails ChartrepoRepoChartsNameVersionGet(ctx, repo, name, version)
Get the specified chart version

Get the specified chart version

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **repo** | **string**| The project name | 
  **name** | **string**| The chart name | 
  **version** | **string**| The chart version | 

### Return type

[**HarborChartVersionDetails**](ChartVersionDetails.md)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoChartsPost**
> ChartrepoRepoChartsPost(ctx, chart, repo)
Upload a chart file to the specified project.

Upload a chart file to the specified project. With this API, the corresponding provance file can be uploaded together with chart file at once.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **chart** | ***os.File*****os.File**|  | 
  **repo** | **string**| The project name | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ChartrepoRepoProvPost**
> ChartrepoRepoProvPost(ctx, prov, repo)
Upload a provance file to the specified project.

Upload a provance file to the specified project. The provance file should be targeted for an existing chart file.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **prov** | ***os.File*****os.File**|  | 
  **repo** | **string**| The project name | 

### Return type

 (empty response body)

### Authorization

[basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

