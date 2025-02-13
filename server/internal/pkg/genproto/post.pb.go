// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v4.25.6
// source: api/proto/post.proto

package genproto

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

type PostByIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PostByIdRequest) Reset() {
	*x = PostByIdRequest{}
	mi := &file_api_proto_post_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostByIdRequest) ProtoMessage() {}

func (x *PostByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostByIdRequest.ProtoReflect.Descriptor instead.
func (*PostByIdRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{0}
}

func (x *PostByIdRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type CreatePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description   *string                `protobuf:"bytes,2,opt,name=description,proto3,oneof" json:"description,omitempty"`
	AuthorName    *string                `protobuf:"bytes,3,opt,name=authorName,proto3,oneof" json:"authorName,omitempty"`
	AuthorId      *int64                 `protobuf:"varint,4,opt,name=authorId,proto3,oneof" json:"authorId,omitempty"`
	AuthSub       []byte                 `protobuf:"bytes,5,opt,name=authSub,proto3" json:"authSub,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePostRequest) Reset() {
	*x = CreatePostRequest{}
	mi := &file_api_proto_post_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostRequest) ProtoMessage() {}

func (x *CreatePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostRequest.ProtoReflect.Descriptor instead.
func (*CreatePostRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePostRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreatePostRequest) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *CreatePostRequest) GetAuthorName() string {
	if x != nil && x.AuthorName != nil {
		return *x.AuthorName
	}
	return ""
}

func (x *CreatePostRequest) GetAuthorId() int64 {
	if x != nil && x.AuthorId != nil {
		return *x.AuthorId
	}
	return 0
}

func (x *CreatePostRequest) GetAuthSub() []byte {
	if x != nil {
		return x.AuthSub
	}
	return nil
}

type CreatePostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePostResponse) Reset() {
	*x = CreatePostResponse{}
	mi := &file_api_proto_post_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostResponse) ProtoMessage() {}

func (x *CreatePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostResponse.ProtoReflect.Descriptor instead.
func (*CreatePostResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{2}
}

func (x *CreatePostResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type CategoryVars struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Metadata      []byte                 `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CategoryVars) Reset() {
	*x = CategoryVars{}
	mi := &file_api_proto_post_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CategoryVars) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategoryVars) ProtoMessage() {}

func (x *CategoryVars) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CategoryVars.ProtoReflect.Descriptor instead.
func (*CategoryVars) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{3}
}

func (x *CategoryVars) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CategoryVars) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tag           string                 `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Tag) Reset() {
	*x = Tag{}
	mi := &file_api_proto_post_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{4}
}

func (x *Tag) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

type Post struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   *string                `protobuf:"bytes,3,opt,name=description,proto3,oneof" json:"description,omitempty"`
	Owner         int64                  `protobuf:"varint,4,opt,name=owner,proto3" json:"owner,omitempty"`
	Author        *int64                 `protobuf:"varint,5,opt,name=author,proto3,oneof" json:"author,omitempty"`
	Vars          []*CategoryVars        `protobuf:"bytes,6,rep,name=vars,proto3" json:"vars,omitempty"`
	Tags          []*Tag                 `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Post) Reset() {
	*x = Post{}
	mi := &file_api_proto_post_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_post_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_api_proto_post_proto_rawDescGZIP(), []int{5}
}

func (x *Post) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Post) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Post) GetDescription() string {
	if x != nil && x.Description != nil {
		return *x.Description
	}
	return ""
}

func (x *Post) GetOwner() int64 {
	if x != nil {
		return x.Owner
	}
	return 0
}

func (x *Post) GetAuthor() int64 {
	if x != nil && x.Author != nil {
		return *x.Author
	}
	return 0
}

func (x *Post) GetVars() []*CategoryVars {
	if x != nil {
		return x.Vars
	}
	return nil
}

func (x *Post) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

var File_api_proto_post_proto protoreflect.FileDescriptor

var file_api_proto_post_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a,
	0x0f, 0x50, 0x6f, 0x73, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64,
	0x22, 0xda, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x00, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01,
	0x01, 0x12, 0x23, 0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x1f, 0x0a, 0x08, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x08, 0x61, 0x75, 0x74, 0x68,
	0x6f, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x53,
	0x75, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x53, 0x75,
	0x62, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x4e, 0x61, 0x6d, 0x65,
	0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x49, 0x64, 0x22, 0x24, 0x0a,
	0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x0c, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x56,
	0x61, 0x72, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x17, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0xe8, 0x01, 0x0a,
	0x04, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01,
	0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x01, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x04, 0x76, 0x61, 0x72, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x56, 0x61, 0x72, 0x73, 0x52, 0x04, 0x76, 0x61, 0x72, 0x73, 0x12, 0x1e, 0x0a, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x42, 0x0e, 0x0a, 0x0c,
	0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x09, 0x0a, 0x07,
	0x5f, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x32, 0x88, 0x01, 0x0a, 0x0b, 0x50, 0x6f, 0x73, 0x74,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x6f,
	0x73, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50,
	0x6f, 0x73, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x22, 0x00, 0x12, 0x43, 0x0a,
	0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x18, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x4d, 0x4b, 0x4b, 0x4c, 0x31, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x74, 0x69, 0x63,
	0x2d, 0x61, 0x70, 0x70, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_proto_post_proto_rawDescOnce sync.Once
	file_api_proto_post_proto_rawDescData = file_api_proto_post_proto_rawDesc
)

func file_api_proto_post_proto_rawDescGZIP() []byte {
	file_api_proto_post_proto_rawDescOnce.Do(func() {
		file_api_proto_post_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_proto_post_proto_rawDescData)
	})
	return file_api_proto_post_proto_rawDescData
}

var file_api_proto_post_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_proto_post_proto_goTypes = []any{
	(*PostByIdRequest)(nil),    // 0: proto.PostByIdRequest
	(*CreatePostRequest)(nil),  // 1: proto.CreatePostRequest
	(*CreatePostResponse)(nil), // 2: proto.CreatePostResponse
	(*CategoryVars)(nil),       // 3: proto.CategoryVars
	(*Tag)(nil),                // 4: proto.Tag
	(*Post)(nil),               // 5: proto.Post
}
var file_api_proto_post_proto_depIdxs = []int32{
	3, // 0: proto.Post.vars:type_name -> proto.CategoryVars
	4, // 1: proto.Post.tags:type_name -> proto.Tag
	0, // 2: proto.PostService.GetPostById:input_type -> proto.PostByIdRequest
	1, // 3: proto.PostService.CreatePost:input_type -> proto.CreatePostRequest
	5, // 4: proto.PostService.GetPostById:output_type -> proto.Post
	2, // 5: proto.PostService.CreatePost:output_type -> proto.CreatePostResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_proto_post_proto_init() }
func file_api_proto_post_proto_init() {
	if File_api_proto_post_proto != nil {
		return
	}
	file_api_proto_post_proto_msgTypes[1].OneofWrappers = []any{}
	file_api_proto_post_proto_msgTypes[5].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_proto_post_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_post_proto_goTypes,
		DependencyIndexes: file_api_proto_post_proto_depIdxs,
		MessageInfos:      file_api_proto_post_proto_msgTypes,
	}.Build()
	File_api_proto_post_proto = out.File
	file_api_proto_post_proto_rawDesc = nil
	file_api_proto_post_proto_goTypes = nil
	file_api_proto_post_proto_depIdxs = nil
}
