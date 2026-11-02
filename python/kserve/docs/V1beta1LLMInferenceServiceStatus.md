# V1beta1LLMInferenceServiceStatus

LLMInferenceServiceStatus defines the observed state of LLMInferenceService.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**address** | [**KnativeAddressable**](KnativeAddressable.md) |  | [optional] 
**addresses** | [**list[V1beta1SourcedAddress]**](V1beta1SourcedAddress.md) | Addresses lists the network endpoints where the service is reachable. Each address may include an origin reference identifying the networking resource that produced it. | [optional] 
**annotations** | **dict(str, str)** | Annotations is additional Status fields for the Resource to save some additional State as well as convey more information to the user. This is roughly akin to Annotations on any k8s resource, just the reconciler conveying richer information outwards. | [optional] 
**applied_configs** | [**list[V1beta1AppliedConfigRef]**](V1beta1AppliedConfigRef.md) | AppliedConfigRefs records which resources were applied during the last successful reconciliation, in merge precedence order. Well-known configs (determined by the deployment pattern) appear first with lower precedence, followed by explicitly referenced baseRefs with higher precedence. The service&#39;s own spec always takes the highest precedence but is not listed here. Atomic because the controller always writes the full list and ordering encodes merge precedence. | [optional] 
**conditions** | [**list[KnativeCondition]**](KnativeCondition.md) | Conditions the latest available observations of a resource&#39;s current state. | [optional] 
**observed_generation** | **int** | ObservedGeneration is the &#39;Generation&#39; of the Service that was last processed by the controller. | [optional] 
**router** | [**V1beta1RouterStatus**](V1beta1RouterStatus.md) |  | [optional] 
**url** | [**KnativeURL**](KnativeURL.md) |  | [optional] 
**workloads** | [**V1beta1WorkloadStatus**](V1beta1WorkloadStatus.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


