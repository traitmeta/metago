package runes

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestDecipher(t *testing.T) {
	type args struct {
		script string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test mint",
			args: args{
				script: "6a5d0b00c0a2338e01b09f1aff07",
			},
			want:    "00c0a2338e01b09f1aff07",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeScript, _ := hex.DecodeString(tt.args.script)
			got, err := Decipher(decodeScript)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decipher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want, _ := hex.DecodeString(tt.want)
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Decipher() = %v, want %v", got, want)
			}
		})
	}
}
