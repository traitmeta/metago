package tools

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestGetBalanceToday(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				address: "bc1p9wf3h8hwafqhmh2djyy2yxqvfnu5njrauteu70pd98s63m5zqzrq23ayqn",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBalanceToday(tt.args.address)
			t.Log(decimal.NewFromInt(got).Div(decimal.New(1, 8)).String())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBalanceToday() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBalanceToday() got = %v, want %v", got, tt.want)
			}
		})
	}
}
