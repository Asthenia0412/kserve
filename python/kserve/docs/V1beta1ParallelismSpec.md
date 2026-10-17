# V1beta1ParallelismSpec

TracingSpec defines the distributed tracing configuration. When present (even as an empty object `{}`), tracing is enabled with sensible defaults. When omitted, no tracing instrumentation is injected.  Example - Enable tracing with defaults:   tracing: {}  Example - Custom configuration:   tracing:    exporterEndpoint: \"http://my-collector:4317\"    sampler: \"parentbased_traceidratio\"    samplerArg: \"0.1\"    exporter: \"otlp\"
## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**data** | **int** | Data parallelism size. | [optional] 
**data_local** | **int** | DataLocal data local parallelism size. | [optional] 
**data_rpc_port** | **int** | DataRPCPort is the data parallelism RPC port. | [optional] 
**expert** | **bool** | Expert enables expert parallelism. | [optional] 
**pipeline** | **int** | Pipeline parallelism size. | [optional] 
**tensor** | **int** | Tensor parallelism size. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


