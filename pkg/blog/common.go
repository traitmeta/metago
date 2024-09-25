package blog

import (
	"context"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/traitmeta/metago/pkg/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LogConf struct {
	ServiceName string `toml:"service_name" mapstructure:"service_name" json:"service_name"`
	Mode        string `toml:"mode" json:"mode"`
	Path        string `toml:"path" json:"path"`
	Level       string `toml:"level" json:"level"`
	Compress    bool   `toml:"compress" json:"compress"`
	KeepDays    int    `toml:"keep_days" mapstructure:"keep_days" json:"keep_days"`
}

type ErrorToCode func(err error) codes.Code

func DefaultErrorToCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	switch e := err.(type) {
	case interface{ GRPCStatus() *status.Status }:
		return status.Code(err)
	case *errcode.Err:
		return codes.Code(e.Code())
	default:
		return codes.Unknown
	}
}

type Decider func(methodName string, err error) bool

func DefaultDeciderMethod(methodName string, err error) bool {
	return true
}

type RecoveryHandlerContextFunc func(ctx context.Context, p interface{}) (err error)

type ServerLoggingDecider func(ctx context.Context, methodName string, servingObject interface{}) bool

type ClientLoggingDecider func(ctx context.Context, methodName string) bool

type JsonPbMarshaler interface {
	Marshal(out io.Writer, pb proto.Message) error
}
