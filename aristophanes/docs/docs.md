# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [aristophanes.proto](#aristophanes-proto)
    - [CloseTraceRequest](#aristophanes-CloseTraceRequest)
    - [DatabaseSpan](#aristophanes-DatabaseSpan)
    - [DatabaseSpanRequest](#aristophanes-DatabaseSpanRequest)
    - [Empty](#aristophanes-Empty)
    - [HealthCheckResponse](#aristophanes-HealthCheckResponse)
    - [ParabasisRequest](#aristophanes-ParabasisRequest)
    - [Span](#aristophanes-Span)
    - [SpanRequest](#aristophanes-SpanRequest)
    - [StartTraceRequest](#aristophanes-StartTraceRequest)
    - [Trace](#aristophanes-Trace)
    - [TraceBare](#aristophanes-TraceBare)
    - [TraceCommon](#aristophanes-TraceCommon)
    - [TraceRequest](#aristophanes-TraceRequest)
    - [TraceResponse](#aristophanes-TraceResponse)
    - [TraceStart](#aristophanes-TraceStart)
    - [TraceStop](#aristophanes-TraceStop)
    - [TracingMetrics](#aristophanes-TracingMetrics)
  
    - [TraceService](#aristophanes-TraceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="aristophanes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## aristophanes.proto



<a name="aristophanes-CloseTraceRequest"></a>

### CloseTraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_body | [string](#string) |  | Optional: Response body data |
| response_code | [int32](#int32) |  |  |






<a name="aristophanes-DatabaseSpan"></a>

### DatabaseSpan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| query | [string](#string) |  | Database query statement |
| time_started | [string](#string) |  | Reuse TraceCommon for common fields |
| time_finished | [string](#string) |  |  |
| hits | [int64](#int64) |  |  |
| took | [string](#string) |  |  |
| common | [TraceCommon](#aristophanes-TraceCommon) |  |  |






<a name="aristophanes-DatabaseSpanRequest"></a>

### DatabaseSpanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | Action performed in the span |
| query | [string](#string) |  | Database query statement |
| hits | [int64](#int64) |  | Number of hits |
| time_took | [int64](#int64) |  |  |






<a name="aristophanes-Empty"></a>

### Empty







<a name="aristophanes-HealthCheckResponse"></a>

### HealthCheckResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [bool](#bool) |  |  |






<a name="aristophanes-ParabasisRequest"></a>

### ParabasisRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| parent_span_id | [string](#string) |  |  |
| span_id | [string](#string) |  |  |
| start_trace | [StartTraceRequest](#aristophanes-StartTraceRequest) |  |  |
| trace | [TraceRequest](#aristophanes-TraceRequest) |  |  |
| close_trace | [CloseTraceRequest](#aristophanes-CloseTraceRequest) |  |  |
| span | [SpanRequest](#aristophanes-SpanRequest) |  |  |
| database_span | [DatabaseSpanRequest](#aristophanes-DatabaseSpanRequest) |  |  |






<a name="aristophanes-Span"></a>

### Span



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | Action performed in the span |
| status | [string](#string) |  |  |
| took | [string](#string) |  |  |
| common | [TraceCommon](#aristophanes-TraceCommon) |  |  |






<a name="aristophanes-SpanRequest"></a>

### SpanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | Action performed in the span |
| status | [string](#string) |  |  |
| took | [string](#string) |  |  |






<a name="aristophanes-StartTraceRequest"></a>

### StartTraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| remote_address | [string](#string) |  | Remote address of the client |
| operation | [string](#string) |  | Graphql operation that generated Trace start |
| root_query | [string](#string) |  | Graphql Root Query |






<a name="aristophanes-Trace"></a>

### Trace



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| common | [TraceCommon](#aristophanes-TraceCommon) |  | Reuse TraceCommon for common fields |
| metrics | [TracingMetrics](#aristophanes-TracingMetrics) |  |  |






<a name="aristophanes-TraceBare"></a>

### TraceBare



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| span_id | [string](#string) |  |  |
| save | [bool](#bool) |  |  |






<a name="aristophanes-TraceCommon"></a>

### TraceCommon
Common message used in various trace-related messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| span_id | [string](#string) |  |  |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| timestamp | [string](#string) |  | Timestamp will be set automatically |
| pod_name | [string](#string) |  | Pod that generated Trace |
| namespace | [string](#string) |  | Namespace that generated Trace |
| item_type | [string](#string) |  | TRACE, SPAN, etc. |






<a name="aristophanes-TraceRequest"></a>

### TraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |






<a name="aristophanes-TraceResponse"></a>

### TraceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ack | [string](#string) |  |  |






<a name="aristophanes-TraceStart"></a>

### TraceStart



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| remote_address | [string](#string) |  | Remote address of the client |
| operation | [string](#string) |  | Graphql operation that generated Trace start |
| root_query | [string](#string) |  | Graphql Root Query |
| common | [TraceCommon](#aristophanes-TraceCommon) |  | Reuse TraceCommon for common fields |
| metrics | [TracingMetrics](#aristophanes-TracingMetrics) |  |  |






<a name="aristophanes-TraceStop"></a>

### TraceStop



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_body | [string](#string) |  | The response generated when a 200 is returned |
| common | [TraceCommon](#aristophanes-TraceCommon) |  | Reuse TraceCommon for common fields |
| metrics | [TracingMetrics](#aristophanes-TracingMetrics) |  |  |






<a name="aristophanes-TracingMetrics"></a>

### TracingMetrics



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| cpu_units | [string](#string) |  |  |
| memory_units | [string](#string) |  |  |
| name | [string](#string) |  |  |
| cpu_raw | [int64](#int64) |  |  |
| memory_raw | [int64](#int64) |  |  |
| cpu_human_readable | [string](#string) |  |  |
| memory_human_readable | [string](#string) |  |  |





 

 

 


<a name="aristophanes-TraceService"></a>

### TraceService
The TraceService service provides operations for managing and tracking traces and spans.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| HealthCheck | [Empty](#aristophanes-Empty) | [HealthCheckResponse](#aristophanes-HealthCheckResponse) |  |
| Chorus | [ParabasisRequest](#aristophanes-ParabasisRequest) stream | [TraceResponse](#aristophanes-TraceResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

