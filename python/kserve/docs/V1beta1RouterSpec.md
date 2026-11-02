# V1beta1RouterSpec

RouterSpec defines the routing configuration for exposing the service. It supports Kubernetes Ingress and the Gateway API. The fields are mutually exclusive.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**gateway** | [**V1beta1GatewaySpec**](V1beta1GatewaySpec.md) |  | [optional] 
**ingress** | [**V1beta1IngressSpec**](V1beta1IngressSpec.md) |  | [optional] 
**route** | [**V1beta1GatewayRoutesSpec**](V1beta1GatewayRoutesSpec.md) |  | [optional] 
**scheduler** | [**V1beta1SchedulerSpec**](V1beta1SchedulerSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


