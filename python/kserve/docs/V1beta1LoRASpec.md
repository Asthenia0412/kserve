# V1beta1LoRASpec

LoRASpec defines the configuration for LoRA adapters.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**adapters** | [**list[V1beta1LLMModelSpec]**](V1beta1LLMModelSpec.md) | Adapters is a list of LoRA (Low-Rank Adaptation) adapters to attach to the base model. Each adapter is specified by name and URI (supports hf://, s3://, and pvc:// schemes). The controller automatically downloads adapters and configures the runtime to use them. This type is recursive https://github.com/kubernetes-sigs/controller-tools/issues/585 | [optional] 
**max_adapters** | **int** | MaxAdapters is the maximum number of LoRA adapters that can be loaded simultaneously (maps to vLLM --max-loras). Defaults to the number of configured adapters. | [optional] 
**max_cpu_adapters** | **int** | MaxCpuAdapters is the maximum number of LoRA adapters stored in CPU memory (maps to vLLM --max-cpu-loras). Defaults to the number of configured adapters. | [optional] 
**max_rank** | **int** | MaxRank is the maximum LoRA rank supported by the runtime (maps to vLLM --max-lora-rank). Higher values allow adapters with higher rank but increase memory usage. If not set, vLLM&#39;s default applies (16). | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


