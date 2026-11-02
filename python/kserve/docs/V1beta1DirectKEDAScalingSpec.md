# V1beta1DirectKEDAScalingSpec

DirectKEDAScalingSpec configures direct KEDA scaling. It embeds KEDAScalingSpec for shared ScaledObject settings and adds user-defined triggers.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**advanced** | [**GithubComKedacoreKedaV2ApisKedaV1alpha1AdvancedConfig**](GithubComKedacoreKedaV2ApisKedaV1alpha1AdvancedConfig.md) |  | [optional] 
**cooldown_period** | **int** | CooldownPeriod is the period in seconds to wait after the last trigger reported active before scaling the resource back to its minimum replica count. A value of 0 means scale down immediately with no cooldown. | [optional] 
**fallback** | [**GithubComKedacoreKedaV2ApisKedaV1alpha1Fallback**](GithubComKedacoreKedaV2ApisKedaV1alpha1Fallback.md) |  | [optional] 
**idle_replica_count** | **int** | IdleReplicaCount is the number of replicas KEDA will scale the resource down to when there are no triggers active. This must be less than minReplicas. Set to 0 for true scale-to-zero: KEDA scales the resource to zero replicas when idle and immediately back to minReplicas once a trigger becomes active. If not set, KEDA will not scale below minReplicas. | [optional] 
**initial_cooldown_period** | **int** | InitialCooldownPeriod is the period in seconds to wait after the ScaledObject is created before KEDA starts evaluating triggers. Useful for LLM deployments where the model takes time to load before it can serve traffic, preventing premature scale-up decisions. | [optional] 
**polling_interval** | **int** | PollingInterval is the interval in seconds to check each trigger on. Must be at least 1 second. | [optional] 
**triggers** | [**list[GithubComKedacoreKedaV2ApisKedaV1alpha1ScaleTriggers]**](GithubComKedacoreKedaV2ApisKedaV1alpha1ScaleTriggers.md) | Triggers defines the KEDA scaling triggers. At least one trigger is required. | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


