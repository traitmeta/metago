package xgrpc

import (
	"context"

	"github.com/golang/protobuf/jsonpb"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/traitmeta/metago/pkg/blog"
	"github.com/traitmeta/metago/pkg/blog/bzap"
)

var (
	SystemField                               = zap.String("system", "grpc")
	ServerField                               = zap.String("span.kind", "server")
	ClientField                               = zap.String("span.kind", "client")
	JsonPbMarshaller     blog.JsonPbMarshaler = &jsonpb.Marshaler{}
	defaultOptions                            = newDefaultOptions()
	defaultClientOptions                      = newDefaultClientOptions()
)

func newDefaultOptions() *bzap.Options {
	options := bzap.NewDefaultOption()
	opts := []bzap.Option{bzap.WithMessageProducer(GrpcMessageProducer)}
	for _, o := range opts {
		o(options)
	}

	return options
}

func newDefaultClientOptions() *bzap.Options {
	options := bzap.NewDefaultClientOption()
	opts := []bzap.Option{bzap.WithMessageProducer(GrpcMessageProducer)}
	for _, o := range opts {
		o(options)
	}

	return options
}

func GrpcRecoveryHandlerFunc(ctx context.Context, p interface{}) (err error) {
	err = status.Errorf(codes.Internal, "%v", p)
	cl := bzap.WithContext(ctx)
	cl.Panic("panic", zap.Error(err))

	return err
}

func GrpcMessageProducer(ctx context.Context, msg string, level zapcore.Level, err error, fields []zapcore.Field) {
	fields = append(fields, zap.Error(err))
	cl := bzap.WithContext(ctx)
	cl.Extract().Check(level, msg).Write(fields...)
}
