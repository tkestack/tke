# {{classname}}

All URIs are relative to *http://localhost/api/v2.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddLabel**](ArtifactApi.md#AddLabel) | **Post** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/labels | Add label to artifact
[**CopyArtifact**](ArtifactApi.md#CopyArtifact) | **Post** /projects/{project_name}/repositories/{repository_name}/artifacts | Copy artifact
[**CreateTag**](ArtifactApi.md#CreateTag) | **Post** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags | Create tag
[**DeleteArtifact**](ArtifactApi.md#DeleteArtifact) | **Delete** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference} | Delete the specific artifact
[**DeleteTag**](ArtifactApi.md#DeleteTag) | **Delete** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags/{tag_name} | Delete tag
[**GetAddition**](ArtifactApi.md#GetAddition) | **Get** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/additions/{addition} | Get the addition of the specific artifact
[**GetArtifact**](ArtifactApi.md#GetArtifact) | **Get** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference} | Get the specific artifact
[**ListArtifacts**](ArtifactApi.md#ListArtifacts) | **Get** /projects/{project_name}/repositories/{repository_name}/artifacts | List artifacts
[**ListTags**](ArtifactApi.md#ListTags) | **Get** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/tags | List tags
[**RemoveLabel**](ArtifactApi.md#RemoveLabel) | **Delete** /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}/labels/{label_id} | Remove label from artifact

# **AddLabel**
> AddLabel(ctx, body, projectName, repositoryName, reference, optional)
Add label to artifact

Add label to the specified artiact.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborLabel**](HarborLabel.md)| The label that added to the artifact. Only the ID property is needed. | 
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ArtifactApiAddLabelOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiAddLabelOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **xRequestId** | **optional.**| An unique ID for the request | 

### Return type

 (empty response body)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CopyArtifact**
> CopyArtifact(ctx, projectName, repositoryName, from, optional)
Copy artifact

Copy the artifact specified in the \"from\" parameter to the repository.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **from** | **string**| The artifact from which the new artifact is copied from, the format should be \&quot;project/repository:tag\&quot; or \&quot;project/repository@digest\&quot;. | 
 **optional** | ***ArtifactApiCopyArtifactOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiCopyArtifactOpts struct
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

# **CreateTag**
> CreateTag(ctx, body, projectName, repositoryName, reference, optional)
Create tag

Create a tag for the specified artifact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborTag**](HarborTag.md)| The JSON object of tag. | 
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ArtifactApiCreateTagOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiCreateTagOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **xRequestId** | **optional.**| An unique ID for the request | 

### Return type

 (empty response body)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteArtifact**
> DeleteArtifact(ctx, projectName, repositoryName, reference, optional)
Delete the specific artifact

Delete the artifact specified by the reference under the project and repository. The reference can be digest or tag

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ArtifactApiDeleteArtifactOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiDeleteArtifactOpts struct
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

# **DeleteTag**
> DeleteTag(ctx, projectName, repositoryName, reference, tagName, optional)
Delete tag

Delete the tag of the specified artifact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
  **tagName** | **string**| The name of the tag | 
 **optional** | ***ArtifactApiDeleteTagOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiDeleteTagOpts struct
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

# **GetAddition**
> string GetAddition(ctx, projectName, repositoryName, reference, addition, optional)
Get the addition of the specific artifact

Get the addition of the artifact specified by the reference under the project and repository.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
  **addition** | **string**| The type of addition. | 
 **optional** | ***ArtifactApiGetAdditionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiGetAdditionOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

**string**

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetArtifact**
> HarborArtifact GetArtifact(ctx, projectName, repositoryName, reference, optional)
Get the specific artifact

Get the artifact specified by the reference under the project and repository. The reference can be digest or tag.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ArtifactApiGetArtifactOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiGetArtifactOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRequestId** | **optional.String**| An unique ID for the request | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **withTag** | **optional.Bool**| Specify whether the tags are inclued inside the returning artifacts | [default to true]
 **withLabel** | **optional.Bool**| Specify whether the labels are inclued inside the returning artifacts | [default to false]
 **withScanOverview** | **optional.Bool**| Specify whether the scan overview is inclued inside the returning artifacts | [default to false]
 **withSignature** | **optional.Bool**| Specify whether the signature is inclued inside the returning artifacts | [default to false]
 **withImmutableStatus** | **optional.Bool**| Specify whether the immutable status is inclued inside the tags of the returning artifacts. Only works when setting \&quot;with_tag&#x3D;true\&quot; | [default to false]

### Return type

[**HarborArtifact**](Artifact.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListArtifacts**
> []HarborArtifact ListArtifacts(ctx, projectName, repositoryName, optional)
List artifacts

List artifacts under the specific project and repository. Except the basic properties, the other supported queries in \"q\" includes \"tags=*\" to list only tagged artifacts, \"tags=nil\" to list only untagged artifacts, \"tags=~v\" to list artifacts whose tag fuzzy matches \"v\", \"tags=v\" to list artifact whose tag exactly matches \"v\", \"labels=(id1, id2)\" to list artifacts that both labels with id1 and id2 are added to

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
 **optional** | ***ArtifactApiListArtifactsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiListArtifactsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xRequestId** | **optional.String**| An unique ID for the request | 
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **withTag** | **optional.Bool**| Specify whether the tags are included inside the returning artifacts | [default to true]
 **withLabel** | **optional.Bool**| Specify whether the labels are included inside the returning artifacts | [default to false]
 **withScanOverview** | **optional.Bool**| Specify whether the scan overview is included inside the returning artifacts | [default to false]
 **withSignature** | **optional.Bool**| Specify whether the signature is included inside the tags of the returning artifacts. Only works when setting \&quot;with_tag&#x3D;true\&quot; | [default to false]
 **withImmutableStatus** | **optional.Bool**| Specify whether the immutable status is included inside the tags of the returning artifacts. Only works when setting \&quot;with_tag&#x3D;true\&quot; | [default to false]

### Return type

[**[]HarborArtifact**](Artifact.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListTags**
> []HarborTag ListTags(ctx, projectName, repositoryName, reference, optional)
List tags

List tags of the specific artifact

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
 **optional** | ***ArtifactApiListTagsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiListTagsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRequestId** | **optional.String**| An unique ID for the request | 
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **withSignature** | **optional.Bool**| Specify whether the signature is included inside the returning tags | [default to false]
 **withImmutableStatus** | **optional.Bool**| Specify whether the immutable status is included inside the returning tags | [default to false]

### Return type

[**[]HarborTag**](Tag.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RemoveLabel**
> RemoveLabel(ctx, projectName, repositoryName, reference, labelId, optional)
Remove label from artifact

Remove the label from the specified artiact.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **repositoryName** | **string**| The name of the repository. If it contains slash, encode it with URL encoding. e.g. a/b -&gt; a%252Fb | 
  **reference** | **string**| The reference of the artifact, can be digest or tag | 
  **labelId** | **int64**| The ID of the label that removed from the artifact. | 
 **optional** | ***ArtifactApiRemoveLabelOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ArtifactApiRemoveLabelOpts struct
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

