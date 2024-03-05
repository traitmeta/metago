package tap

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
)

func TestElement_IsAvailable(t *testing.T) {
	type fields struct {
		name    string
		pattern string
		field   string
	}
	type args struct {
		block *wire.MsgBlock
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test",
			fields: fields{
				name:    "bit",
				pattern: "3b",
				field:   "11",
			},
			args: args{
				block: &wire.MsgBlock{
					Header: wire.BlockHeader{
						Bits: 386120285,
					},
					Transactions: nil,
				},
			},
			want: true,
		},
		{
			name: "test",
			fields: fields{
				name:    "bit",
				pattern: "",
				field:   "11",
			},
			args: args{
				block: &wire.MsgBlock{
					Header: wire.BlockHeader{
						Bits: 386120285,
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Element{
				name:    tt.fields.name,
				pattern: tt.fields.pattern,
				field:   tt.fields.field,
			}
			if got := e.IsAvailable(tt.args.block); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}
