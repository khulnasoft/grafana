// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: secret.proto

package access

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

type SecretRequest struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	Namespace string                 `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	// A list of the names we want decoded
	Names         []string `protobuf:"bytes,2,rep,name=names,proto3" json:"names,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecretRequest) Reset() {
	*x = SecretRequest{}
	mi := &file_secret_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretRequest) ProtoMessage() {}

func (x *SecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_secret_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretRequest.ProtoReflect.Descriptor instead.
func (*SecretRequest) Descriptor() ([]byte, []int) {
	return file_secret_proto_rawDescGZIP(), []int{0}
}

func (x *SecretRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *SecretRequest) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

type SecretResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The value will be empty if the user can not read it
	// This could be that it exists but the request does not have permissions
	// OR that the values does not exist.
	// We avoid 404 so you can not scan for valid values
	Values        map[string]string `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecretResponse) Reset() {
	*x = SecretResponse{}
	mi := &file_secret_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretResponse) ProtoMessage() {}

func (x *SecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_secret_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecretResponse.ProtoReflect.Descriptor instead.
func (*SecretResponse) Descriptor() ([]byte, []int) {
	return file_secret_proto_rawDescGZIP(), []int{1}
}

func (x *SecretResponse) GetValues() map[string]string {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_secret_proto protoreflect.FileDescriptor

var file_secret_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x43, 0x0a, 0x0d, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x87, 0x01, 0x0a, 0x0e,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x39, 0x0a, 0x0b, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x51, 0x0a, 0x0d, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x44,
	0x65, 0x63, 0x6f, 0x64, 0x65, 0x72, 0x12, 0x40, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x65, 0x63,
	0x72, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x15, 0x2e, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x61, 0x66, 0x61, 0x6e, 0x61, 0x2f, 0x67,
	0x72, 0x61, 0x66, 0x61, 0x6e, 0x61, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_secret_proto_rawDescOnce sync.Once
	file_secret_proto_rawDescData = file_secret_proto_rawDesc
)

func file_secret_proto_rawDescGZIP() []byte {
	file_secret_proto_rawDescOnce.Do(func() {
		file_secret_proto_rawDescData = protoimpl.X.CompressGZIP(file_secret_proto_rawDescData)
	})
	return file_secret_proto_rawDescData
}

var file_secret_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_secret_proto_goTypes = []any{
	(*SecretRequest)(nil),  // 0: access.SecretRequest
	(*SecretResponse)(nil), // 1: access.SecretResponse
	nil,                    // 2: access.SecretResponse.ValuesEntry
}
var file_secret_proto_depIdxs = []int32{
	2, // 0: access.SecretResponse.values:type_name -> access.SecretResponse.ValuesEntry
	0, // 1: access.SecretDecoder.GetSecretValues:input_type -> access.SecretRequest
	1, // 2: access.SecretDecoder.GetSecretValues:output_type -> access.SecretResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_secret_proto_init() }
func file_secret_proto_init() {
	if File_secret_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_secret_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_secret_proto_goTypes,
		DependencyIndexes: file_secret_proto_depIdxs,
		MessageInfos:      file_secret_proto_msgTypes,
	}.Build()
	File_secret_proto = out.File
	file_secret_proto_rawDesc = nil
	file_secret_proto_goTypes = nil
	file_secret_proto_depIdxs = nil
}
