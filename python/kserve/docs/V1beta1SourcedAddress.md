# V1beta1SourcedAddress

SourcedAddress extends Addressable with the networking resource that produced this address, enabling consumers to select endpoints by origin.
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ca_certs** | **str** | CACerts is the Certification Authority (CA) certificates in PEM format according to https://www.rfc-editor.org/rfc/rfc7468. | [optional] 
**audience** | **str** | Audience is the OIDC audience for this address. | [optional] 
**models** | [**list[V1beta1ModelSourcedAddressStatus]**](V1beta1ModelSourcedAddressStatus.md) |  | [optional] 
**name** | **str** | Name is the name of the address. | [optional] 
**origin** | [**SigsK8sIoGatewayApiApisV1ObjectReference**](SigsK8sIoGatewayApiApisV1ObjectReference.md) |  | [optional] 
**url** | [**KnativeURL**](KnativeURL.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


