# V1beta1StorageInitializerSpec

StorageInitializerSpec defines the configuration for the storage initializer. The storage initializer is an initContainer responsible for downloading model artifacts from remote storage (s3://, hf://) before the main container starts.  Example - Disable storage initializer:   storageInitializer:    enabled: false  Example - Explicitly enable (same as default):   storageInitializer:    enabled: true
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** | Enabled controls whether the storage-initializer initContainer is created. When nil or true, storage-initializer is created for applicable URIs (s3://, hf://). When explicitly set to false, storage-initializer creation is skipped. This is useful when models are pre-loaded via alternative mechanisms (e.g., custom init containers, modelcars). Default: true (nil is treated as true for backward compatibility) | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


