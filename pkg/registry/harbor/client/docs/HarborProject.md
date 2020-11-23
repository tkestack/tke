# HarborProject

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ProjectId** | **int32** | Project ID | [optional] [default to null]
**OwnerId** | **int32** | The owner ID of the project always means the creator of the project. | [optional] [default to null]
**Name** | **string** | The name of the project. | [optional] [default to null]
**RegistryId** | **int64** | The ID of referenced registry when the project is a proxy cache project. | [optional] [default to null]
**CreationTime** | [**time.Time**](time.Time.md) | The creation time of the project. | [optional] [default to null]
**UpdateTime** | [**time.Time**](time.Time.md) | The update time of the project. | [optional] [default to null]
**Deleted** | **bool** | A deletion mark of the project. | [optional] [default to null]
**OwnerName** | **string** | The owner name of the project. | [optional] [default to null]
**Togglable** | **bool** | Correspond to the UI about whether the project&#x27;s publicity is  updatable (for UI) | [optional] [default to null]
**CurrentUserRoleId** | **int32** | The role ID with highest permission of the current user who triggered the API (for UI).  This attribute is deprecated and will be removed in future versions. | [optional] [default to null]
**CurrentUserRoleIds** | **[]int32** | The list of role ID of the current user who triggered the API (for UI) | [optional] [default to null]
**RepoCount** | **int32** | The number of the repositories under this project. | [optional] [default to null]
**ChartCount** | **int32** | The total number of charts under this project. | [optional] [default to null]
**Metadata** | [***HarborProjectMetadata**](ProjectMetadata.md) |  | [optional] [default to null]
**CveAllowlist** | [***HarborCveAllowlist**](CVEAllowlist.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

