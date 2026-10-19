# V1beta1KVCacheOffloadingSpec

KVCacheOffloadingSpec configures KV cache offloading via vLLM's OffloadingConnector.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**cpu** | [**ResourceQuantity**](ResourceQuantity.md) |  | 
**eviction_policy** | **str** | EvictionPolicy for the primary CPU KV cache tier. Defaults to \&quot;lru\&quot;. | [optional] 
**secondary** | [**list[V1beta1SecondaryTierSpec]**](V1beta1SecondaryTierSpec.md) | Secondary is an ordered list of secondary KV cache tiers. vLLM cascades through tiers in the order listed. Currently only fileSystem tiers are supported; the array shape is designed to accommodate object store tiers in a future follow-up without an API break. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


