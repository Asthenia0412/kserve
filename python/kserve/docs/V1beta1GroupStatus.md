# V1beta1GroupStatus

GroupStatus reports the observed state of the routing group this LLMISVC belongs to.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**members** | [**list[V1beta1GroupMemberStatus]**](V1beta1GroupMemberStatus.md) | Members lists all observed members of the group with their weights and resolved backend references. | [optional] 
**name** | **str** | Name of the routing group. | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


