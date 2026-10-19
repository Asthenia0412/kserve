# V1beta1LLMInferenceServiceSpec

LLMInferenceServiceSpec defines the desired state of LLMInferenceService.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict(str, str)** | Annotations that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | [optional] 
**base_refs** | [**list[V1LocalObjectReference]**](https://github.com/kubernetes-client/python/blob/master/kubernetes/docs/V1LocalObjectReference.md) | BaseRefs allows inheriting and overriding configurations from one or more LLMInferenceServiceConfig instances. The controller merges these base configurations, with the current LLMInferenceService spec taking the highest precedence. When multiple baseRefs are provided, the last one in the list overrides previous ones. | [optional] 
**kv_cache_offloading** | [**V1beta1KVCacheOffloadingSpec**](V1beta1KVCacheOffloadingSpec.md) |  | [optional] 
**labels** | **dict(str, str)** | Labels that will be added to the component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ | [optional] 
**model** | [**V1beta1LLMModelSpec**](V1beta1LLMModelSpec.md) |  | [optional] 
**parallelism** | [**V1beta1ParallelismSpec**](V1beta1ParallelismSpec.md) |  | [optional] 
**prefill** | [**V1beta1WorkloadSpec**](V1beta1WorkloadSpec.md) |  | [optional] 
**replicas** | **int** | Number of replicas for the deployment. | [optional] 
**rollout_strategy** | [**V1beta1RolloutStrategy**](V1beta1RolloutStrategy.md) |  | [optional] 
**router** | [**V1beta1RouterSpec**](V1beta1RouterSpec.md) |  | [optional] 
**runtime** | **str** | Runtime is the name of a ServingRuntime (namespaced) or ClusterServingRuntime that supplies the base container spec — primarily the container image — for the inference workload. The controller resolves this name against ServingRuntime in the LLMInferenceService&#39;s namespace first, then falls back to ClusterServingRuntime. The resolved container spec is prepended as the lowest-priority layer in the merge chain, so LLMInferenceServiceConfig baseRefs and spec.template still override it.  When omitted, the controller uses the default vLLM template (kserve-config-llm-template) which ships its own image. | [optional] 
**scaling** | [**V1beta1ScalingSpec**](V1beta1ScalingSpec.md) |  | [optional] 
**storage_initializer** | [**V1beta1StorageInitializerSpec**](V1beta1StorageInitializerSpec.md) |  | [optional] 
**template** | [**V1PodSpec**](V1PodSpec.md) |  | [optional] 
**tracing** | [**V1beta1TracingSpec**](V1beta1TracingSpec.md) |  | [optional] 
**worker** | [**V1PodSpec**](V1PodSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


