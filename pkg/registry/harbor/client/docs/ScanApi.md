# {{classname}}

All URIs are relative to *http://localhost/api/v2.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetReportLog**](ScanApi.md#GetReportLog) | **Get** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/scan/{report_id}/log | Get the log of the scan report
[**ScanArtifact**](ScanApi.md#ScanArtifact) | **Post** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/scan | Scan the artifact

# **GetReportLog**
> string GetReportLog(ctx, projectName, repositoryName, reference, reportId, optional)
Get the log of the scan report

Get the log of the scan report

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
  **reportId** | **string**| The report id to get the log | 
 **optional** | ***ScanApiGetReportLogOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ScanApiGetReportLogOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

**string**

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ScanArtifact**
> ScanArtifact(ctx, projectName, repositoryName, reference, optional)
Scan the artifact

Scan the specified artifact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ScanApiScanArtifactOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ScanApiScanArtifactOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

 (empty response body)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

