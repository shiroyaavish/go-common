package create_svc

import (
	"github.com/IntelXLabs-LLC/go-common/models/wrappers"
	"github.com/IntelXLabs-LLC/go-common/models/wrappers/vpn_operations/create_svc/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type CreateClientRequest struct {
	*wrappers.GRPCError
	*proto.CreateClientRequest
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeviceID    string `json:"device_id" protobuf:"bytes,1,opt,name=device_id"`
	OperationID string `json:"operation_id" protobuf:"bytes,2,opt,name=operation_id"`
}

type CreateClientResponse struct {
	*wrappers.GRPCError
	*proto.CreateClientResponse
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status            string `json:"status" protobuf:"bytes,1,opt,name=status"`
	Code              int64  `json:"code" protobuf:"varint,2,opt,name=code"`
	Base64EncodedConf string `json:"base_64_encoded_conf" protobuf:"bytes,3,opt,name=base_64_encoded_conf"`
	OperationID       string `json:"operation_id" protobuf:"bytes,4,opt,name=operation_id"`
}
