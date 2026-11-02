# V1beta1WVASpec

WVASpec is the deprecated Workload Variant Autoscaler configuration retained for upgrade compatibility. New WVA configuration is rejected and existing WVA configuration is not reconciled.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hpa** | [**V1beta1HPAScalingSpec**](V1beta1HPAScalingSpec.md) |  | [optional] 
**keda** | [**V1beta1KEDAScalingSpec**](V1beta1KEDAScalingSpec.md) |  | [optional] 
**variant_cost** | **str** | VariantCost specifies the cost per replica for this variant (used in saturation analysis). | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


