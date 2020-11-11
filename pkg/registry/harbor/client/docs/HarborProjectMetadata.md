# HarborProjectMetadata

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Public** | **string** | The public status of the project. The valid values are \&quot;true\&quot;, \&quot;false\&quot;. | [optional] [default to null]
**EnableContentTrust** | **string** | Whether content trust is enabled or not. If it is enabled, user can&#x27;t pull unsigned images from this project. The valid values are \&quot;true\&quot;, \&quot;false\&quot;. | [optional] [default to null]
**PreventVul** | **string** | Whether prevent the vulnerable images from running. The valid values are \&quot;true\&quot;, \&quot;false\&quot;. | [optional] [default to null]
**Severity** | **string** | If the vulnerability is high than severity defined here, the images can&#x27;t be pulled. The valid values are \&quot;none\&quot;, \&quot;low\&quot;, \&quot;medium\&quot;, \&quot;high\&quot;, \&quot;critical\&quot;. | [optional] [default to null]
**AutoScan** | **string** | Whether scan images automatically when pushing. The valid values are \&quot;true\&quot;, \&quot;false\&quot;. | [optional] [default to null]
**ReuseSysCveAllowlist** | **string** | Whether this project reuse the system level CVE allowlist as the allowlist of its own.  The valid values are \&quot;true\&quot;, \&quot;false\&quot;. If it is set to \&quot;true\&quot; the actual allowlist associate with this project, if any, will be ignored. | [optional] [default to null]
**RetentionId** | **string** | The ID of the tag retention policy for the project | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

