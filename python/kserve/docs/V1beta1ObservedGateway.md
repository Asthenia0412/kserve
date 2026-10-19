# V1beta1ObservedGateway

ObservedGateway is a Gateway reference with the listeners and HTTPRoutes bound to this service through it. Used in status to record observed routing topology.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**group** | **str** | Group is the group of the referent. For example, \&quot;gateway.networking.k8s.io\&quot;. When set to the empty string, core API group is inferred. | [default to '']
**http_routes** | [**list[SigsK8sIoGatewayApiApisV1ObjectReference]**](SigsK8sIoGatewayApiApisV1ObjectReference.md) | HTTPRoutes lists the HTTPRoutes bound to this service through this Gateway. | [optional] 
**kind** | **str** | Kind is kind of the referent. For example \&quot;ConfigMap\&quot; or \&quot;Service\&quot;. | [default to '']
**listeners** | **list[str]** | Listeners lists the SectionNames of the Gateway listeners that accepted routes from this service. Nil means the route targets all listeners (no SectionName was specified on the parentRef). | [optional] 
**name** | **str** | Name is the name of the referent. | [default to '']
**namespace** | **str** | Namespace is the namespace of the referenced object. When unspecified, the local namespace is inferred.  Note that when a namespace different than the local namespace is specified, a ReferenceGrant object is required in the referent namespace to allow that namespace&#39;s owner to accept the reference. See the ReferenceGrant documentation for details.  Support: Core | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


