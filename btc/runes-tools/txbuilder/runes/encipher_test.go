package runes

import (
	"reflect"
	"testing"
)

func Test_encipherMint(t *testing.T) {
	type args struct {
		runesId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				runesId: "2587266:780",
			},
			want:    []byte{20, 130, 245, 157, 1, 20, 140, 6, 22, 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encipherMint(tt.args.runesId)
			if (err != nil) != tt.wantErr {
				t.Errorf("encipherMint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encipherMint() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encipher(t *testing.T) {
	type args struct {
		runesId string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				runesId: "2587266:780",
			},
			want:    []byte{106, 13, 10, 20, 130, 245, 157, 1, 20, 140, 6, 22, 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encipher(tt.args.runesId)
			if (err != nil) != tt.wantErr {
				t.Errorf("encipher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encipher() got = %v, want %v", got, tt.want)
			}
		})
	}
}
