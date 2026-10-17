# V1beta1LLMModelSpec

ConfidentialSpec enables confidential model serving with encrypted model artifacts. When enabled, the storage initializer will decrypt model files using keys obtained from a Key Broker Service (KBS) via TEE attestation.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**confidential** | [**V1beta1ConfidentialSpec**](V1beta1ConfidentialSpec.md) |  | [optional] 
**lora** | [**V1beta1LoRASpec**](V1beta1LoRASpec.md) |  | [optional] 
**name** | **str** | Name is the name of the model as it will be set in the \&quot;model\&quot; parameter for an incoming request. If omitted, it will default to &#x60;metadata.name&#x60;. For LoRA adapters, this field is required. | [optional] 
**uri** | [**KnativeURL**](KnativeURL.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


