package tools

import (
	"reflect"
	"testing"
)

func TestGetPrivateKeyFromWIF(t *testing.T) {
	type args struct {
		wifPrivate string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test",
			args: args{
				wifPrivate: "XV8mYI8C3ELcD2v3+OjiQw4PDz8WF7tXOub38INVxaY=",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPrivateKeyFromWIF(tt.args.wifPrivate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPrivateKeyFromWIF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFromBase64(t *testing.T) {
	type args struct {
		priv string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test ",
			args: args{
				priv: "XV8mYI8C3ELcD2v3+OjiQw4PDz8WF7tXOub38INVxaY=",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFromBase64(tt.args.priv)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFromBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Log(string(got))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFromBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}
