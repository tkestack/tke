# HarborAuditLog

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | The ID of the audit log entry. | [optional] [default to null]
**Username** | **string** | Username of the user in this log entry. | [optional] [default to null]
**Resource** | **string** | Name of the repository in this log entry. | [optional] [default to null]
**ResourceType** | **string** | Tag of the repository in this log entry. | [optional] [default to null]
**Operation** | **string** | The operation against the repository in this log entry. | [optional] [default to null]
**OpTime** | [**time.Time**](time.Time.md) | The time when this operation is triggered. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

