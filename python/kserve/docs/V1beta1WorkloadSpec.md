# V1beta1WorkloadSpec

WorkloadSpec defines the configuration for a deployment workload, such as replicas and pod specifications.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict(str, str)** | Annotations that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | [optional] 
**kv_cache_offloading** | [**V1beta1KVCacheOffloadingSpec**](V1beta1KVCacheOffloadingSpec.md) |  | [optional] 
**labels** | **dict(str, str)** | Labels that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ | [optional] 
**parallelism** | [**V1beta1ParallelismSpec**](V1beta1ParallelismSpec.md) |  | [optional] 
**replicas** | **int** | Number of replicas for the deployment. | [optional] 
**rollout_strategy** | [**V1beta1RolloutStrategy**](V1beta1RolloutStrategy.md) |  | [optional] 
**scaling** | [**V1beta1ScalingSpec**](V1beta1ScalingSpec.md) |  | [optional] 
**template** | [**V1PodSpec**](V1PodSpec.md) |  | [optional] 
**worker** | [**V1PodSpec**](V1PodSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


