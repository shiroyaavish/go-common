package delete_svc

import (
	"github.com/shiroyaavish/go-common/models/wrappers"
	"github.com/shiroyaavish/go-common/models/wrappers/vpn_operations/delete_svc/proto"
)

type DeleteClientRequest struct {
	*wrappers.GRPCError
	*proto.DeleteClientRequest
}

type DeleteClientResponse struct {
	*wrappers.GRPCError
	*proto.DeleteClientResponse
}

type Request struct {
	DeviceID string `json:"device_id" protobuf:"bytes,1,opt,name=device_id"`
}

type Response struct {
	Status string `json:"status" protobuf:"bytes,1,opt,name=status"`
	Code   int64  `json:"code" protobuf:"varint,2,opt,name=code"`
}
