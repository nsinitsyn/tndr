// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.2
// source: match_service.proto

package tinderpbv1

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

type GetMatchesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Version       int64                  `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMatchesRequest) Reset() {
	*x = GetMatchesRequest{}
	mi := &file_match_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMatchesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMatchesRequest) ProtoMessage() {}

func (x *GetMatchesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMatchesRequest.ProtoReflect.Descriptor instead.
func (*GetMatchesRequest) Descriptor() ([]byte, []int) {
	return file_match_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetMatchesRequest) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

type GetMatchesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Error:
	//
	//	*GetMatchesResponse_NoUpdates
	//	*GetMatchesResponse_IncorrectVersion
	Error         isGetMatchesResponse_Error `protobuf_oneof:"error"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMatchesResponse) Reset() {
	*x = GetMatchesResponse{}
	mi := &file_match_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMatchesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMatchesResponse) ProtoMessage() {}

func (x *GetMatchesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMatchesResponse.ProtoReflect.Descriptor instead.
func (*GetMatchesResponse) Descriptor() ([]byte, []int) {
	return file_match_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetMatchesResponse) GetError() isGetMatchesResponse_Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *GetMatchesResponse) GetNoUpdates() *NoUpdatesError {
	if x != nil {
		if x, ok := x.Error.(*GetMatchesResponse_NoUpdates); ok {
			return x.NoUpdates
		}
	}
	return nil
}

func (x *GetMatchesResponse) GetIncorrectVersion() *IncorrectVersionError {
	if x != nil {
		if x, ok := x.Error.(*GetMatchesResponse_IncorrectVersion); ok {
			return x.IncorrectVersion
		}
	}
	return nil
}

type isGetMatchesResponse_Error interface {
	isGetMatchesResponse_Error()
}

type GetMatchesResponse_NoUpdates struct {
	NoUpdates *NoUpdatesError `protobuf:"bytes,1,opt,name=no_updates,json=noUpdates,proto3,oneof"`
}

type GetMatchesResponse_IncorrectVersion struct {
	IncorrectVersion *IncorrectVersionError `protobuf:"bytes,2,opt,name=incorrect_version,json=incorrectVersion,proto3,oneof"`
}

func (*GetMatchesResponse_NoUpdates) isGetMatchesResponse_Error() {}

func (*GetMatchesResponse_IncorrectVersion) isGetMatchesResponse_Error() {}

type IncorrectVersionError struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *IncorrectVersionError) Reset() {
	*x = IncorrectVersionError{}
	mi := &file_match_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IncorrectVersionError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncorrectVersionError) ProtoMessage() {}

func (x *IncorrectVersionError) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncorrectVersionError.ProtoReflect.Descriptor instead.
func (*IncorrectVersionError) Descriptor() ([]byte, []int) {
	return file_match_service_proto_rawDescGZIP(), []int{2}
}

type NoUpdatesError struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *NoUpdatesError) Reset() {
	*x = NoUpdatesError{}
	mi := &file_match_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NoUpdatesError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoUpdatesError) ProtoMessage() {}

func (x *NoUpdatesError) ProtoReflect() protoreflect.Message {
	mi := &file_match_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoUpdatesError.ProtoReflect.Descriptor instead.
func (*NoUpdatesError) Descriptor() ([]byte, []int) {
	return file_match_service_proto_rawDescGZIP(), []int{3}
}

var File_match_service_proto protoreflect.FileDescriptor

var file_match_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x2d, 0x0a,
	0x11, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xa4, 0x01, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x6e, 0x6f, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72,
	0x2e, 0x4e, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48,
	0x00, 0x52, 0x09, 0x6e, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x4c, 0x0a, 0x11,
	0x69, 0x6e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72,
	0x2e, 0x49, 0x6e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x48, 0x00, 0x52, 0x10, 0x69, 0x6e, 0x63, 0x6f, 0x72, 0x72,
	0x65, 0x63, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x22, 0x17, 0x0a, 0x15, 0x49, 0x6e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x10, 0x0a, 0x0e,
	0x4e, 0x6f, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x53,
	0x0a, 0x0c, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x12, 0x19, 0x2e, 0x74,
	0x69, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x3b, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72,
	0x70, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_match_service_proto_rawDescOnce sync.Once
	file_match_service_proto_rawDescData = file_match_service_proto_rawDesc
)

func file_match_service_proto_rawDescGZIP() []byte {
	file_match_service_proto_rawDescOnce.Do(func() {
		file_match_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_match_service_proto_rawDescData)
	})
	return file_match_service_proto_rawDescData
}

var file_match_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_match_service_proto_goTypes = []any{
	(*GetMatchesRequest)(nil),     // 0: tinder.GetMatchesRequest
	(*GetMatchesResponse)(nil),    // 1: tinder.GetMatchesResponse
	(*IncorrectVersionError)(nil), // 2: tinder.IncorrectVersionError
	(*NoUpdatesError)(nil),        // 3: tinder.NoUpdatesError
}
var file_match_service_proto_depIdxs = []int32{
	3, // 0: tinder.GetMatchesResponse.no_updates:type_name -> tinder.NoUpdatesError
	2, // 1: tinder.GetMatchesResponse.incorrect_version:type_name -> tinder.IncorrectVersionError
	0, // 2: tinder.MatchService.GetMatches:input_type -> tinder.GetMatchesRequest
	1, // 3: tinder.MatchService.GetMatches:output_type -> tinder.GetMatchesResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_match_service_proto_init() }
func file_match_service_proto_init() {
	if File_match_service_proto != nil {
		return
	}
	file_match_service_proto_msgTypes[1].OneofWrappers = []any{
		(*GetMatchesResponse_NoUpdates)(nil),
		(*GetMatchesResponse_IncorrectVersion)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_match_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_match_service_proto_goTypes,
		DependencyIndexes: file_match_service_proto_depIdxs,
		MessageInfos:      file_match_service_proto_msgTypes,
	}.Build()
	File_match_service_proto = out.File
	file_match_service_proto_rawDesc = nil
	file_match_service_proto_goTypes = nil
	file_match_service_proto_depIdxs = nil
}
