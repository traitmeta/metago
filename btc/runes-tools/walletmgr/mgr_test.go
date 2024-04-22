package walletmgr

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestInitAndCacheWalletWif(t *testing.T) {
	type args struct {
		cacheDir string
		net      *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				cacheDir: "./runes/wallet",
				net:      &chaincfg.TestNet3Params,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitAndCacheWalletWif(tt.args.cacheDir, tt.args.net); (err != nil) != tt.wantErr {
				t.Errorf("InitAndCacheWalletWif() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitWalletMgr(t *testing.T) {
	type args struct {
		cacheDir string
		net      *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		want    []Wallet
		wantErr bool
	}{
		{
			name: "",
			args: args{
				cacheDir: "./runes/wallet",
				net:      &chaincfg.TestNet3Params,
			},
			want:    []Wallet{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := InitWalletMgr(tt.args.cacheDir, tt.args.net)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitWalletMgr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got.wallets) != 10 {
				t.Errorf("InitWalletMgr() got = %v, want %v", got, tt.want)
				return
			}

			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("InitWalletMgr() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
