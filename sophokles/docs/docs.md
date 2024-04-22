# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [sophokles.proto](#sophokles-proto)
    - [ContainerMetrics](#proto-ContainerMetrics)
    - [Empty](#proto-Empty)
    - [HealthCheckResponse](#proto-HealthCheckResponse)
    - [MetricsResponse](#proto-MetricsResponse)
    - [PodMetrics](#proto-PodMetrics)
  
    - [MetricsService](#proto-MetricsService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="sophokles-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## sophokles.proto



<a name="proto-ContainerMetrics"></a>

### ContainerMetrics



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| container_name | [string](#string) |  |  |
| container_cpu_raw | [int64](#int64) |  |  |
| container_memory_raw | [int64](#int64) |  |  |
| container_cpu_human_readable | [string](#string) |  |  |
| container_memory_human_readable | [string](#string) |  |  |






<a name="proto-Empty"></a>

### Empty







<a name="proto-HealthCheckResponse"></a>

### HealthCheckResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [bool](#bool) |  |  |






<a name="proto-MetricsResponse"></a>

### MetricsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pod | [PodMetrics](#proto-PodMetrics) |  |  |
| cpu_units | [string](#string) |  |  |
| memory_units | [string](#string) |  |  |






<a name="proto-PodMetrics"></a>

### PodMetrics



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| cpu_raw | [int64](#int64) |  |  |
| memory_raw | [int64](#int64) |  |  |
| cpu_human_readable | [string](#string) |  |  |
| memory_human_readable | [string](#string) |  |  |
| containers | [ContainerMetrics](#proto-ContainerMetrics) | repeated |  |





 

 

 


<a name="proto-MetricsService"></a>

### MetricsService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| HealthCheck | [Empty](#proto-Empty) | [HealthCheckResponse](#proto-HealthCheckResponse) |  |
| FetchMetrics | [Empty](#proto-Empty) | [MetricsResponse](#proto-MetricsResponse) |  |

 



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

