# V1beta1AppliedConfigRef

AppliedConfigRef identifies a resource that contributed to the final merged configuration during reconciliation.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | Name of the LLMInferenceServiceConfig, ServingRuntime, or ClusterServingRuntime resource that was applied. | [default to '']
**namespace** | **str** | Namespace where the LLMInferenceServiceConfig or ServingRuntime was resolved from. Omitted for cluster-scoped resources such as ClusterServingRuntime. | [optional] 
**source** | **str** | Source indicates where this config came from: one KServe ships and selected itself, one the service referenced, or a container spec contributed by a ServingRuntime. A config resolved outside the KServe namespace reports UserRef even when it answers to a well-known preset name. | [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


