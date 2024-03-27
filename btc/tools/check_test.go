package tools

import "testing"

func TestDustCheck(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{

		{
			name: "test 1 prefix",
			args: args{
				address: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			},
			want: 546,
		},
		{
			name: "test 3 prefix",
			args: args{
				address: "3GfvGfUg5yezijpiBPXRYnW8qAd4DiT3gL",
			},
			want: 540,
		},
		{
			name: "test bc1p prefix",
			args: args{
				address: "bc1p07pv3956cuq4xlhm4skwkha4kdvr43qz7t5zjd8uyvhq2jsewvtqmqvhx0",
			},
			want: 330,
		},
		{
			name: "test bc1q prefix",
			args: args{
				address: "bc1qd3eljyyzdagslxrzlpl57fkrmuztgr5mnrjr04",
			},
			want: 294,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DustCheck(tt.args.address); got != tt.want {
				t.Errorf("DustCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
