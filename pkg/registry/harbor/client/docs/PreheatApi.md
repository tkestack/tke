# {{classname}}

All URIs are relative to *http://localhost/api/v2.0*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateInstance**](PreheatApi.md#CreateInstance) | **Post** /p2p/preheat/instances | Create p2p provider instances
[**CreatePolicy**](PreheatApi.md#CreatePolicy) | **Post** /projects/{project_name}/preheat/policies | Create a preheat policy under a project
[**DeleteInstance**](PreheatApi.md#DeleteInstance) | **Delete** /p2p/preheat/instances/{preheat_instance_name} | Delete the specified P2P provider instance
[**DeletePolicy**](PreheatApi.md#DeletePolicy) | **Delete** /projects/{project_name}/preheat/policies/{preheat_policy_name} | Delete a preheat policy
[**GetExecution**](PreheatApi.md#GetExecution) | **Get** /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions/{execution_id} | Get a execution detail by id
[**GetInstance**](PreheatApi.md#GetInstance) | **Get** /p2p/preheat/instances/{preheat_instance_name} | Get a P2P provider instance
[**GetPolicy**](PreheatApi.md#GetPolicy) | **Get** /projects/{project_name}/preheat/policies/{preheat_policy_name} | Get a preheat policy
[**GetPreheatLog**](PreheatApi.md#GetPreheatLog) | **Get** /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions/{execution_id}/tasks/{task_id}/logs | Get the log text stream of the specified task for the given execution
[**ListExecutions**](PreheatApi.md#ListExecutions) | **Get** /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions | List executions for the given policy
[**ListInstances**](PreheatApi.md#ListInstances) | **Get** /p2p/preheat/instances | List P2P provider instances
[**ListPolicies**](PreheatApi.md#ListPolicies) | **Get** /projects/{project_name}/preheat/policies | List preheat policies
[**ListProviders**](PreheatApi.md#ListProviders) | **Get** /p2p/preheat/providers | List P2P providers
[**ListProvidersUnderProject**](PreheatApi.md#ListProvidersUnderProject) | **Get** /projects/{project_name}/preheat/providers | Get all providers at project level
[**ListTasks**](PreheatApi.md#ListTasks) | **Get** /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions/{execution_id}/tasks | List all the related tasks for the given execution
[**ManualPreheat**](PreheatApi.md#ManualPreheat) | **Post** /projects/{project_name}/preheat/policies/{preheat_policy_name} | Manual preheat
[**PingInstances**](PreheatApi.md#PingInstances) | **Post** /p2p/preheat/instances/ping | Ping status of a instance.
[**StopExecution**](PreheatApi.md#StopExecution) | **Patch** /projects/{project_name}/preheat/policies/{preheat_policy_name}/executions/{execution_id} | Stop a execution
[**UpdateInstance**](PreheatApi.md#UpdateInstance) | **Put** /p2p/preheat/instances/{preheat_instance_name} | Update the specified P2P provider instance
[**UpdatePolicy**](PreheatApi.md#UpdatePolicy) | **Put** /projects/{project_name}/preheat/policies/{preheat_policy_name} | Update preheat policy

# **CreateInstance**
> HarborInstanceCreatedResp CreateInstance(ctx, body, optional)
Create p2p provider instances

Create p2p provider instances

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborInstance**](HarborInstance.md)| The JSON object of instance. | 
 **optional** | ***PreheatApiCreateInstanceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiCreateInstanceOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRequestId** | **optional.**| An unique ID for the request | 

### Return type

[**HarborInstanceCreatedResp**](InstanceCreatedResp.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreatePolicy**
> CreatePolicy(ctx, body, projectName, optional)
Create a preheat policy under a project

Create a preheat policy under a project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborPreheatPolicy**](HarborPreheatPolicy.md)| The policy schema info | 
  **projectName** | **string**| The name of the project | 
 **optional** | ***PreheatApiCreatePolicyOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiCreatePolicyOpts struct
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

# **DeleteInstance**
> HarborInstanceDeletedResp DeleteInstance(ctx, preheatInstanceName, optional)
Delete the specified P2P provider instance

Delete the specified P2P provider instance

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **preheatInstanceName** | **string**| Instance Name | 
 **optional** | ***PreheatApiDeleteInstanceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiDeleteInstanceOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**HarborInstanceDeletedResp**](InstanceDeletedResp.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeletePolicy**
> DeletePolicy(ctx, projectName, preheatPolicyName, optional)
Delete a preheat policy

Delete a preheat policy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
 **optional** | ***PreheatApiDeletePolicyOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiDeletePolicyOpts struct
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

# **GetExecution**
> HarborExecution GetExecution(ctx, projectName, preheatPolicyName, executionId, optional)
Get a execution detail by id

Get a execution detail by id

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
  **executionId** | **int32**| Execution ID | 
 **optional** | ***PreheatApiGetExecutionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiGetExecutionOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**HarborExecution**](Execution.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstance**
> HarborInstance GetInstance(ctx, preheatInstanceName, optional)
Get a P2P provider instance

Get a P2P provider instance

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **preheatInstanceName** | **string**| Instance Name | 
 **optional** | ***PreheatApiGetInstanceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiGetInstanceOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**HarborInstance**](Instance.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPolicy**
> HarborPreheatPolicy GetPolicy(ctx, projectName, preheatPolicyName, optional)
Get a preheat policy

Get a preheat policy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
 **optional** | ***PreheatApiGetPolicyOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiGetPolicyOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**HarborPreheatPolicy**](PreheatPolicy.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetPreheatLog**
> string GetPreheatLog(ctx, projectName, preheatPolicyName, executionId, taskId, optional)
Get the log text stream of the specified task for the given execution

Get the log text stream of the specified task for the given execution

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
  **executionId** | **int32**| Execution ID | 
  **taskId** | **int32**| Task ID | 
 **optional** | ***PreheatApiGetPreheatLogOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiGetPreheatLogOpts struct
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

# **ListExecutions**
> []HarborExecution ListExecutions(ctx, projectName, preheatPolicyName, optional)
List executions for the given policy

List executions for the given policy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
 **optional** | ***PreheatApiListExecutionsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListExecutionsOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xRequestId** | **optional.String**| An unique ID for the request | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 

### Return type

[**[]HarborExecution**](Execution.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListInstances**
> []HarborInstance ListInstances(ctx, optional)
List P2P provider instances

List P2P provider instances

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***PreheatApiListInstancesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListInstancesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestId** | **optional.String**| An unique ID for the request | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 

### Return type

[**[]HarborInstance**](Instance.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListPolicies**
> []HarborPreheatPolicy ListPolicies(ctx, projectName, optional)
List preheat policies

List preheat policies

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
 **optional** | ***PreheatApiListPoliciesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListPoliciesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRequestId** | **optional.String**| An unique ID for the request | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 

### Return type

[**[]HarborPreheatPolicy**](PreheatPolicy.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListProviders**
> []HarborMetadata ListProviders(ctx, optional)
List P2P providers

List P2P providers

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***PreheatApiListProvidersOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListProvidersOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**[]HarborMetadata**](Metadata.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListProvidersUnderProject**
> []HarborProviderUnderProject ListProvidersUnderProject(ctx, projectName, optional)
Get all providers at project level

Get all providers at project level

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
 **optional** | ***PreheatApiListProvidersUnderProjectOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListProvidersUnderProjectOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xRequestId** | **optional.String**| An unique ID for the request | 

### Return type

[**[]HarborProviderUnderProject**](ProviderUnderProject.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListTasks**
> []HarborTask ListTasks(ctx, projectName, preheatPolicyName, executionId, optional)
List all the related tasks for the given execution

List all the related tasks for the given execution

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
  **executionId** | **int32**| Execution ID | 
 **optional** | ***PreheatApiListTasksOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiListTasksOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **xRequestId** | **optional.String**| An unique ID for the request | 
 **page** | **optional.Int64**| The page number | [default to 1]
 **pageSize** | **optional.Int64**| The size of per page | [default to 10]
 **q** | **optional.String**| Query string to query resources. Supported query patterns are \&quot;exact match(k&#x3D;v)\&quot;, \&quot;fuzzy match(k&#x3D;~v)\&quot;, \&quot;range(k&#x3D;[min~max])\&quot;, \&quot;list with union releationship(k&#x3D;{v1 v2 v3})\&quot; and \&quot;list with intersetion relationship(k&#x3D;(v1 v2 v3))\&quot;. The value of range and list can be string(enclosed by \&quot; or &#x27;), integer or time(in format \&quot;2020-04-09 02:36:00\&quot;). All of these query patterns should be put in the query string \&quot;q&#x3D;xxx\&quot; and splitted by \&quot;,\&quot;. e.g. q&#x3D;k1&#x3D;v1,k2&#x3D;~v2,k3&#x3D;[min~max] | 

### Return type

[**[]HarborTask**](Task.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ManualPreheat**
> ManualPreheat(ctx, body, projectName, preheatPolicyName, optional)
Manual preheat

Manual preheat

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborPreheatPolicy**](HarborPreheatPolicy.md)| The policy schema info | 
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
 **optional** | ***PreheatApiManualPreheatOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiManualPreheatOpts struct
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

# **PingInstances**
> PingInstances(ctx, body, optional)
Ping status of a instance.

This endpoint checks status of a instance, the instance can be given by ID or Endpoint URL (together with credential) 

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborInstance**](HarborInstance.md)| The JSON object of instance. | 
 **optional** | ***PreheatApiPingInstancesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiPingInstancesOpts struct
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

# **StopExecution**
> StopExecution(ctx, body, projectName, preheatPolicyName, executionId, optional)
Stop a execution

Stop a execution

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborExecution**](HarborExecution.md)| The data of execution | 
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
  **executionId** | **int32**| Execution ID | 
 **optional** | ***PreheatApiStopExecutionOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiStopExecutionOpts struct
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

# **UpdateInstance**
> HarborInstanceUpdateResp UpdateInstance(ctx, body, preheatInstanceName, optional)
Update the specified P2P provider instance

Update the specified P2P provider instance

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborInstance**](HarborInstance.md)| The instance to update | 
  **preheatInstanceName** | **string**| Instance Name | 
 **optional** | ***PreheatApiUpdateInstanceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiUpdateInstanceOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **xRequestId** | **optional.**| An unique ID for the request | 

### Return type

[**HarborInstanceUpdateResp**](InstanceUpdateResp.md)

### Authorization

[basic](../README.md#basic)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdatePolicy**
> UpdatePolicy(ctx, body, projectName, preheatPolicyName, optional)
Update preheat policy

Update preheat policy

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**HarborPreheatPolicy**](HarborPreheatPolicy.md)| The policy schema info | 
  **projectName** | **string**| The name of the project | 
  **preheatPolicyName** | **string**| Preheat Policy Name | 
 **optional** | ***PreheatApiUpdatePolicyOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a PreheatApiUpdatePolicyOpts struct
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

