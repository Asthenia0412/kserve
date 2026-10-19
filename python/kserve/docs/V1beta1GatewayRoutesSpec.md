# V1beta1GatewayRoutesSpec

GatewayRoutesSpec defines the configuration for a Gateway API route.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**group** | **str** | Group identifies the routing group this LLMISVC belongs to. All members with the same group share weighted traffic distribution. Must be explicit - not inferred from spec.model.name. | [optional] 
**http** | [**V1beta1HTTPRouteSpec**](V1beta1HTTPRouteSpec.md) |  | [optional] 
**weight** | **int** | Weight is the proportional traffic share for this member. Follows Gateway API backendRef weight semantics (proportional, not percentage). | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


