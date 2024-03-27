package tools

import (
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func Test_pubkeyToAddr(t *testing.T) {
	type args struct {
		pubkeyHex   string
		addressType AddressType
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test compress p2tr address",
			args: args{
				pubkeyHex:   "03ac18ec81d4b2e3ce4445aa1f403b0c03d75e725eed4a1608e1d57b6747ccafc4",
				addressType: TaprootPubKey,
			},
			want:    "bc1pwq6j83lxu552qj43y455e3mxvpt3l52ur3v8400882hn8tqwuh9s5fredq",
			wantErr: false,
		},
		{
			name: "test compress p2sh address",
			args: args{
				pubkeyHex:   "037870606049bbe5f59581c022af6864cdc122b9927e5febe57c14c5150b164919",
				addressType: WitnessPubKey,
			},
			want:    "bc1q2y7gm4g90h508upkwdava36uqysxwa2hxrx82x",
			wantErr: false,
		},
		{
			name: "test compress p2pkh address",
			args: args{
				pubkeyHex:   "037870606049bbe5f59581c022af6864cdc122b9927e5febe57c14c5150b164919",
				addressType: PubKeyHash,
			},
			want:    "18QYGuZH593aQN1KWmYTm9PGjC8UPvtSjt",
			wantErr: false,
		},
		{
			name: "test compress p2wpkh address",
			args: args{
				pubkeyHex:   "037870606049bbe5f59581c022af6864cdc122b9927e5febe57c14c5150b164919",
				addressType: WitnessPubKey,
			},
			want:    "bc1q2y7gm4g90h508upkwdava36uqysxwa2hxrx82x",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BtcPubKeyToBtcAddress(tt.args.pubkeyHex, tt.args.addressType, &chaincfg.MainNetParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("pubkeyToAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pubkeyToAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_publickKeyToWalletAddress(t *testing.T) {
	type args struct {
		btcAddress string
		nonce      uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "test address bc1p84havkv55wkvfmpfela97ulju2key6tw8z6jyzsae4ajzjwmqdyqcd2hes",
			args: args{
				btcAddress: "bc1p84havkv55wkvfmpfela97ulju2key6tw8z6jyzsae4ajzjwmqdyqcd2hes",
				nonce:      1,
			},
			want: "0x7b646235EC56e6468084a1F6bF70a4f94a4515A5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BtcAddressToContractAddress(tt.args.btcAddress, tt.args.nonce); got != tt.want {
				t.Errorf("publickKeyToWalletAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toHexAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test address",
			args: args{
				address: "bc1p9wf3h8hwafqhmh2djyy2yxqvfnu5njrauteu70pd98s63m5zqzrq23ayqn",
			},
			want: "0x0218a6A1954Cb08246CFd0730043Eb0022c84c50",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BtcAddressToBvmAddress(tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BtcAddressToBvmAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPublicKeyBVMAddress(t *testing.T) {
	type args struct {
		pubKey string
		addr   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bc1q false",
			args: args{
				pubKey: "020852f2e311c060f473e9e00e1118905191021b934a2722c3542bdfd032fd2661",
				addr:   "0x10046269dcbc7d51d2e6845c46a2b178bd5421cb",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := VerifyPublicKeyBVMAddress(tt.args.pubKey, tt.args.addr, &chaincfg.MainNetParams)
			if res != tt.want {
				t.Errorf("CheckPublicKeyAndAddress() error = %v, want %v", res, tt.want)
				return
			}
		})
	}
}
