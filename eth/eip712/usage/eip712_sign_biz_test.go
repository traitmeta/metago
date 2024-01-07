package usage

import (
	"testing"
)

func TestEIP712SignBiz_recover(t *testing.T) {
	type args struct {
		dataHash  string
		signature string
	}
	tests := []struct {
		name              string
		wantRecoveredAddr string
		args              args
		wantErr           bool
	}{
		{
			name: "test",
			args: args{
				dataHash:  "0x846c2aa6277c50980556cccc77d2c9bcde1258b00228ce062da733268802fa01",
				signature: "0xc5f7a27fb56690c5ca607b2ddc5efd58ca8b0290dab78a72a478e5d32e7facb562284309264b9cc0a09473e20ba5271587e3aa0bf9bc00b2296c560eb7b6035b1b",
			},
			wantRecoveredAddr: "0x33A4F229Bd34ea7783302c99FFD6e26324Bd2789",
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecoveredAddr, err := recover(tt.args.signature, tt.args.dataHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.recover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRecoveredAddr != tt.wantRecoveredAddr {
				t.Errorf("EIP712SignBiz.recover() = %v, want %v", gotRecoveredAddr, tt.wantRecoveredAddr)
			}
		})
	}
}

func TestEIP712SignBiz_SignForOrder(t *testing.T) {
	type args struct {
		priv      string
		eip712Str string
	}
	tests := []struct {
		name          string
		args          args
		wantSignature string
		wantAddress   string
		wantErr       bool
	}{
		{
			name: "test remix",
			args: args{
				priv:      "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				eip712Str: `{"types":{"OrderComponents":[{"name":"offerer","type":"address"},{"name":"zone","type":"address"},{"name":"offer","type":"OfferItem[]"},{"name":"consideration","type":"ConsiderationItem[]"},{"name":"orderType","type":"uint8"},{"name":"startTime","type":"uint256"},{"name":"endTime","type":"uint256"},{"name":"zoneHash","type":"bytes32"},{"name":"salt","type":"uint256"},{"name":"conduitKey","type":"bytes32"},{"name":"counter","type":"uint256"}],"OfferItem":[{"name":"itemType","type":"uint8"},{"name":"token","type":"address"},{"name":"identifierOrCriteria","type":"uint256"},{"name":"startAmount","type":"uint256"},{"name":"endAmount","type":"uint256"}],"ConsiderationItem":[{"name":"itemType","type":"uint8"},{"name":"token","type":"address"},{"name":"identifierOrCriteria","type":"uint256"},{"name":"startAmount","type":"uint256"},{"name":"endAmount","type":"uint256"},{"name":"recipient","type":"address"}],"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}]},"domain":{"name":"Seaport","version":"1.5","chainId":"1","verifyingContract":"0xd8b934580fce35a11b58c6d73adee468a2833fa8"},"primaryType":"OrderComponents","message":{"offerer":"0xab8483f64d9c6d1ecf9b849ae677dd3315835cb2","zone":"0x0000000000000000000000000000000000000000","offer":[{"itemType":"2","token":"0xf8e81d47203a594245e36c48e151709f0c19fbe8","identifierOrCriteria":"1","startAmount":"1","endAmount":"1"}],"consideration":[{"itemType":"0","token":"0x0000000000000000000000000000000000000000","identifierOrCriteria":"0","startAmount":"20000000","endAmount":"20000000","recipient":"0xab8483f64d9c6d1ecf9b849ae677dd3315835cb2"},{"itemType":"0","token":"0x0000000000000000000000000000000000000000","identifierOrCriteria":"0","startAmount":"1000000","endAmount":"1000000","recipient":"0x5b38da6a701c568545dcfcb03fcb875f56beddc4"}],"orderType":"0","startTime":"1291123844","endTime":"1791123844","zoneHash":"0x0000000000000000000000000000000000000000000000000000000000000000","salt":"24446860302761739304752683030156737591518664810215442929816108075358245614181","conduitKey":"0x0000000000000000000000000000000000000000000000000000000000000000","counter":"0"}}`,
			},
			wantSignature: "0x78fe25a72faa95d43a10cc16cd4de40531d5e649f5ff6bdcc935380b7997a71146aa50e35c4822cd8023715e95a131e19873df629437514c468e9defd35316371b",
			wantAddress:   "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
			wantErr:       false,
		},
		{
			name: "test testOrderData",
			args: args{
				priv:      "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				eip712Str: testOrderData,
			},
			wantSignature: "0x9b4bde85527149134cca002cab51df8a22f1ab16fd5b6e8a7597a091daddf587561d17b937b3c2bcc149fe2ff9628533628021661e139e00d4d4baa9baa904e61b",
			wantAddress:   "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
			wantErr:       false,
		},
		{
			name: "test testOfferData",
			args: args{
				priv:      "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				eip712Str: testOfferData,
			},
			wantSignature: "0x02ac50cbb04f03952edfbc5c6285a3b2a800ecf0e953f8bfdae94ce8d5af00ad58494b6efb9c9952e0d9613241b719c2a9804ea6635d63c592635eca3ce630971b",
			wantAddress:   "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
			wantErr:       false,
		},
		{
			name: "test testZoneSigner",
			args: args{
				priv:      "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				eip712Str: testZoneSigner,
			},
			wantSignature: "0x9ecc03c46989634fe019223c9aa94871a8e161e8173ae758930fd7b5b77d8e784f75d938a975d35c377752f97ec5b31e516abb98bcaa7bc3bedc78a640da623f1c",
			wantAddress:   "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
			wantErr:       false,
		},
		{
			name: "test testOrderNft",
			args: args{
				priv:      "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				eip712Str: testOrderNft,
			},
			wantSignature: "0xe10f98052ecbdb31a0fb564b98ec3171a6996f8c8618b62c688130ea8e36b54a67a6e14db9a8788a73f481f50e1bbe22ba0d4364fb5731f26b859e21797384041b",
			wantAddress:   "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, err := SignData(tt.args.priv, tt.args.eip712Str)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.SignForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("EIP712SignBiz.SignForOrder() = %v, want %v", gotSignature, tt.wantSignature)
			}

			addr, err := VerifySignature(gotSignature, tt.args.eip712Str)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.SignForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if addr != tt.wantAddress {
				t.Errorf("EIP712SignBiz.SignForOrder() = %v, want %v", addr, tt.wantAddress)
			}
		})
	}
}
