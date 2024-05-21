package runes

import (
	"math"
	"reflect"
	"testing"
)

func Test_encodeVarint(t *testing.T) {
	type args struct {
		n uint128
		v *[]byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test 0",
			args: args{
				n: uint128{
					lo: 0,
					hi: 0,
				},
				v: &[]byte{},
			},
			want: []byte{0},
		},
		{
			name: "test max 128",
			args: args{
				n: uint128{
					lo: math.MaxUint64,
					hi: math.MaxUint64,
				},
				v: &[]byte{},
			},
			want: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeVarint(tt.args.n, tt.args.v)
			if !reflect.DeepEqual(*tt.args.v, tt.want) {
				t.Errorf("Decipher() = %v, want %v", *tt.args.v, tt.want)
			}
		})
	}
}
