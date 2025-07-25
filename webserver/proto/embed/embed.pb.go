// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.0
// source: embed.proto

package embed

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EmbedRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Text          string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Category      string                 `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	Group         int64                  `protobuf:"varint,3,opt,name=group,proto3" json:"group,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmbedRequest) Reset() {
	*x = EmbedRequest{}
	mi := &file_embed_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmbedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbedRequest) ProtoMessage() {}

func (x *EmbedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_embed_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbedRequest.ProtoReflect.Descriptor instead.
func (*EmbedRequest) Descriptor() ([]byte, []int) {
	return file_embed_proto_rawDescGZIP(), []int{0}
}

func (x *EmbedRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *EmbedRequest) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *EmbedRequest) GetGroup() int64 {
	if x != nil {
		return x.Group
	}
	return 0
}

type Match struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Filename      string                 `protobuf:"bytes,1,opt,name=filename,proto3" json:"filename,omitempty"`
	PageNumber    int32                  `protobuf:"varint,2,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	Title         string                 `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Category      string                 `protobuf:"bytes,4,opt,name=category,proto3" json:"category,omitempty"`
	Content       string                 `protobuf:"bytes,5,opt,name=content,proto3" json:"content,omitempty"`
	Html          string                 `protobuf:"bytes,6,opt,name=html,proto3" json:"html,omitempty"`
	Score         float32                `protobuf:"fixed32,7,opt,name=score,proto3" json:"score,omitempty"`
	Meta          string                 `protobuf:"bytes,8,opt,name=meta,proto3" json:"meta,omitempty"`
	FileUrl       string                 `protobuf:"bytes,9,opt,name=file_url,json=fileUrl,proto3" json:"file_url,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Match) Reset() {
	*x = Match{}
	mi := &file_embed_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Match) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Match) ProtoMessage() {}

func (x *Match) ProtoReflect() protoreflect.Message {
	mi := &file_embed_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Match.ProtoReflect.Descriptor instead.
func (*Match) Descriptor() ([]byte, []int) {
	return file_embed_proto_rawDescGZIP(), []int{1}
}

func (x *Match) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *Match) GetPageNumber() int32 {
	if x != nil {
		return x.PageNumber
	}
	return 0
}

func (x *Match) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Match) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Match) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Match) GetHtml() string {
	if x != nil {
		return x.Html
	}
	return ""
}

func (x *Match) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *Match) GetMeta() string {
	if x != nil {
		return x.Meta
	}
	return ""
}

func (x *Match) GetFileUrl() string {
	if x != nil {
		return x.FileUrl
	}
	return ""
}

type EmbedResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Result        []*Match               `protobuf:"bytes,1,rep,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EmbedResponse) Reset() {
	*x = EmbedResponse{}
	mi := &file_embed_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EmbedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmbedResponse) ProtoMessage() {}

func (x *EmbedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_embed_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmbedResponse.ProtoReflect.Descriptor instead.
func (*EmbedResponse) Descriptor() ([]byte, []int) {
	return file_embed_proto_rawDescGZIP(), []int{2}
}

func (x *EmbedResponse) GetResult() []*Match {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_embed_proto protoreflect.FileDescriptor

const file_embed_proto_rawDesc = "" +
	"\n" +
	"\vembed.proto\x12\x05embed\"T\n" +
	"\fEmbedRequest\x12\x12\n" +
	"\x04text\x18\x01 \x01(\tR\x04text\x12\x1a\n" +
	"\bcategory\x18\x02 \x01(\tR\bcategory\x12\x14\n" +
	"\x05group\x18\x03 \x01(\x03R\x05group\"\xe9\x01\n" +
	"\x05Match\x12\x1a\n" +
	"\bfilename\x18\x01 \x01(\tR\bfilename\x12\x1f\n" +
	"\vpage_number\x18\x02 \x01(\x05R\n" +
	"pageNumber\x12\x14\n" +
	"\x05title\x18\x03 \x01(\tR\x05title\x12\x1a\n" +
	"\bcategory\x18\x04 \x01(\tR\bcategory\x12\x18\n" +
	"\acontent\x18\x05 \x01(\tR\acontent\x12\x12\n" +
	"\x04html\x18\x06 \x01(\tR\x04html\x12\x14\n" +
	"\x05score\x18\a \x01(\x02R\x05score\x12\x12\n" +
	"\x04meta\x18\b \x01(\tR\x04meta\x12\x19\n" +
	"\bfile_url\x18\t \x01(\tR\afileUrl\"5\n" +
	"\rEmbedResponse\x12$\n" +
	"\x06result\x18\x01 \x03(\v2\f.embed.MatchR\x06result2>\n" +
	"\bEmbedder\x122\n" +
	"\x05Embed\x12\x13.embed.EmbedRequest\x1a\x14.embed.EmbedResponseB\x0eZ\f/embed;embedb\x06proto3"

var (
	file_embed_proto_rawDescOnce sync.Once
	file_embed_proto_rawDescData []byte
)

func file_embed_proto_rawDescGZIP() []byte {
	file_embed_proto_rawDescOnce.Do(func() {
		file_embed_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_embed_proto_rawDesc), len(file_embed_proto_rawDesc)))
	})
	return file_embed_proto_rawDescData
}

var file_embed_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_embed_proto_goTypes = []any{
	(*EmbedRequest)(nil),  // 0: embed.EmbedRequest
	(*Match)(nil),         // 1: embed.Match
	(*EmbedResponse)(nil), // 2: embed.EmbedResponse
}
var file_embed_proto_depIdxs = []int32{
	1, // 0: embed.EmbedResponse.result:type_name -> embed.Match
	0, // 1: embed.Embedder.Embed:input_type -> embed.EmbedRequest
	2, // 2: embed.Embedder.Embed:output_type -> embed.EmbedResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_embed_proto_init() }
func file_embed_proto_init() {
	if File_embed_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_embed_proto_rawDesc), len(file_embed_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_embed_proto_goTypes,
		DependencyIndexes: file_embed_proto_depIdxs,
		MessageInfos:      file_embed_proto_msgTypes,
	}.Build()
	File_embed_proto = out.File
	file_embed_proto_goTypes = nil
	file_embed_proto_depIdxs = nil
}
