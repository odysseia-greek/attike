# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [aristophanes.proto](#aristophanes-proto)
    - [CloseTraceRequest](#proto-CloseTraceRequest)
    - [DatabaseSpan](#proto-DatabaseSpan)
    - [DatabaseSpanRequest](#proto-DatabaseSpanRequest)
    - [Span](#proto-Span)
    - [SpanRequest](#proto-SpanRequest)
    - [StartSpanRequest](#proto-StartSpanRequest)
    - [StartTraceRequest](#proto-StartTraceRequest)
    - [Trace](#proto-Trace)
    - [TraceCommon](#proto-TraceCommon)
    - [TraceRequest](#proto-TraceRequest)
    - [TraceResponse](#proto-TraceResponse)
    - [TraceStart](#proto-TraceStart)
    - [TraceStop](#proto-TraceStop)
  
    - [TraceService](#proto-TraceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="aristophanes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## aristophanes.proto



<a name="proto-CloseTraceRequest"></a>

### CloseTraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  | Root trace_id which where item will be added |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| response_body | [string](#string) |  | Optional: Response body data |
| response_code | [int32](#int32) |  |  |






<a name="proto-DatabaseSpan"></a>

### DatabaseSpan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| query | [string](#string) |  | Database query statement |
| result_json | [string](#string) |  | Query result data as JSON string |
| common | [TraceCommon](#proto-TraceCommon) |  | Reuse TraceCommon for common fields |






<a name="proto-DatabaseSpanRequest"></a>

### DatabaseSpanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| action | [string](#string) |  | Action performed in the span |
| query | [string](#string) |  | Database query statement |
| result_json | [string](#string) |  | Query result data as JSON string |






<a name="proto-Span"></a>

### Span



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | Action performed in the span |
| request_body | [string](#string) |  | Optional: Request body data |
| response_body | [string](#string) |  | Optional: Response body data |
| common | [TraceCommon](#proto-TraceCommon) |  | Reuse TraceCommon for common fields |






<a name="proto-SpanRequest"></a>

### SpanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  | Root trace_id which where item will be added |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| action | [string](#string) |  | Action performed in the span |
| request_body | [string](#string) |  | Optional: Request body data |
| response_body | [string](#string) |  | Optional: Response body data |






<a name="proto-StartSpanRequest"></a>

### StartSpanRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  | Root trace_id which where item will be added |






<a name="proto-StartTraceRequest"></a>

### StartTraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| remote_address | [string](#string) |  | Remote address of the client |
| operation | [string](#string) |  | Graphql operation that generated Trace start |
| root_query | [string](#string) |  | Graphql Root Query |






<a name="proto-Trace"></a>

### Trace



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| common | [TraceCommon](#proto-TraceCommon) |  | Reuse TraceCommon for common fields |






<a name="proto-TraceCommon"></a>

### TraceCommon
Common message used in various trace-related messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| span_id | [string](#string) |  | The span_id will be generated automatically |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| timestamp | [string](#string) |  | Timestamp will be set automatically |
| pod_name | [string](#string) |  | Pod that generated Trace |
| namespace | [string](#string) |  | Namespace that generated Trace |
| item_type | [string](#string) |  | TRACE, SPAN, etc. |






<a name="proto-TraceRequest"></a>

### TraceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  | Root trace_id which where item will be added |
| parent_span_id | [string](#string) |  | The root parent_span_id |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |






<a name="proto-TraceResponse"></a>

### TraceResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| combined_id | [string](#string) |  | Combination of the trace_id, parent_span_id and sampling bool for example: 841a4f73-ba5b-4c38-9237-e1ad91459028&#43;70b993de1e2f879d&#43;1 |






<a name="proto-TraceStart"></a>

### TraceStart



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | GET POST PUT, etc. |
| url | [string](#string) |  | The URL called by a client |
| host | [string](#string) |  | The host that generated the trace |
| remote_address | [string](#string) |  | Remote address of the client |
| operation | [string](#string) |  | Graphql operation that generated Trace start |
| root_query | [string](#string) |  | Graphql Root Query |
| common | [TraceCommon](#proto-TraceCommon) |  | Reuse TraceCommon for common fields |






<a name="proto-TraceStop"></a>

### TraceStop



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_body | [string](#string) |  | The response generated when a 200 is returned |
| common | [TraceCommon](#proto-TraceCommon) |  | Reuse TraceCommon for common fields |





 

 

 


<a name="proto-TraceService"></a>

### TraceService
The TraceService service provides operations for managing and tracking traces and spans.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StartTrace | [StartTraceRequest](#proto-StartTraceRequest) | [TraceResponse](#proto-TraceResponse) | Start a new trace. |
| Trace | [TraceRequest](#proto-TraceRequest) | [TraceResponse](#proto-TraceResponse) | Record a new trace within an existing trace. |
| StartNewSpan | [StartSpanRequest](#proto-StartSpanRequest) | [TraceResponse](#proto-TraceResponse) | Start a new span within an existing trace. |
| Span | [SpanRequest](#proto-SpanRequest) | [TraceResponse](#proto-TraceResponse) | Record a span with details of an action performed. |
| DatabaseSpan | [DatabaseSpanRequest](#proto-DatabaseSpanRequest) | [TraceResponse](#proto-TraceResponse) | Record a span related to a database query. |
| CloseTrace | [CloseTraceRequest](#proto-CloseTraceRequest) | [TraceResponse](#proto-TraceResponse) | Close an existing trace. |

 



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

