# HarborProjectReq

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectName** | **string** | The name of the project. | [optional] [default to null]
**Public** | **bool** | deprecated, reserved for project creation in replication | [optional] [default to null]
**Metadata** | [***HarborProjectMetadata**](ProjectMetadata.md) |  | [optional] [default to null]
**CveAllowlist** | [***HarborCveAllowlist**](CVEAllowlist.md) |  | [optional] [default to null]
**StorageLimit** | **int64** | The storage quota of the project. | [optional] [default to null]
**RegistryId** | **int64** | The ID of referenced registry when creating the proxy cache project | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

