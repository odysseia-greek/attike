// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/aristophanes.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Trace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpanId       string `protobuf:"bytes,1,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	Method       string `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	Url          string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	Timestamp    string `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	PodName      string `protobuf:"bytes,6,opt,name=pod_name,json=podName,proto3" json:"pod_name,omitempty"`
	Namespace    string `protobuf:"bytes,7,opt,name=namespace,proto3" json:"namespace,omitempty"`
	ItemType     string `protobuf:"bytes,8,opt,name=item_type,json=itemType,proto3" json:"item_type,omitempty"`
}

func (x *Trace) Reset() {
	*x = Trace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trace) ProtoMessage() {}

func (x *Trace) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trace.ProtoReflect.Descriptor instead.
func (*Trace) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{0}
}

func (x *Trace) GetSpanId() string {
	if x != nil {
		return x.SpanId
	}
	return ""
}

func (x *Trace) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *Trace) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *Trace) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Trace) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Trace) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

func (x *Trace) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Trace) GetItemType() string {
	if x != nil {
		return x.ItemType
	}
	return ""
}

type Span struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpanId       string `protobuf:"bytes,1,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	Timestamp    string `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	PodName      string `protobuf:"bytes,4,opt,name=pod_name,json=podName,proto3" json:"pod_name,omitempty"`
	Namespace    string `protobuf:"bytes,5,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Action       string `protobuf:"bytes,6,opt,name=action,proto3" json:"action,omitempty"`                                 // Action performed in the span
	RequestBody  string `protobuf:"bytes,7,opt,name=request_body,json=requestBody,proto3" json:"request_body,omitempty"`    // Optional: Request body data
	ResponseBody string `protobuf:"bytes,8,opt,name=response_body,json=responseBody,proto3" json:"response_body,omitempty"` // Optional: Response body data
	ItemType     string `protobuf:"bytes,9,opt,name=item_type,json=itemType,proto3" json:"item_type,omitempty"`
}

func (x *Span) Reset() {
	*x = Span{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Span) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Span) ProtoMessage() {}

func (x *Span) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Span.ProtoReflect.Descriptor instead.
func (*Span) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{1}
}

func (x *Span) GetSpanId() string {
	if x != nil {
		return x.SpanId
	}
	return ""
}

func (x *Span) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *Span) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Span) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

func (x *Span) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *Span) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *Span) GetRequestBody() string {
	if x != nil {
		return x.RequestBody
	}
	return ""
}

func (x *Span) GetResponseBody() string {
	if x != nil {
		return x.ResponseBody
	}
	return ""
}

func (x *Span) GetItemType() string {
	if x != nil {
		return x.ItemType
	}
	return ""
}

type DatabaseSpan struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpanId       string `protobuf:"bytes,1,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	Timestamp    string `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	PodName      string `protobuf:"bytes,4,opt,name=pod_name,json=podName,proto3" json:"pod_name,omitempty"`
	Namespace    string `protobuf:"bytes,5,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Query        string `protobuf:"bytes,6,opt,name=query,proto3" json:"query,omitempty"`                             // Database query statement
	ResultJson   string `protobuf:"bytes,7,opt,name=result_json,json=resultJson,proto3" json:"result_json,omitempty"` // Query result data as JSON string
	ItemType     string `protobuf:"bytes,8,opt,name=item_type,json=itemType,proto3" json:"item_type,omitempty"`
}

func (x *DatabaseSpan) Reset() {
	*x = DatabaseSpan{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseSpan) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseSpan) ProtoMessage() {}

func (x *DatabaseSpan) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseSpan.ProtoReflect.Descriptor instead.
func (*DatabaseSpan) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{2}
}

func (x *DatabaseSpan) GetSpanId() string {
	if x != nil {
		return x.SpanId
	}
	return ""
}

func (x *DatabaseSpan) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *DatabaseSpan) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *DatabaseSpan) GetPodName() string {
	if x != nil {
		return x.PodName
	}
	return ""
}

func (x *DatabaseSpan) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *DatabaseSpan) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *DatabaseSpan) GetResultJson() string {
	if x != nil {
		return x.ResultJson
	}
	return ""
}

func (x *DatabaseSpan) GetItemType() string {
	if x != nil {
		return x.ItemType
	}
	return ""
}

type StartTraceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method    string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Url       string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	SaveTrace bool   `protobuf:"varint,3,opt,name=save_trace,json=saveTrace,proto3" json:"save_trace,omitempty"` // Indicates whether the trace should be saved
}

func (x *StartTraceRequest) Reset() {
	*x = StartTraceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartTraceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartTraceRequest) ProtoMessage() {}

func (x *StartTraceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartTraceRequest.ProtoReflect.Descriptor instead.
func (*StartTraceRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{3}
}

func (x *StartTraceRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *StartTraceRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *StartTraceRequest) GetSaveTrace() bool {
	if x != nil {
		return x.SaveTrace
	}
	return false
}

type StartSpanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraceId string `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
}

func (x *StartSpanRequest) Reset() {
	*x = StartSpanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartSpanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartSpanRequest) ProtoMessage() {}

func (x *StartSpanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartSpanRequest.ProtoReflect.Descriptor instead.
func (*StartSpanRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{4}
}

func (x *StartSpanRequest) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

type SpanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraceId      string `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	SaveTrace    bool   `protobuf:"varint,3,opt,name=save_trace,json=saveTrace,proto3" json:"save_trace,omitempty"`         // Indicates whether the trace should be saved
	Action       string `protobuf:"bytes,4,opt,name=action,proto3" json:"action,omitempty"`                                 // Action performed in the span
	RequestBody  string `protobuf:"bytes,5,opt,name=request_body,json=requestBody,proto3" json:"request_body,omitempty"`    // Optional: Request body data
	ResponseBody string `protobuf:"bytes,6,opt,name=response_body,json=responseBody,proto3" json:"response_body,omitempty"` // Optional: Response body data
}

func (x *SpanRequest) Reset() {
	*x = SpanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SpanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SpanRequest) ProtoMessage() {}

func (x *SpanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SpanRequest.ProtoReflect.Descriptor instead.
func (*SpanRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{5}
}

func (x *SpanRequest) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *SpanRequest) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *SpanRequest) GetSaveTrace() bool {
	if x != nil {
		return x.SaveTrace
	}
	return false
}

func (x *SpanRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *SpanRequest) GetRequestBody() string {
	if x != nil {
		return x.RequestBody
	}
	return ""
}

func (x *SpanRequest) GetResponseBody() string {
	if x != nil {
		return x.ResponseBody
	}
	return ""
}

type DatabaseSpanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraceId      string `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	SaveTrace    bool   `protobuf:"varint,3,opt,name=save_trace,json=saveTrace,proto3" json:"save_trace,omitempty"`   // Indicates whether the trace should be saved
	Action       string `protobuf:"bytes,4,opt,name=action,proto3" json:"action,omitempty"`                           // Action performed in the span
	Query        string `protobuf:"bytes,6,opt,name=query,proto3" json:"query,omitempty"`                             // Database query statement
	ResultJson   string `protobuf:"bytes,7,opt,name=result_json,json=resultJson,proto3" json:"result_json,omitempty"` // Query result data as JSON string
}

func (x *DatabaseSpanRequest) Reset() {
	*x = DatabaseSpanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatabaseSpanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatabaseSpanRequest) ProtoMessage() {}

func (x *DatabaseSpanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatabaseSpanRequest.ProtoReflect.Descriptor instead.
func (*DatabaseSpanRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{6}
}

func (x *DatabaseSpanRequest) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *DatabaseSpanRequest) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *DatabaseSpanRequest) GetSaveTrace() bool {
	if x != nil {
		return x.SaveTrace
	}
	return false
}

func (x *DatabaseSpanRequest) GetAction() string {
	if x != nil {
		return x.Action
	}
	return ""
}

func (x *DatabaseSpanRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *DatabaseSpanRequest) GetResultJson() string {
	if x != nil {
		return x.ResultJson
	}
	return ""
}

type TraceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CombinedId string `protobuf:"bytes,1,opt,name=combined_id,json=combinedId,proto3" json:"combined_id,omitempty"`
}

func (x *TraceResponse) Reset() {
	*x = TraceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TraceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceResponse) ProtoMessage() {}

func (x *TraceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceResponse.ProtoReflect.Descriptor instead.
func (*TraceResponse) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{7}
}

func (x *TraceResponse) GetCombinedId() string {
	if x != nil {
		return x.CombinedId
	}
	return ""
}

type TraceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TraceId      string `protobuf:"bytes,1,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	ParentSpanId string `protobuf:"bytes,2,opt,name=parent_span_id,json=parentSpanId,proto3" json:"parent_span_id,omitempty"`
	Method       string `protobuf:"bytes,3,opt,name=method,proto3" json:"method,omitempty"`
	Url          string `protobuf:"bytes,4,opt,name=url,proto3" json:"url,omitempty"`
	SaveTrace    bool   `protobuf:"varint,5,opt,name=save_trace,json=saveTrace,proto3" json:"save_trace,omitempty"` // Indicates whether the trace should be saved
}

func (x *TraceRequest) Reset() {
	*x = TraceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_aristophanes_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TraceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TraceRequest) ProtoMessage() {}

func (x *TraceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_aristophanes_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TraceRequest.ProtoReflect.Descriptor instead.
func (*TraceRequest) Descriptor() ([]byte, []int) {
	return file_proto_aristophanes_proto_rawDescGZIP(), []int{8}
}

func (x *TraceRequest) GetTraceId() string {
	if x != nil {
		return x.TraceId
	}
	return ""
}

func (x *TraceRequest) GetParentSpanId() string {
	if x != nil {
		return x.ParentSpanId
	}
	return ""
}

func (x *TraceRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *TraceRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *TraceRequest) GetSaveTrace() bool {
	if x != nil {
		return x.SaveTrace
	}
	return false
}

var File_proto_aristophanes_proto protoreflect.FileDescriptor

var file_proto_aristophanes_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74, 0x6f, 0x70, 0x68,
	0x61, 0x6e, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xe4, 0x01, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x73,
	0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x70,
	0x61, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73,
	0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x75, 0x72, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69,
	0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x22, 0x99, 0x02, 0x0a, 0x04, 0x53, 0x70, 0x61,
	0x6e, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x61, 0x6e, 0x49, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x19,
	0x0a, 0x08, 0x70, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x6f,
	0x64, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x22, 0xf8, 0x01, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73,
	0x65, 0x53, 0x70, 0x61, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x24,
	0x0a, 0x0e, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x70,
	0x61, 0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6a, 0x73, 0x6f, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4a, 0x73,
	0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x22,
	0x5c, 0x0a, 0x11, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1d,
	0x0a, 0x0a, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x09, 0x73, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x63, 0x65, 0x22, 0x2d, 0x0a,
	0x10, 0x53, 0x74, 0x61, 0x72, 0x74, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x22, 0xcd, 0x01, 0x0a,
	0x0b, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x53, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x09, 0x73, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x5f, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x22, 0xc4, 0x01, 0x0a,
	0x13, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12,
	0x24, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x70, 0x61, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x53,
	0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x74, 0x72,
	0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x73, 0x61, 0x76, 0x65, 0x54,
	0x72, 0x61, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05,
	0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x6a, 0x73, 0x6f,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4a,
	0x73, 0x6f, 0x6e, 0x22, 0x30, 0x0a, 0x0d, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x62, 0x69, 0x6e, 0x65, 0x64,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6d, 0x62, 0x69,
	0x6e, 0x65, 0x64, 0x49, 0x64, 0x22, 0x98, 0x01, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x24, 0x0a, 0x0e, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x70, 0x61, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x61, 0x72, 0x65, 0x6e,
	0x74, 0x53, 0x70, 0x61, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x73, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x63, 0x65,
	0x32, 0xb3, 0x02, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x63, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3c, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12,
	0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x72, 0x61,
	0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x63, 0x65, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4e, 0x65, 0x77, 0x53,
	0x70, 0x61, 0x6e, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x53, 0x70, 0x61, 0x6e, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65,
	0x53, 0x70, 0x61, 0x6e, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74,
	0x61, 0x62, 0x61, 0x73, 0x65, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x64, 0x79, 0x73, 0x73, 0x65, 0x69, 0x61, 0x2d, 0x67, 0x72,
	0x65, 0x65, 0x6b, 0x2f, 0x61, 0x74, 0x74, 0x69, 0x6b, 0x65, 0x2f, 0x61, 0x72, 0x69, 0x73, 0x74,
	0x6f, 0x70, 0x68, 0x61, 0x6e, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_aristophanes_proto_rawDescOnce sync.Once
	file_proto_aristophanes_proto_rawDescData = file_proto_aristophanes_proto_rawDesc
)

func file_proto_aristophanes_proto_rawDescGZIP() []byte {
	file_proto_aristophanes_proto_rawDescOnce.Do(func() {
		file_proto_aristophanes_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_aristophanes_proto_rawDescData)
	})
	return file_proto_aristophanes_proto_rawDescData
}

var file_proto_aristophanes_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_aristophanes_proto_goTypes = []interface{}{
	(*Trace)(nil),               // 0: proto.Trace
	(*Span)(nil),                // 1: proto.Span
	(*DatabaseSpan)(nil),        // 2: proto.DatabaseSpan
	(*StartTraceRequest)(nil),   // 3: proto.StartTraceRequest
	(*StartSpanRequest)(nil),    // 4: proto.StartSpanRequest
	(*SpanRequest)(nil),         // 5: proto.SpanRequest
	(*DatabaseSpanRequest)(nil), // 6: proto.DatabaseSpanRequest
	(*TraceResponse)(nil),       // 7: proto.TraceResponse
	(*TraceRequest)(nil),        // 8: proto.TraceRequest
}
var file_proto_aristophanes_proto_depIdxs = []int32{
	3, // 0: proto.TraceService.StartTrace:input_type -> proto.StartTraceRequest
	8, // 1: proto.TraceService.Trace:input_type -> proto.TraceRequest
	4, // 2: proto.TraceService.StartNewSpan:input_type -> proto.StartSpanRequest
	5, // 3: proto.TraceService.Span:input_type -> proto.SpanRequest
	6, // 4: proto.TraceService.DatabaseSpan:input_type -> proto.DatabaseSpanRequest
	7, // 5: proto.TraceService.StartTrace:output_type -> proto.TraceResponse
	7, // 6: proto.TraceService.Trace:output_type -> proto.TraceResponse
	7, // 7: proto.TraceService.StartNewSpan:output_type -> proto.TraceResponse
	7, // 8: proto.TraceService.Span:output_type -> proto.TraceResponse
	7, // 9: proto.TraceService.DatabaseSpan:output_type -> proto.TraceResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_aristophanes_proto_init() }
func file_proto_aristophanes_proto_init() {
	if File_proto_aristophanes_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_aristophanes_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trace); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Span); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseSpan); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartTraceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartSpanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SpanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatabaseSpanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TraceResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_aristophanes_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TraceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_aristophanes_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_aristophanes_proto_goTypes,
		DependencyIndexes: file_proto_aristophanes_proto_depIdxs,
		MessageInfos:      file_proto_aristophanes_proto_msgTypes,
	}.Build()
	File_proto_aristophanes_proto = out.File
	file_proto_aristophanes_proto_rawDesc = nil
	file_proto_aristophanes_proto_goTypes = nil
	file_proto_aristophanes_proto_depIdxs = nil
}
