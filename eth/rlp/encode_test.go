package rlp

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEncodeToBytes(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test encode with rlp",
			args: args{
				item: Entity{
					WalletAddress: common.HexToAddress("0xc427059055842b4dd032faf00ccbe859a83d3234"),
					AccountNonce:  0,
				},
			},
			want:    []byte{0xd6, 0x94, 0xc4, 0x27, 0x05, 0x90, 0x55, 0x84, 0x2b, 0x4d, 0xd0, 0x32, 0xfa, 0xf0, 0x0c, 0xcb, 0xe8, 0x59, 0xa8, 0x3d, 0x32, 0x34, 0x80},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToBytes(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContractAddress(t *testing.T) {
	type args struct {
		input []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test generator contract address",
			args: args{
				input: []byte{0xd6, 0x94, 0xc4, 0x27, 0x05, 0x90, 0x55, 0x84, 0x2b, 0x4d, 0xd0, 0x32, 0xfa, 0xf0, 0x0c, 0xcb, 0xe8, 0x59, 0xa8, 0x3d, 0x32, 0x34, 0x80},
			},
			want: "0xf72ca65f9CC9E993fb60b48A6AE62e7E48eF8a19",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContractAddress(tt.args.input); got != tt.want {
				t.Errorf("ContractAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
