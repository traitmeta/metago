package txbuilder

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"github.com/traitmeta/metago/btc/runes-tools/txbuilder/runes"
)

func TestMintTxBuilder_BuildMintTx(t *testing.T) {
	type fields struct {
		tx                *wire.MsgTx
		net               *chaincfg.Params
		runesCli          *runes.Client
		privateKey        *btcec.PrivateKey
		prevOutputFetcher *txscript.MultiPrevOutFetcher
		req               runes.EtchRequest
	}
	type args struct {
		prev PrevInfo
		req  runes.EtchRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *wire.MsgTx
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				tx:                nil,
				net:               &chaincfg.TestNet3Params,
				runesCli:          nil,
				privateKey:        nil,
				prevOutputFetcher: nil,
				req: runes.EtchRequest{
					FeeRate:     40,
					RuneID:      "322:1",
					Destination: "",
				},
			},
			args:    args{},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			privateKey, err := btcec.NewPrivateKey() // note: 创建一个密钥对，用来构建reveal tx
			if err != nil {
				t.Errorf("NewPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			mintBuilder := NewMintTxBuilder(privateKey, tt.fields.net)
			got, err := mintBuilder.BuildMintTx(tt.args.prev, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildMintTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildMintTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
