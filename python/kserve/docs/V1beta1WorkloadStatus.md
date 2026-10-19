# V1beta1WorkloadStatus

WorkloadStatus records the workload resources observed during the last successful reconciliation. Nil when no workload resources have been created yet, or when the service is stopped.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**prefill** | [**V1beta1ObservedWorkloadStatus**](V1beta1ObservedWorkloadStatus.md) |  | [optional] 
**primary** | [**V1beta1ObservedWorkloadStatus**](V1beta1ObservedWorkloadStatus.md) |  | [optional] 
**scheduler** | [**V1beta1ObservedWorkloadStatus**](V1beta1ObservedWorkloadStatus.md) |  | [optional] 
**service** | [**V1TypedLocalObjectReference**](V1TypedLocalObjectReference.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


