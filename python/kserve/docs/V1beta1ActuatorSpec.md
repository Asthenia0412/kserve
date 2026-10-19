# V1beta1ActuatorSpec

ActuatorSpec defines the autoscaling actuator backend for WVA. Exactly one of HPA or KEDA must be specified.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hpa** | [**V1beta1HPAScalingSpec**](V1beta1HPAScalingSpec.md) |  | [optional] 
**keda** | [**V1beta1KEDAScalingSpec**](V1beta1KEDAScalingSpec.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


