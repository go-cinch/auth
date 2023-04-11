// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.17.3
// source: reason-proto/reason.proto

package reason

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type ErrorReason int32

const (
	ErrorReason_TOO_MANY_REQUESTS ErrorReason = 0
	ErrorReason_ILLEGAL_PARAMETER ErrorReason = 1
	ErrorReason_NOT_FOUND         ErrorReason = 2
	ErrorReason_UNAUTHORIZED      ErrorReason = 3
	ErrorReason_FORBIDDEN         ErrorReason = 4
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0: "TOO_MANY_REQUESTS",
		1: "ILLEGAL_PARAMETER",
		2: "NOT_FOUND",
		3: "UNAUTHORIZED",
		4: "FORBIDDEN",
	}
	ErrorReason_value = map[string]int32{
		"TOO_MANY_REQUESTS": 0,
		"ILLEGAL_PARAMETER": 1,
		"NOT_FOUND":         2,
		"UNAUTHORIZED":      3,
		"FORBIDDEN":         4,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_reason_proto_reason_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_reason_proto_reason_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_reason_proto_reason_proto_rawDescGZIP(), []int{0}
}

var File_reason_proto_reason_proto protoreflect.FileDescriptor

var file_reason_proto_reason_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x72, 0x65, 0x61,
	0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x8f, 0x01, 0x0a, 0x0b,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x11, 0x54,
	0x4f, 0x4f, 0x5f, 0x4d, 0x41, 0x4e, 0x59, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x53,
	0x10, 0x00, 0x1a, 0x04, 0xa8, 0x45, 0xad, 0x03, 0x12, 0x1b, 0x0a, 0x11, 0x49, 0x4c, 0x4c, 0x45,
	0x47, 0x41, 0x4c, 0x5f, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x45, 0x54, 0x45, 0x52, 0x10, 0x01, 0x1a,
	0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x13, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55,
	0x4e, 0x44, 0x10, 0x02, 0x1a, 0x04, 0xa8, 0x45, 0x90, 0x03, 0x12, 0x16, 0x0a, 0x0c, 0x55, 0x4e,
	0x41, 0x55, 0x54, 0x48, 0x4f, 0x52, 0x49, 0x5a, 0x45, 0x44, 0x10, 0x03, 0x1a, 0x04, 0xa8, 0x45,
	0x91, 0x03, 0x12, 0x13, 0x0a, 0x09, 0x46, 0x4f, 0x52, 0x42, 0x49, 0x44, 0x44, 0x45, 0x4e, 0x10,
	0x04, 0x1a, 0x04, 0xa8, 0x45, 0x93, 0x03, 0x1a, 0x04, 0xa0, 0x45, 0xf4, 0x03, 0x42, 0x30, 0x0a,
	0x09, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x13, 0x2e, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x3b, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0xa2, 0x02, 0x0b, 0x41, 0x50, 0x49, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reason_proto_reason_proto_rawDescOnce sync.Once
	file_reason_proto_reason_proto_rawDescData = file_reason_proto_reason_proto_rawDesc
)

func file_reason_proto_reason_proto_rawDescGZIP() []byte {
	file_reason_proto_reason_proto_rawDescOnce.Do(func() {
		file_reason_proto_reason_proto_rawDescData = protoimpl.X.CompressGZIP(file_reason_proto_reason_proto_rawDescData)
	})
	return file_reason_proto_reason_proto_rawDescData
}

var file_reason_proto_reason_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_reason_proto_reason_proto_goTypes = []interface{}{
	(ErrorReason)(0), // 0: reason.v1.ErrorReason
}
var file_reason_proto_reason_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_reason_proto_reason_proto_init() }
func file_reason_proto_reason_proto_init() {
	if File_reason_proto_reason_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_reason_proto_reason_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_reason_proto_reason_proto_goTypes,
		DependencyIndexes: file_reason_proto_reason_proto_depIdxs,
		EnumInfos:         file_reason_proto_reason_proto_enumTypes,
	}.Build()
	File_reason_proto_reason_proto = out.File
	file_reason_proto_reason_proto_rawDesc = nil
	file_reason_proto_reason_proto_goTypes = nil
	file_reason_proto_reason_proto_depIdxs = nil
}
