package sync

import "testing"

func TestMatchElementPattern(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test normal without pattern",
			args: args{
				content: "dmt.11.element",
			},
			want: true,
		},
		{
			name: "test normal with pattern",
			args: args{
				content: "dmt.11.11.element",
			},
			want: true,
		},
		{
			name: "test normal with mix name",
			args: args{
				content: "123dmt123.11.11.element",
			},
			want: true,
		},
		{
			name: "test normal with mix pattern",
			args: args{
				content: "dmt.83ddms.11.element",
			},
			want: true,
		},
		{
			name: "test when name has empty in middle",
			args: args{
				content: "123  dmt123.11.11.element",
			},
			want: true,
		},
		{
			name: "test when name has empty in end",
			args: args{
				content: "dmt .83ddms.11.element",
			},
			want: true,
		},
		{
			name: "test with music symbol",
			args: args{
				content: "ᚰ-runa.ᚰ.11.element",
			},
			want: true,
		},
		{
			name: "test with frog name",
			args: args{
				content: "🐸🐸🐸.83_ddms.11.element",
			},
			want: true,
		},
		{
			name: "test  with frog name 2",
			args: args{
				content: "🐸🐸🐸.83ddms.11.element",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchElementPattern(tt.args.content); got != tt.want {
				t.Errorf("MatchElementPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
