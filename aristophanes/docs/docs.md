# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [v1/aristophanes.proto](#v1_aristophanes-proto)
    - [ActionEvent](#aristophanes-v1-ActionEvent)
    - [AttikeEvent](#aristophanes-v1-AttikeEvent)
    - [DatabaseSpanEvent](#aristophanes-v1-DatabaseSpanEvent)
    - [Empty](#aristophanes-v1-Empty)
    - [GraphQLEvent](#aristophanes-v1-GraphQLEvent)
    - [HealthCheckResponse](#aristophanes-v1-HealthCheckResponse)
    - [ObserveAction](#aristophanes-v1-ObserveAction)
    - [ObserveDbSpan](#aristophanes-v1-ObserveDbSpan)
    - [ObserveGraphQL](#aristophanes-v1-ObserveGraphQL)
    - [ObserveRequest](#aristophanes-v1-ObserveRequest)
    - [ObserveResponse](#aristophanes-v1-ObserveResponse)
    - [ObserveTraceHop](#aristophanes-v1-ObserveTraceHop)
    - [ObserveTraceHopStop](#aristophanes-v1-ObserveTraceHopStop)
    - [ObserveTraceStart](#aristophanes-v1-ObserveTraceStart)
    - [ObserveTraceStop](#aristophanes-v1-ObserveTraceStop)
    - [TraceBare](#aristophanes-v1-TraceBare)
    - [TraceCommon](#aristophanes-v1-TraceCommon)
    - [TraceHopEvent](#aristophanes-v1-TraceHopEvent)
    - [TraceHopStopEvent](#aristophanes-v1-TraceHopStopEvent)
    - [TraceStartEvent](#aristophanes-v1-TraceStartEvent)
    - [TraceStopEvent](#aristophanes-v1-TraceStopEvent)
  
    - [ItemType](#aristophanes-v1-ItemType)
  
    - [TraceService](#aristophanes-v1-TraceService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="v1_aristophanes-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## v1/aristophanes.proto



<a name="aristophanes-v1-ActionEvent"></a>

### ActionEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  |  |
| status | [string](#string) |  |  |
| took_ms | [int64](#int64) |  |  |






<a name="aristophanes-v1-AttikeEvent"></a>

### AttikeEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| common | [TraceCommon](#aristophanes-v1-TraceCommon) |  |  |
| trace_start | [TraceStartEvent](#aristophanes-v1-TraceStartEvent) |  |  |
| trace_hop | [TraceHopEvent](#aristophanes-v1-TraceHopEvent) |  |  |
| graphql | [GraphQLEvent](#aristophanes-v1-GraphQLEvent) |  |  |
| action | [ActionEvent](#aristophanes-v1-ActionEvent) |  |  |
| db_span | [DatabaseSpanEvent](#aristophanes-v1-DatabaseSpanEvent) |  |  |
| trace_stop | [TraceStopEvent](#aristophanes-v1-TraceStopEvent) |  |  |
| trace_hop_stop | [TraceHopStopEvent](#aristophanes-v1-TraceHopStopEvent) |  |  |






<a name="aristophanes-v1-DatabaseSpanEvent"></a>

### DatabaseSpanEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  |  |
| query | [string](#string) |  |  |
| hits | [int64](#int64) |  |  |
| took_ms | [int64](#int64) |  |  |
| target | [string](#string) |  |  |
| index | [string](#string) |  |  |






<a name="aristophanes-v1-Empty"></a>

### Empty







<a name="aristophanes-v1-GraphQLEvent"></a>

### GraphQLEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [string](#string) |  |  |
| root_query | [string](#string) |  |  |






<a name="aristophanes-v1-HealthCheckResponse"></a>

### HealthCheckResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [bool](#bool) |  |  |






<a name="aristophanes-v1-ObserveAction"></a>

### ObserveAction



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | A single in-service action |
| status | [string](#string) |  |  |
| took_ms | [int64](#int64) |  |  |






<a name="aristophanes-v1-ObserveDbSpan"></a>

### ObserveDbSpan



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| action | [string](#string) |  | A database/search query (Elastic etc.)

e.g. &#34;elastic.search&#34; |
| query | [string](#string) |  | query body/string (careful with size) |
| hits | [int64](#int64) |  |  |
| took_ms | [int64](#int64) |  |  |
| target | [string](#string) |  | &#34;elasticsearch&#34;, &#34;cassandra&#34;, ... |
| index | [string](#string) |  | optional index/table |






<a name="aristophanes-v1-ObserveGraphQL"></a>

### ObserveGraphQL



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operation | [string](#string) |  | GraphQL-specific context. Emit once near the start (or whenever known).

GraphQL operation name |
| root_query | [string](#string) |  | GraphQL query string (careful with size) |






<a name="aristophanes-v1-ObserveRequest"></a>

### ObserveRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| span_id | [string](#string) |  | unique id for this event |
| parent_span_id | [string](#string) |  | used to form a tree (optional) |
| trace_start | [ObserveTraceStart](#aristophanes-v1-ObserveTraceStart) |  |  |
| trace_hop | [ObserveTraceHop](#aristophanes-v1-ObserveTraceHop) |  |  |
| trace_hop_stop | [ObserveTraceHopStop](#aristophanes-v1-ObserveTraceHopStop) |  |  |
| graphql | [ObserveGraphQL](#aristophanes-v1-ObserveGraphQL) |  |  |
| action | [ObserveAction](#aristophanes-v1-ObserveAction) |  |  |
| db_span | [ObserveDbSpan](#aristophanes-v1-ObserveDbSpan) |  |  |
| trace_stop | [ObserveTraceStop](#aristophanes-v1-ObserveTraceStop) |  |  |






<a name="aristophanes-v1-ObserveResponse"></a>

### ObserveResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ack | [string](#string) |  |  |






<a name="aristophanes-v1-ObserveTraceHop"></a>

### ObserveTraceHop



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | A recorded step during the trace (service-to-service call, internal hop, etc.)

&#34;GET&#34;, &#34;POST&#34;, &#34;grpc&#34;, &#34;internal&#34; |
| url | [string](#string) |  | path or URL (for grpc you can use &#34;/pkg.Service/Method&#34;) |
| host | [string](#string) |  | destination service/host |






<a name="aristophanes-v1-ObserveTraceHopStop"></a>

### ObserveTraceHopStop



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_code | [int32](#int32) |  |  |
| took_ms | [int64](#int64) |  | duration of this hop |






<a name="aristophanes-v1-ObserveTraceStart"></a>

### ObserveTraceStart



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | Usually the inbound request that created the trace |
| url | [string](#string) |  |  |
| host | [string](#string) |  |  |
| remote_address | [string](#string) |  |  |
| root_query | [string](#string) |  |  |
| operation | [string](#string) |  |  |






<a name="aristophanes-v1-ObserveTraceStop"></a>

### ObserveTraceStop



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_body | [string](#string) |  |  |
| response_code | [int32](#int32) |  |  |






<a name="aristophanes-v1-TraceBare"></a>

### TraceBare
Helper for middleware adapater


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| span_id | [string](#string) |  |  |
| save | [bool](#bool) |  |  |






<a name="aristophanes-v1-TraceCommon"></a>

### TraceCommon



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| trace_id | [string](#string) |  |  |
| span_id | [string](#string) |  |  |
| parent_span_id | [string](#string) |  |  |
| timestamp | [string](#string) |  | UTC &#34;2006-01-02T15:04:05.000&#34; |
| pod_name | [string](#string) |  |  |
| namespace | [string](#string) |  |  |
| item_type | [ItemType](#aristophanes-v1-ItemType) |  |  |






<a name="aristophanes-v1-TraceHopEvent"></a>

### TraceHopEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  |  |
| url | [string](#string) |  |  |
| host | [string](#string) |  |  |
| response_code | [int32](#int32) |  |  |
| took_ms | [int64](#int64) |  |  |






<a name="aristophanes-v1-TraceHopStopEvent"></a>

### TraceHopStopEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_code | [int32](#int32) |  |  |
| took_ms | [int64](#int64) |  |  |






<a name="aristophanes-v1-TraceStartEvent"></a>

### TraceStartEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  |  |
| url | [string](#string) |  |  |
| host | [string](#string) |  |  |
| remote_address | [string](#string) |  |  |
| root_query | [string](#string) |  |  |
| operation | [string](#string) |  |  |






<a name="aristophanes-v1-TraceStopEvent"></a>

### TraceStopEvent



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| response_body | [string](#string) |  |  |
| response_code | [int32](#int32) |  |  |
| time_started | [string](#string) |  | Finalization fields so Aiskhylos doesn&#39;t need extra state |
| time_ended | [string](#string) |  |  |
| total_time_ms | [int64](#int64) |  |  |
| is_closed | [bool](#bool) |  |  |





 


<a name="aristophanes-v1-ItemType"></a>

### ItemType


| Name | Number | Description |
| ---- | ------ | ----------- |
| ITEM_TYPE_UNSPECIFIED | 0 |  |
| TRACE_START | 1 |  |
| TRACE_HOP | 2 |  |
| GRAPHQL | 3 |  |
| ACTION | 4 |  |
| DB_SPAN | 5 |  |
| TRACE_STOP | 6 |  |
| TRACE_HOP_STOP | 7 |  |


 

 


<a name="aristophanes-v1-TraceService"></a>

### TraceService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| HealthCheck | [Empty](#aristophanes-v1-Empty) | [HealthCheckResponse](#aristophanes-v1-HealthCheckResponse) |  |
| Chorus | [ObserveRequest](#aristophanes-v1-ObserveRequest) stream | [ObserveResponse](#aristophanes-v1-ObserveResponse) |  |

 



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

