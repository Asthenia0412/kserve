# V1beta1GatewayObjectReference

GatewayObjectReference is a reference to a Gateway resource. It extends UntypedObjectReference with Gateway-specific fields.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | Name of the referenced object. | [optional] 
**namespace** | **str** | Namespace of the referenced object. | [optional] 
**section_name** | **str** | SectionName is the name of a section within the target resource. When set on a Gateway reference, it targets a specific listener by name. When unset, the route is attached to all listeners on the referenced Gateway that support the route type. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


