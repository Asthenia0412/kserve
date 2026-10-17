# V1beta1FileSystemTierSpec

FileSystemTierSpec configures a POSIX disk secondary tier. Exactly one of EmptyDir or PVC must be set.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**empty_dir** | [**V1beta1EmptyDirTierSpec**](V1beta1EmptyDirTierSpec.md) |  | [optional] 
**pvc** | [**V1beta1PVCTierSpec**](V1beta1PVCTierSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


