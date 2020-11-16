# HarborCveAllowlist

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int32** | ID of the allowlist | [optional] [default to null]
**ProjectId** | **int32** | ID of the project which the allowlist belongs to.  For system level allowlist this attribute is zero. | [optional] [default to null]
**ExpiresAt** | **int32** | the time for expiration of the allowlist, in the form of seconds since epoch.  This is an optional attribute, if it&#x27;s not set the CVE allowlist does not expire. | [optional] [default to null]
**Items** | [**[]HarborCveAllowlistItem**](CVEAllowlistItem.md) |  | [optional] [default to null]
**CreationTime** | [**time.Time**](time.Time.md) | The creation time of the allowlist. | [optional] [default to null]
**UpdateTime** | [**time.Time**](time.Time.md) | The update time of the allowlist. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

