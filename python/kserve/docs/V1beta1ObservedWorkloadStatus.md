# V1beta1ObservedWorkloadStatus

ObservedWorkloadStatus identifies a workload resource and its observed replica state.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_group** | **str** | APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required. | [optional] 
**kind** | **str** | Kind is the type of resource being referenced | [default to '']
**name** | **str** | Name is the name of resource being referenced | [default to '']
**ready_replicas** | **int** | ReadyReplicas is the number of pods available to serve traffic. Copied from the workload&#39;s status on each reconcile. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


