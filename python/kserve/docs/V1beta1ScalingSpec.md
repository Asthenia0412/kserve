# V1beta1ScalingSpec

ScalingSpec configures autoscaling for the LLM inference deployment. When scaling is configured, the controller creates and manages autoscaling resources (ServiceMonitor and the selected actuator — HPA or KEDA ScaledObject). Direct KEDA scaling uses user-defined triggers. The deprecated WVA fields are retained only so existing resources can be decoded during upgrades; new WVA configuration is rejected.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**keda** | [**V1beta1DirectKEDAScalingSpec**](V1beta1DirectKEDAScalingSpec.md) |  | [optional] 
**max_replicas** | **int** | MaxReplicas is the maximum number of replicas for the deployment. | [default to 0]
**min_replicas** | **int** | MinReplicas is the minimum number of replicas for the deployment during active scaling. This is the scaling floor when triggers are active. For idle scale-down, use KEDA&#39;s idleReplicaCount instead. Defaults to 1 if not specified. | [optional] 
**wva** | [**V1beta1WVASpec**](V1beta1WVASpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


