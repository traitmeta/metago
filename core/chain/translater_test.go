package chain

import (
	"testing"

	"github.com/traitmeta/metago/core/common"
)

func Test_truncateAddressHash(t *testing.T) {
	type args struct {
		addressHash string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test zero address",
			args: args{
				addressHash: "",
			},
			want: common.ZeroAddress,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateAddressHash(tt.args.addressHash); got != tt.want {
				t.Errorf("truncateAddressHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
