// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: post_entity.proto

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

type PostEntity struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Desc          *string                `protobuf:"bytes,3,opt,name=desc,proto3,oneof" json:"desc,omitempty"`
	Owner         int64                  `protobuf:"varint,4,opt,name=owner,proto3" json:"owner,omitempty"`
	AName         *string                `protobuf:"bytes,5,opt,name=aName,proto3,oneof" json:"aName,omitempty"`
	AId           *int64                 `protobuf:"varint,6,opt,name=aId,proto3,oneof" json:"aId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PostEntity) Reset() {
	*x = PostEntity{}
	mi := &file_post_entity_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostEntity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostEntity) ProtoMessage() {}

func (x *PostEntity) ProtoReflect() protoreflect.Message {
	mi := &file_post_entity_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostEntity.ProtoReflect.Descriptor instead.
func (*PostEntity) Descriptor() ([]byte, []int) {
	return file_post_entity_proto_rawDescGZIP(), []int{0}
}

func (x *PostEntity) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PostEntity) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PostEntity) GetDesc() string {
	if x != nil && x.Desc != nil {
		return *x.Desc
	}
	return ""
}

func (x *PostEntity) GetOwner() int64 {
	if x != nil {
		return x.Owner
	}
	return 0
}

func (x *PostEntity) GetAName() string {
	if x != nil && x.AName != nil {
		return *x.AName
	}
	return ""
}

func (x *PostEntity) GetAId() int64 {
	if x != nil && x.AId != nil {
		return *x.AId
	}
	return 0
}

var File_post_entity_proto protoreflect.FileDescriptor

var file_post_entity_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x01, 0x0a, 0x0a, 0x50,
	0x6f, 0x73, 0x74, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a,
	0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x88, 0x01, 0x01, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x19, 0x0a, 0x05,
	0x61, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x05, 0x61,
	0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x61, 0x49, 0x64, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x03, 0x61, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x07,
	0x0a, 0x05, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x61, 0x4e, 0x61, 0x6d,
	0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x61, 0x49, 0x64, 0x42, 0x58, 0x5a, 0x56, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x4b, 0x4b, 0x4c, 0x31, 0x2f, 0x73, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x73, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_post_entity_proto_rawDescOnce sync.Once
	file_post_entity_proto_rawDescData = file_post_entity_proto_rawDesc
)

func file_post_entity_proto_rawDescGZIP() []byte {
	file_post_entity_proto_rawDescOnce.Do(func() {
		file_post_entity_proto_rawDescData = protoimpl.X.CompressGZIP(file_post_entity_proto_rawDescData)
	})
	return file_post_entity_proto_rawDescData
}

var file_post_entity_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_post_entity_proto_goTypes = []any{
	(*PostEntity)(nil), // 0: proto.PostEntity
}
var file_post_entity_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_post_entity_proto_init() }
func file_post_entity_proto_init() {
	if File_post_entity_proto != nil {
		return
	}
	file_post_entity_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_post_entity_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_post_entity_proto_goTypes,
		DependencyIndexes: file_post_entity_proto_depIdxs,
		MessageInfos:      file_post_entity_proto_msgTypes,
	}.Build()
	File_post_entity_proto = out.File
	file_post_entity_proto_rawDesc = nil
	file_post_entity_proto_goTypes = nil
	file_post_entity_proto_depIdxs = nil
}
