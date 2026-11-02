# V1beta1RouterStatus

RouterStatus records the networking resources observed during the last successful routing reconciliation. Nil when routing is not configured or the service is stopped.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**gateways** | [**list[V1beta1ObservedGateway]**](V1beta1ObservedGateway.md) | Gateways lists the Gateway resources observed as attached to this service, each with the listeners and HTTPRoutes bound through them. | [optional] 
**group** | [**V1beta1GroupStatus**](V1beta1GroupStatus.md) |  | [optional] 
**scheduler** | [**V1beta1ObservedSchedulerStatus**](V1beta1ObservedSchedulerStatus.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


