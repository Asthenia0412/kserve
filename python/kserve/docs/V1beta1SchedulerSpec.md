# V1beta1SchedulerSpec

SchedulerSpec defines the Inference Gateway extension configuration.  The SchedulerSpec configures the connection from the Gateway to the model deployment leveraging the LLM optimized request Scheduler, also known as the Endpoint Picker (EPP). The EPP determines the exact pod that should handle the request and responds back to the Gateway with the target pod. The Gateway will then forward the request to the chosen pod.  The Scheduler is only effective when having multiple inference pod replicas.  Step 1 (endpoint selection): Gateway <-- ExtProc --> EPP (select the optimal replica to handle the request) Step 2 (endpoint routing): Gateway <-- forward request/response --> Inference Pod X
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict(str, str)** | Annotations that will be added to the scheduler component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ | [optional] 
**config** | [**V1beta1SchedulerConfigSpec**](V1beta1SchedulerConfigSpec.md) |  | [optional] 
**labels** | **dict(str, str)** | Labels that will be added to the scheduler component pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ | [optional] 
**pool** | [**V1beta1InferencePoolSpec**](V1beta1InferencePoolSpec.md) |  | [optional] 
**replicas** | **int** | Replicas is the number of replicas for the scheduler. | [optional] 
**template** | [**V1PodSpec**](V1PodSpec.md) |  | [optional] 
**tokenizer** | [**V1beta1TokenizerSpec**](V1beta1TokenizerSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


