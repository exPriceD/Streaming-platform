// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.20.3
// source: proto/chat.proto

package proto

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

// Сообщение чата
type ChatMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`       // ID пользователя
	StreamId      string                 `protobuf:"bytes,2,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"` // ID стрима
	Content       string                 `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`                   // Текст сообщения
	Timestamp     int64                  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`              // Временная метка
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatMessage) Reset() {
	*x = ChatMessage{}
	mi := &file_proto_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatMessage) ProtoMessage() {}

func (x *ChatMessage) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatMessage.ProtoReflect.Descriptor instead.
func (*ChatMessage) Descriptor() ([]byte, []int) {
	return file_proto_chat_proto_rawDescGZIP(), []int{0}
}

func (x *ChatMessage) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ChatMessage) GetStreamId() string {
	if x != nil {
		return x.StreamId
	}
	return ""
}

func (x *ChatMessage) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *ChatMessage) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

// Ответ на запрос чата
type ChatResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Status        string                 `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`                                 // Статус (например, "success" или "error")
	ErrorMessage  string                 `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"` // Сообщение об ошибке (если есть)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatResponse) Reset() {
	*x = ChatResponse{}
	mi := &file_proto_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatResponse) ProtoMessage() {}

func (x *ChatResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatResponse.ProtoReflect.Descriptor instead.
func (*ChatResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_proto_rawDescGZIP(), []int{1}
}

func (x *ChatResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ChatResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

// Запрос истории сообщений
type ChatHistoryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StreamId      string                 `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"` // ID стрима
	Limit         int32                  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`                      // Количество сообщений (по умолчанию 50)
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatHistoryRequest) Reset() {
	*x = ChatHistoryRequest{}
	mi := &file_proto_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatHistoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatHistoryRequest) ProtoMessage() {}

func (x *ChatHistoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatHistoryRequest.ProtoReflect.Descriptor instead.
func (*ChatHistoryRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_proto_rawDescGZIP(), []int{2}
}

func (x *ChatHistoryRequest) GetStreamId() string {
	if x != nil {
		return x.StreamId
	}
	return ""
}

func (x *ChatHistoryRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// Ответ с историей сообщений
type ChatHistoryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Messages      []*ChatMessage         `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ChatHistoryResponse) Reset() {
	*x = ChatHistoryResponse{}
	mi := &file_proto_chat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ChatHistoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatHistoryResponse) ProtoMessage() {}

func (x *ChatHistoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatHistoryResponse.ProtoReflect.Descriptor instead.
func (*ChatHistoryResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_proto_rawDescGZIP(), []int{3}
}

func (x *ChatHistoryResponse) GetMessages() []*ChatMessage {
	if x != nil {
		return x.Messages
	}
	return nil
}

// Стрим сообщений в реальном времени
type StreamMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	StreamId      string                 `protobuf:"bytes,1,opt,name=stream_id,json=streamId,proto3" json:"stream_id,omitempty"` // ID стрима
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamMessageRequest) Reset() {
	*x = StreamMessageRequest{}
	mi := &file_proto_chat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamMessageRequest) ProtoMessage() {}

func (x *StreamMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamMessageRequest.ProtoReflect.Descriptor instead.
func (*StreamMessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_proto_rawDescGZIP(), []int{4}
}

func (x *StreamMessageRequest) GetStreamId() string {
	if x != nil {
		return x.StreamId
	}
	return ""
}

var File_proto_chat_proto protoreflect.FileDescriptor

var file_proto_chat_proto_rawDesc = string([]byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x63, 0x68, 0x61, 0x74, 0x22, 0x7b, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x4b, 0x0a, 0x0c, 0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x23, 0x0a,
	0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x47, 0x0a, 0x12, 0x43, 0x68, 0x61, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x72, 0x65,
	0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x44, 0x0a, 0x13, 0x43,
	0x68, 0x61, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x2d, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x22, 0x33, 0x0a, 0x14, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x49, 0x64, 0x32, 0xcd, 0x01, 0x0a, 0x0b, 0x43, 0x68, 0x61, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x68, 0x61,
	0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x12, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x18,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e,
	0x43, 0x68, 0x61, 0x74, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x74, 0x72,
	0x65, 0x61, 0x6d, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x30, 0x01, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x78, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x2f, 0x53, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x2d, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_chat_proto_rawDescOnce sync.Once
	file_proto_chat_proto_rawDescData []byte
)

func file_proto_chat_proto_rawDescGZIP() []byte {
	file_proto_chat_proto_rawDescOnce.Do(func() {
		file_proto_chat_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_chat_proto_rawDesc), len(file_proto_chat_proto_rawDesc)))
	})
	return file_proto_chat_proto_rawDescData
}

var file_proto_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_chat_proto_goTypes = []any{
	(*ChatMessage)(nil),          // 0: chat.ChatMessage
	(*ChatResponse)(nil),         // 1: chat.ChatResponse
	(*ChatHistoryRequest)(nil),   // 2: chat.ChatHistoryRequest
	(*ChatHistoryResponse)(nil),  // 3: chat.ChatHistoryResponse
	(*StreamMessageRequest)(nil), // 4: chat.StreamMessageRequest
}
var file_proto_chat_proto_depIdxs = []int32{
	0, // 0: chat.ChatHistoryResponse.messages:type_name -> chat.ChatMessage
	0, // 1: chat.ChatService.SendMessage:input_type -> chat.ChatMessage
	2, // 2: chat.ChatService.GetChatHistory:input_type -> chat.ChatHistoryRequest
	4, // 3: chat.ChatService.StreamMessages:input_type -> chat.StreamMessageRequest
	1, // 4: chat.ChatService.SendMessage:output_type -> chat.ChatResponse
	3, // 5: chat.ChatService.GetChatHistory:output_type -> chat.ChatHistoryResponse
	0, // 6: chat.ChatService.StreamMessages:output_type -> chat.ChatMessage
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_chat_proto_init() }
func file_proto_chat_proto_init() {
	if File_proto_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_chat_proto_rawDesc), len(file_proto_chat_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chat_proto_goTypes,
		DependencyIndexes: file_proto_chat_proto_depIdxs,
		MessageInfos:      file_proto_chat_proto_msgTypes,
	}.Build()
	File_proto_chat_proto = out.File
	file_proto_chat_proto_goTypes = nil
	file_proto_chat_proto_depIdxs = nil
}
