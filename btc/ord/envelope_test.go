package ord

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestFromTapScript(t *testing.T) {
	type args struct {
		tapScript string
		input     int
	}
	tests := []struct {
		name    string
		args    args
		want    []Envelope
		wantLen int
		wantErr bool
	}{
		{
			name: "test tap pointer",
			args: args{
				// from txID 831b8d6637f1c7cf82859203b7a644104e8b40ce026bca848b383f0e3c820fdc
				tapScript: witnessList[0],
				input:     0,
			},
			want:    nil,
			wantLen: 1000,
			wantErr: false,
		},
		{
			name: "test tap pointer",
			args: args{
				tapScript: witnessList[1],
				input:     0,
			},
			want:    nil,
			wantLen: 1,
			wantErr: false,
		},
		{
			// txid : 84defb5e7db11f047b4bc58685b65654980be4a55e254a510048a71e0e43e532i0
			name: "test tap Provenance",
			args: args{
				tapScript: witnessList[2],
				input:     0,
			},
			want:    []Envelope{},
			wantLen: 1,
			wantErr: false,
		},
	}
	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tapScript, err := hex.DecodeString(tt.args.tapScript)
			if err != nil {
				t.Errorf("Decode TapScript err = %v", err)
			}
			got, err := FromTapScript(tapScript, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromTapScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if idx%2 == 0 {
				if len(got) != tt.wantLen {
					t.Errorf("FromTapScript() got_len = %v, want_len %v", len(got), tt.wantLen)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("FromTapScript() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEnvelope_GetProvenance(t *testing.T) {
	type args struct {
		tapScript string
		input     int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantLen int
		wantErr bool
	}{
		{
			// txid : 84defb5e7db11f047b4bc58685b65654980be4a55e254a510048a71e0e43e532i0
			name: "test tap Provenance",
			args: args{
				tapScript: witnessList[2],
				input:     0,
			},
			want:    "593015b9a76a11554f0a05c3b77a4723c6baaefb8bdd4175712a7320714b8ea8i0",
			wantLen: 1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tapScript, err := hex.DecodeString(tt.args.tapScript)
			if err != nil {
				t.Errorf("Decode TapScript err = %v", err)
			}
			got, err := FromTapScript(tapScript, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromTapScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 1 {
				t.Errorf("GetProvenance() = %v, want len %v", len(got), tt.want)
			}
			if provenance := got[0].GetProvenance(); provenance != tt.want {
				t.Errorf("GetProvenance() = %v, want %v", provenance, tt.want)
			}
		})
	}
}
