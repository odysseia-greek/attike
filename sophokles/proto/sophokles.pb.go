// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.1
// source: proto/sophokles.proto

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_sophokles_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_sophokles_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_sophokles_proto_rawDescGZIP(), []int{0}
}

type HealthCheckResponseMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HealthCheckResponseMetrics) Reset() {
	*x = HealthCheckResponseMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_sophokles_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthCheckResponseMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthCheckResponseMetrics) ProtoMessage() {}

func (x *HealthCheckResponseMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_sophokles_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthCheckResponseMetrics.ProtoReflect.Descriptor instead.
func (*HealthCheckResponseMetrics) Descriptor() ([]byte, []int) {
	return file_proto_sophokles_proto_rawDescGZIP(), []int{1}
}

func (x *HealthCheckResponseMetrics) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

type MetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Pod         *PodMetrics `protobuf:"bytes,1,opt,name=pod,proto3" json:"pod,omitempty"`
	CpuUnits    string      `protobuf:"bytes,3,opt,name=cpu_units,json=cpuUnits,proto3" json:"cpu_units,omitempty"`
	MemoryUnits string      `protobuf:"bytes,4,opt,name=memory_units,json=memoryUnits,proto3" json:"memory_units,omitempty"`
}

func (x *MetricsResponse) Reset() {
	*x = MetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_sophokles_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricsResponse) ProtoMessage() {}

func (x *MetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_sophokles_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricsResponse.ProtoReflect.Descriptor instead.
func (*MetricsResponse) Descriptor() ([]byte, []int) {
	return file_proto_sophokles_proto_rawDescGZIP(), []int{2}
}

func (x *MetricsResponse) GetPod() *PodMetrics {
	if x != nil {
		return x.Pod
	}
	return nil
}

func (x *MetricsResponse) GetCpuUnits() string {
	if x != nil {
		return x.CpuUnits
	}
	return ""
}

func (x *MetricsResponse) GetMemoryUnits() string {
	if x != nil {
		return x.MemoryUnits
	}
	return ""
}

type PodMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                string              `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	CpuRaw              int64               `protobuf:"varint,2,opt,name=cpu_raw,json=cpuRaw,proto3" json:"cpu_raw,omitempty"`
	MemoryRaw           int64               `protobuf:"varint,3,opt,name=memory_raw,json=memoryRaw,proto3" json:"memory_raw,omitempty"`
	CpuHumanReadable    string              `protobuf:"bytes,4,opt,name=cpu_human_readable,json=cpuHumanReadable,proto3" json:"cpu_human_readable,omitempty"`
	MemoryHumanReadable string              `protobuf:"bytes,5,opt,name=memory_human_readable,json=memoryHumanReadable,proto3" json:"memory_human_readable,omitempty"`
	Containers          []*ContainerMetrics `protobuf:"bytes,6,rep,name=containers,proto3" json:"containers,omitempty"`
}

func (x *PodMetrics) Reset() {
	*x = PodMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_sophokles_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PodMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PodMetrics) ProtoMessage() {}

func (x *PodMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_sophokles_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PodMetrics.ProtoReflect.Descriptor instead.
func (*PodMetrics) Descriptor() ([]byte, []int) {
	return file_proto_sophokles_proto_rawDescGZIP(), []int{3}
}

func (x *PodMetrics) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PodMetrics) GetCpuRaw() int64 {
	if x != nil {
		return x.CpuRaw
	}
	return 0
}

func (x *PodMetrics) GetMemoryRaw() int64 {
	if x != nil {
		return x.MemoryRaw
	}
	return 0
}

func (x *PodMetrics) GetCpuHumanReadable() string {
	if x != nil {
		return x.CpuHumanReadable
	}
	return ""
}

func (x *PodMetrics) GetMemoryHumanReadable() string {
	if x != nil {
		return x.MemoryHumanReadable
	}
	return ""
}

func (x *PodMetrics) GetContainers() []*ContainerMetrics {
	if x != nil {
		return x.Containers
	}
	return nil
}

type ContainerMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ContainerName                string `protobuf:"bytes,1,opt,name=container_name,json=containerName,proto3" json:"container_name,omitempty"`
	ContainerCpuRaw              int64  `protobuf:"varint,2,opt,name=container_cpu_raw,json=containerCpuRaw,proto3" json:"container_cpu_raw,omitempty"`
	ContainerMemoryRaw           int64  `protobuf:"varint,3,opt,name=container_memory_raw,json=containerMemoryRaw,proto3" json:"container_memory_raw,omitempty"`
	ContainerCpuHumanReadable    string `protobuf:"bytes,4,opt,name=container_cpu_human_readable,json=containerCpuHumanReadable,proto3" json:"container_cpu_human_readable,omitempty"`
	ContainerMemoryHumanReadable string `protobuf:"bytes,5,opt,name=container_memory_human_readable,json=containerMemoryHumanReadable,proto3" json:"container_memory_human_readable,omitempty"`
}

func (x *ContainerMetrics) Reset() {
	*x = ContainerMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_sophokles_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerMetrics) ProtoMessage() {}

func (x *ContainerMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_sophokles_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerMetrics.ProtoReflect.Descriptor instead.
func (*ContainerMetrics) Descriptor() ([]byte, []int) {
	return file_proto_sophokles_proto_rawDescGZIP(), []int{4}
}

func (x *ContainerMetrics) GetContainerName() string {
	if x != nil {
		return x.ContainerName
	}
	return ""
}

func (x *ContainerMetrics) GetContainerCpuRaw() int64 {
	if x != nil {
		return x.ContainerCpuRaw
	}
	return 0
}

func (x *ContainerMetrics) GetContainerMemoryRaw() int64 {
	if x != nil {
		return x.ContainerMemoryRaw
	}
	return 0
}

func (x *ContainerMetrics) GetContainerCpuHumanReadable() string {
	if x != nil {
		return x.ContainerCpuHumanReadable
	}
	return ""
}

func (x *ContainerMetrics) GetContainerMemoryHumanReadable() string {
	if x != nil {
		return x.ContainerMemoryHumanReadable
	}
	return ""
}

var File_proto_sophokles_proto protoreflect.FileDescriptor

var file_proto_sophokles_proto_rawDesc = []byte{
	0x0a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x6f, 0x70, 0x68, 0x6f, 0x6b, 0x6c, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x07,
	0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x34, 0x0a, 0x1a, 0x48, 0x65, 0x61, 0x6c, 0x74,
	0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x76, 0x0a,
	0x0f, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x23, 0x0a, 0x03, 0x70, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x64, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x52, 0x03, 0x70, 0x6f, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x70, 0x75, 0x5f, 0x75, 0x6e, 0x69,
	0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x70, 0x75, 0x55, 0x6e, 0x69,
	0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x75, 0x6e, 0x69,
	0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x55, 0x6e, 0x69, 0x74, 0x73, 0x22, 0xf3, 0x01, 0x0a, 0x0a, 0x50, 0x6f, 0x64, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x63, 0x70, 0x75, 0x5f,
	0x72, 0x61, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x63, 0x70, 0x75, 0x52, 0x61,
	0x77, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x72, 0x61, 0x77, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x61, 0x77,
	0x12, 0x2c, 0x0a, 0x12, 0x63, 0x70, 0x75, 0x5f, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x5f, 0x72, 0x65,
	0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x63, 0x70,
	0x75, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52, 0x65, 0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x32,
	0x0a, 0x15, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x5f, 0x72,
	0x65, 0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x6d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52, 0x65, 0x61, 0x64, 0x61, 0x62,
	0x6c, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52,
	0x0a, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x73, 0x22, 0x9f, 0x02, 0x0a, 0x10,
	0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x74, 0x61,
	0x69, 0x6e, 0x65, 0x72, 0x5f, 0x63, 0x70, 0x75, 0x5f, 0x72, 0x61, 0x77, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x43, 0x70, 0x75,
	0x52, 0x61, 0x77, 0x12, 0x30, 0x0a, 0x14, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x72, 0x61, 0x77, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x12, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x6d, 0x6f,
	0x72, 0x79, 0x52, 0x61, 0x77, 0x12, 0x3f, 0x0a, 0x1c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x5f, 0x63, 0x70, 0x75, 0x5f, 0x68, 0x75, 0x6d, 0x61, 0x6e, 0x5f, 0x72, 0x65, 0x61,
	0x64, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x19, 0x63, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x43, 0x70, 0x75, 0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52, 0x65,
	0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x45, 0x0a, 0x1f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x68, 0x75, 0x6d, 0x61, 0x6e,
	0x5f, 0x72, 0x65, 0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x1c, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x48, 0x75, 0x6d, 0x61, 0x6e, 0x52, 0x65, 0x61, 0x64, 0x61, 0x62, 0x6c, 0x65, 0x32, 0x8d, 0x01,
	0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x45, 0x0a, 0x12, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61,
	0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x34, 0x0a, 0x0c, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x32, 0x5a,
	0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x64, 0x79, 0x73,
	0x73, 0x65, 0x69, 0x61, 0x2d, 0x67, 0x72, 0x65, 0x65, 0x6b, 0x2f, 0x61, 0x74, 0x74, 0x69, 0x6b,
	0x65, 0x2f, 0x73, 0x6f, 0x70, 0x68, 0x6f, 0x6b, 0x6c, 0x65, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_sophokles_proto_rawDescOnce sync.Once
	file_proto_sophokles_proto_rawDescData = file_proto_sophokles_proto_rawDesc
)

func file_proto_sophokles_proto_rawDescGZIP() []byte {
	file_proto_sophokles_proto_rawDescOnce.Do(func() {
		file_proto_sophokles_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_sophokles_proto_rawDescData)
	})
	return file_proto_sophokles_proto_rawDescData
}

var file_proto_sophokles_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_sophokles_proto_goTypes = []interface{}{
	(*Empty)(nil),                      // 0: proto.Empty
	(*HealthCheckResponseMetrics)(nil), // 1: proto.HealthCheckResponseMetrics
	(*MetricsResponse)(nil),            // 2: proto.MetricsResponse
	(*PodMetrics)(nil),                 // 3: proto.PodMetrics
	(*ContainerMetrics)(nil),           // 4: proto.ContainerMetrics
}
var file_proto_sophokles_proto_depIdxs = []int32{
	3, // 0: proto.MetricsResponse.pod:type_name -> proto.PodMetrics
	4, // 1: proto.PodMetrics.containers:type_name -> proto.ContainerMetrics
	0, // 2: proto.MetricsService.HealthCheckMetrics:input_type -> proto.Empty
	0, // 3: proto.MetricsService.FetchMetrics:input_type -> proto.Empty
	1, // 4: proto.MetricsService.HealthCheckMetrics:output_type -> proto.HealthCheckResponseMetrics
	2, // 5: proto.MetricsService.FetchMetrics:output_type -> proto.MetricsResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_sophokles_proto_init() }
func file_proto_sophokles_proto_init() {
	if File_proto_sophokles_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_sophokles_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_proto_sophokles_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HealthCheckResponseMetrics); i {
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
		file_proto_sophokles_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricsResponse); i {
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
		file_proto_sophokles_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PodMetrics); i {
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
		file_proto_sophokles_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerMetrics); i {
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
			RawDescriptor: file_proto_sophokles_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_sophokles_proto_goTypes,
		DependencyIndexes: file_proto_sophokles_proto_depIdxs,
		MessageInfos:      file_proto_sophokles_proto_msgTypes,
	}.Build()
	File_proto_sophokles_proto = out.File
	file_proto_sophokles_proto_rawDesc = nil
	file_proto_sophokles_proto_goTypes = nil
	file_proto_sophokles_proto_depIdxs = nil
}
