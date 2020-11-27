# HarborArtifact

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **int64** | The ID of the artifact | [optional] [default to null]
**Type_** | **string** | The type of the artifact, e.g. image, chart, etc | [optional] [default to null]
**MediaType** | **string** | The media type of the artifact | [optional] [default to null]
**ManifestMediaType** | **string** | The manifest media type of the artifact | [optional] [default to null]
**ProjectId** | **int64** | The ID of the project that the artifact belongs to | [optional] [default to null]
**RepositoryId** | **int64** | The ID of the repository that the artifact belongs to | [optional] [default to null]
**Digest** | **string** | The digest of the artifact | [optional] [default to null]
**Size** | **int64** | The size of the artifact | [optional] [default to null]
**Icon** | **string** | The digest of the icon | [optional] [default to null]
**PushTime** | [**time.Time**](time.Time.md) | The push time of the artifact | [optional] [default to null]
**PullTime** | [**time.Time**](time.Time.md) | The latest pull time of the artifact | [optional] [default to null]
**ExtraAttrs** | [***map[string]interface{}**](map.md) |  | [optional] [default to null]
**Annotations** | [***map[string]string**](map.md) |  | [optional] [default to null]
**References** | [**[]HarborReference**](Reference.md) |  | [optional] [default to null]
**Tags** | [**[]HarborTag**](Tag.md) |  | [optional] [default to null]
**AdditionLinks** | [***map[string]HarborAdditionLink**](map.md) |  | [optional] [default to null]
**Labels** | [**[]HarborLabel**](Label.md) |  | [optional] [default to null]
**ScanOverview** | [***map[string]HarborNativeReportSummary**](map.md) |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

