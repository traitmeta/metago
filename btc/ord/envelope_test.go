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
			want: []Envelope{
				{
					Input:  0,
					Offset: 0,
					TypeDataMap: map[int][]byte{
						1: []byte("text/plain;charset=utf-8"),
						0: []byte("{ \n  \"p\": \"brc-20\",\n  \"op\": \"deploy\",\n  \"tick\": \"ordi\",\n  \"max\": \"21000000\",\n  \"lim\": \"1000\"\n}"),
					},
					Payload: []byte("\u0001text/plain;charset=utf-8{ \n  \"p\": \"brc-20\",\n  \"op\": \"deploy\",\n  \"tick\": \"ordi\",\n  \"max\": \"21000000\",\n  \"lim\": \"1000\"\n}"),
					Pushnum: false,
					Stutter: false,
				},
			},
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
