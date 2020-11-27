# HarborNativeReportSummary

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ReportId** | **string** | id of the native scan report | [optional] [default to null]
**ScanStatus** | **string** | The status of the report generating process | [optional] [default to null]
**Severity** | **string** | The overall severity | [optional] [default to null]
**Duration** | **int64** | The seconds spent for generating the report | [optional] [default to null]
**Summary** | [***HarborVulnerabilitySummary**](VulnerabilitySummary.md) |  | [optional] [default to null]
**StartTime** | [**time.Time**](time.Time.md) | The start time of the scan process that generating report | [optional] [default to null]
**EndTime** | [**time.Time**](time.Time.md) | The end time of the scan process that generating report | [optional] [default to null]
**CompletePercent** | **int32** | The complete percent of the scanning which value is between 0 and 100 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

