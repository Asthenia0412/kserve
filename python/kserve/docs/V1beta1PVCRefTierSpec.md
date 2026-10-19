# V1beta1PVCRefTierSpec

PVCRefTierSpec mounts a pre-existing user-managed PVC as a KV cache tier.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | Name of the pre-existing PersistentVolumeClaim. | [optional] [default to '']
**path** | **str** | Path is a subdirectory within the PVC used as VolumeMount.subPath. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


