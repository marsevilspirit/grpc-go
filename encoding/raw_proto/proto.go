package raw_proto

import (
	"fmt"

	"github.com/dubbogo/grpc-go/encoding"
	"github.com/dubbogo/grpc-go/mem"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

// codec is a CodecV2 implementation with protobuf. It is the default codec for
// gRPC.
type ProtobufCodec struct{}

func (c *ProtobufCodec) Marshal(v any) (data mem.BufferSlice, err error) {
	vv := messageV2Of(v)
	if vv == nil {
		return nil, fmt.Errorf("proto: failed to marshal, message is %T, want proto.Message", v)
	}

	size := proto.Size(vv)
	if mem.IsBelowBufferPoolingThreshold(size) {
		buf, err := proto.Marshal(vv)
		if err != nil {
			return nil, err
		}
		data = append(data, mem.SliceBuffer(buf))
	} else {
		pool := mem.DefaultBufferPool()
		buf := pool.Get(size)
		if _, err := (proto.MarshalOptions{}).MarshalAppend((*buf)[:0], vv); err != nil {
			pool.Put(buf)
			return nil, err
		}
		data = append(data, mem.NewBuffer(buf, pool))
	}

	return data, nil
}

func (c *ProtobufCodec) Unmarshal(data mem.BufferSlice, v any) (err error) {
	vv := messageV2Of(v)
	if vv == nil {
		return fmt.Errorf("failed to unmarshal, message is %T, want proto.Message", v)
	}

	buf := data.MaterializeToBuffer(mem.DefaultBufferPool())
	defer buf.Free()
	// TODO: Upgrade proto.Unmarshal to support mem.BufferSlice. Right now, it's not
	//  really possible without a major overhaul of the proto package, but the
	//  vtprotobuf library may be able to support this.
	return proto.Unmarshal(buf.ReadOnlyData(), vv)
}

func messageV2Of(v any) proto.Message {
	switch v := v.(type) {
	case protoadapt.MessageV1:
		return protoadapt.MessageV2Of(v)
	case protoadapt.MessageV2:
		return v
	}

	return nil
}

func (c *ProtobufCodec) Name() string {
	return "raw_proto"
}

func NewProtobufCodec() encoding.CodecV2 {
	return &ProtobufCodec{}
}
