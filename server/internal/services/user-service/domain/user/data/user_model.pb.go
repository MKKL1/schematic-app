// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v5.29.3
// source: user_model.proto

package data

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

type UserModel struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	OidcSub       []byte                 `protobuf:"bytes,3,opt,name=oidc_sub,json=oidcSub,proto3" json:"oidc_sub,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UserModel) Reset() {
	*x = UserModel{}
	mi := &file_user_model_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UserModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserModel) ProtoMessage() {}

func (x *UserModel) ProtoReflect() protoreflect.Message {
	mi := &file_user_model_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserModel.ProtoReflect.Descriptor instead.
func (*UserModel) Descriptor() ([]byte, []int) {
	return file_user_model_proto_rawDescGZIP(), []int{0}
}

func (x *UserModel) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserModel) GetOidcSub() []byte {
	if x != nil {
		return x.OidcSub
	}
	return nil
}

var File_user_model_proto protoreflect.FileDescriptor

var file_user_model_proto_rawDesc = []byte{
	0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x4a, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x69, 0x64,
	0x63, 0x5f, 0x73, 0x75, 0x62, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x6f, 0x69, 0x64,
	0x63, 0x53, 0x75, 0x62, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x3b, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_model_proto_rawDescOnce sync.Once
	file_user_model_proto_rawDescData = file_user_model_proto_rawDesc
)

func file_user_model_proto_rawDescGZIP() []byte {
	file_user_model_proto_rawDescOnce.Do(func() {
		file_user_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_model_proto_rawDescData)
	})
	return file_user_model_proto_rawDescData
}

var file_user_model_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_user_model_proto_goTypes = []any{
	(*UserModel)(nil), // 0: data.UserModel
}
var file_user_model_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_user_model_proto_init() }
func file_user_model_proto_init() {
	if File_user_model_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_model_proto_goTypes,
		DependencyIndexes: file_user_model_proto_depIdxs,
		MessageInfos:      file_user_model_proto_msgTypes,
	}.Build()
	File_user_model_proto = out.File
	file_user_model_proto_rawDesc = nil
	file_user_model_proto_goTypes = nil
	file_user_model_proto_depIdxs = nil
}
