// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v5.29.2
// source: geo_service.proto

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

type GetProfilesByLocationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Latitude      float64                `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude     float64                `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProfilesByLocationRequest) Reset() {
	*x = GetProfilesByLocationRequest{}
	mi := &file_geo_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProfilesByLocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProfilesByLocationRequest) ProtoMessage() {}

func (x *GetProfilesByLocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_geo_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProfilesByLocationRequest.ProtoReflect.Descriptor instead.
func (*GetProfilesByLocationRequest) Descriptor() ([]byte, []int) {
	return file_geo_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetProfilesByLocationRequest) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *GetProfilesByLocationRequest) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type GetProfilesByLocationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Profiles      []*LocationProfileDto  `protobuf:"bytes,1,rep,name=profiles,proto3" json:"profiles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProfilesByLocationResponse) Reset() {
	*x = GetProfilesByLocationResponse{}
	mi := &file_geo_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProfilesByLocationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProfilesByLocationResponse) ProtoMessage() {}

func (x *GetProfilesByLocationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_geo_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProfilesByLocationResponse.ProtoReflect.Descriptor instead.
func (*GetProfilesByLocationResponse) Descriptor() ([]byte, []int) {
	return file_geo_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetProfilesByLocationResponse) GetProfiles() []*LocationProfileDto {
	if x != nil {
		return x.Profiles
	}
	return nil
}

type LocationProfileDto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProfileId     int64                  `protobuf:"varint,1,opt,name=profile_id,json=profileId,proto3" json:"profile_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	PhotoUrls     []string               `protobuf:"bytes,4,rep,name=photo_urls,json=photoUrls,proto3" json:"photo_urls,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LocationProfileDto) Reset() {
	*x = LocationProfileDto{}
	mi := &file_geo_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LocationProfileDto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocationProfileDto) ProtoMessage() {}

func (x *LocationProfileDto) ProtoReflect() protoreflect.Message {
	mi := &file_geo_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocationProfileDto.ProtoReflect.Descriptor instead.
func (*LocationProfileDto) Descriptor() ([]byte, []int) {
	return file_geo_service_proto_rawDescGZIP(), []int{2}
}

func (x *LocationProfileDto) GetProfileId() int64 {
	if x != nil {
		return x.ProfileId
	}
	return 0
}

func (x *LocationProfileDto) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LocationProfileDto) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *LocationProfileDto) GetPhotoUrls() []string {
	if x != nil {
		return x.PhotoUrls
	}
	return nil
}

type ChangeLocationRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Latitude      float64                `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude     float64                `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangeLocationRequest) Reset() {
	*x = ChangeLocationRequest{}
	mi := &file_geo_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangeLocationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeLocationRequest) ProtoMessage() {}

func (x *ChangeLocationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_geo_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeLocationRequest.ProtoReflect.Descriptor instead.
func (*ChangeLocationRequest) Descriptor() ([]byte, []int) {
	return file_geo_service_proto_rawDescGZIP(), []int{3}
}

func (x *ChangeLocationRequest) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *ChangeLocationRequest) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

type ChangeLocationResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChangeLocationResponse) Reset() {
	*x = ChangeLocationResponse{}
	mi := &file_geo_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChangeLocationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeLocationResponse) ProtoMessage() {}

func (x *ChangeLocationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_geo_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeLocationResponse.ProtoReflect.Descriptor instead.
func (*ChangeLocationResponse) Descriptor() ([]byte, []int) {
	return file_geo_service_proto_rawDescGZIP(), []int{4}
}

var File_geo_service_proto protoreflect.FileDescriptor

var file_geo_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x67, 0x65, 0x6f, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x22, 0x58, 0x0a, 0x1c, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c,
	0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69,
	0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67,
	0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x57, 0x0a, 0x1d, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65,
	0x72, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x44, 0x74, 0x6f, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x88,
	0x01, 0x0a, 0x12, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x44, 0x74, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x68,
	0x6f, 0x74, 0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x68, 0x6f, 0x74, 0x6f, 0x55, 0x72, 0x6c, 0x73, 0x22, 0x51, 0x0a, 0x15, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x22, 0x18, 0x0a, 0x16,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc3, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x6f, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x64, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x24,
	0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x42, 0x79, 0x4c, 0x6f, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x43,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e,
	0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x74,
	0x69, 0x6e, 0x64, 0x65, 0x72, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0f, 0x5a, 0x0d,
	0x2e, 0x2f, 0x3b, 0x74, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x70, 0x62, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_geo_service_proto_rawDescOnce sync.Once
	file_geo_service_proto_rawDescData = file_geo_service_proto_rawDesc
)

func file_geo_service_proto_rawDescGZIP() []byte {
	file_geo_service_proto_rawDescOnce.Do(func() {
		file_geo_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_geo_service_proto_rawDescData)
	})
	return file_geo_service_proto_rawDescData
}

var file_geo_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_geo_service_proto_goTypes = []any{
	(*GetProfilesByLocationRequest)(nil),  // 0: tinder.GetProfilesByLocationRequest
	(*GetProfilesByLocationResponse)(nil), // 1: tinder.GetProfilesByLocationResponse
	(*LocationProfileDto)(nil),            // 2: tinder.LocationProfileDto
	(*ChangeLocationRequest)(nil),         // 3: tinder.ChangeLocationRequest
	(*ChangeLocationResponse)(nil),        // 4: tinder.ChangeLocationResponse
}
var file_geo_service_proto_depIdxs = []int32{
	2, // 0: tinder.GetProfilesByLocationResponse.profiles:type_name -> tinder.LocationProfileDto
	0, // 1: tinder.GeoService.GetProfilesByLocation:input_type -> tinder.GetProfilesByLocationRequest
	3, // 2: tinder.GeoService.ChangeLocation:input_type -> tinder.ChangeLocationRequest
	1, // 3: tinder.GeoService.GetProfilesByLocation:output_type -> tinder.GetProfilesByLocationResponse
	4, // 4: tinder.GeoService.ChangeLocation:output_type -> tinder.ChangeLocationResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_geo_service_proto_init() }
func file_geo_service_proto_init() {
	if File_geo_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_geo_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_geo_service_proto_goTypes,
		DependencyIndexes: file_geo_service_proto_depIdxs,
		MessageInfos:      file_geo_service_proto_msgTypes,
	}.Build()
	File_geo_service_proto = out.File
	file_geo_service_proto_rawDesc = nil
	file_geo_service_proto_goTypes = nil
	file_geo_service_proto_depIdxs = nil
}
