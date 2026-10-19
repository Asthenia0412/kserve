# V1beta1GroupMemberStatus

GroupMemberStatus reports the observed state of a single group member.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**backend_ref** | [**SigsK8sIoGatewayApiApisV1BackendObjectReference**](SigsK8sIoGatewayApiApisV1BackendObjectReference.md) |  | [optional] 
**name** | **str** | Name of the group member (LLMInferenceService name). | [default to '']
**stopped** | **bool** | Stopped is true when the member has the force-stop annotation set. A stopped member remains in the group but receives no traffic (weight is overridden to 0 regardless of the spec value). | [optional] 
**weight** | **int** | Weight is the effective traffic weight for this member. | [default to 0]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


