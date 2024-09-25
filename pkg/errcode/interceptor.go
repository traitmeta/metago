package errcode

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WrapErr(err error) error {
	if err == nil {
		return nil
	}

	switch e := err.(type) {
	case interface{ GRPCStatus() *status.Status }:
		return e.GRPCStatus().Err()
	case *Err:
		return status.Error(codes.Code(e.Code()), e.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

func UnwrapErr(err error) error {
	if err == nil {
		return nil
	}

	s, _ := status.FromError(err)
	c := uint32(s.Code())
	if c == CodeCustom {
		return NewCustomErr(s.Message())
	}

	if e, ok := codeToErr[c]; ok {
		return e
	}

	return err
}

func ErrInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	return resp, WrapErr(err)
}

func ErrStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, newWrappedServerStream(ss))
	return WrapErr(err)
}

func ErrClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return UnwrapErr(err)
}

func ErrStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	cs, err := streamer(ctx, desc, cc, method, opts...)
	return newWrappedClientStream(cs), UnwrapErr(err)
}

type wrappedServerStream struct {
	grpc.ServerStream
}

func newWrappedServerStream(ss grpc.ServerStream) *wrappedServerStream {
	if existing, ok := ss.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{ServerStream: ss}
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	return UnwrapErr(w.ServerStream.SendMsg(m))
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	return UnwrapErr(w.ServerStream.RecvMsg(m))
}

type wrappedClientStream struct {
	grpc.ClientStream
}

func newWrappedClientStream(cs grpc.ClientStream) *wrappedClientStream {
	if existing, ok := cs.(*wrappedClientStream); ok {
		return existing
	}
	return &wrappedClientStream{ClientStream: cs}
}

func (w *wrappedClientStream) SendMsg(m interface{}) error {
	return UnwrapErr(w.ClientStream.SendMsg(m))
}

func (w *wrappedClientStream) RecvMsg(m interface{}) error {
	return UnwrapErr(w.ClientStream.RecvMsg(m))
}

func (w *wrappedClientStream) CloseSend() error {
	return UnwrapErr(w.ClientStream.CloseSend())
}
