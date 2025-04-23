package encoding

import (
	"github.com/dubbogo/grpc-go/encoding/proto_wrapper_api"
	"github.com/dubbogo/grpc-go/mem"
	perrors "github.com/pkg/errors"
)

// PBWrapperTwoWayCodec is codec impl of pb
type PBWrapperTwoWayCodec struct {
	pbCodec CodecV2
	codec   Codec
	name    string
}

// NewPBWrapperTwoWayCodec new common.TwoWayCodec PBWrapperTwoWayCodec with @codecName defined Codec inside
func NewPBWrapperTwoWayCodec(name string, innerCodec Codec, pbCodec CodecV2) TwoWayCodecV2 {
	return &PBWrapperTwoWayCodec{
		codec:   innerCodec,
		name:    name,
		pbCodec: pbCodec,
	}
}

func (h *PBWrapperTwoWayCodec) Name() string {
	return h.name
}

// MarshalRequest marshal interface @v to []byte
func (h *PBWrapperTwoWayCodec) MarshalRequest(v any) (mem.BufferSlice, error) {
	reqList := v.([]any)
	argsBytes := make([][]byte, 0, len(reqList))
	argsTypes := make([]string, 0, len(reqList))
	for _, value := range reqList {
		data, err := h.codec.Marshal(value)
		if err != nil {
			return nil, err
		}
		argsBytes = append(argsBytes, data)
		argsTypes = append(argsTypes, GetArgType(value))
	}

	wrapperRequest := &proto_wrapper_api.TripleRequestWrapper{
		SerializeType: h.Name(),
		Args:          argsBytes,
		ArgTypes:      argsTypes,
	}
	return h.pbCodec.Marshal(wrapperRequest)
}

// UnmarshalRequest unmarshal bytes @data to interface
func (h *PBWrapperTwoWayCodec) UnmarshalRequest(data mem.BufferSlice, v any) error {
	wrapperRequest := proto_wrapper_api.TripleRequestWrapper{}
	err := h.pbCodec.Unmarshal(data, &wrapperRequest)
	if err != nil {
		return err
	}

	paramsInterfaces := v.([]any)
	if len(paramsInterfaces) != len(wrapperRequest.Args) {
		return perrors.Errorf("error ,request params len is %d, but exported method has %d", len(wrapperRequest.Args), len(paramsInterfaces))
	}

	for idx, value := range wrapperRequest.Args {
		if err := h.codec.Unmarshal(value, paramsInterfaces[idx]); err != nil {
			return err
		}
	}

	return nil
}

// MarshalResponse marshal interface @v to []byte
func (h *PBWrapperTwoWayCodec) MarshalResponse(v interface{}) (mem.BufferSlice, error) {
	data, err := h.codec.Marshal(v)
	if err != nil {
		return nil, err
	}

	wrapperRequest := &proto_wrapper_api.TripleResponseWrapper{
		SerializeType: h.name,
		Data:          data,
		Type:          GetArgType(v),
	}
	return h.pbCodec.Marshal(wrapperRequest)
}

// UnmarshalResponse unmarshal bytes @data to interface
func (h *PBWrapperTwoWayCodec) UnmarshalResponse(data mem.BufferSlice, v interface{}) error {
	wrapperResponse := proto_wrapper_api.TripleResponseWrapper{}
	err := h.pbCodec.Unmarshal(data, &wrapperResponse)
	if err != nil {
		return err
	}
	if v == nil { // empty respose
		return nil
	}
	return h.codec.Unmarshal(wrapperResponse.Data, v)
}

// PBTwoWayCodec is pb impl of TwoWayCodec
type PBTwoWayCodec struct {
	codec Codec
}
