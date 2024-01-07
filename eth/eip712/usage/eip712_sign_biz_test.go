package usage

import (
	"context"
	"testing"
)

func TestVerifyHecoSigForBinding(t *testing.T) {
	tronAddress := "TRgqV7yHuqzkw878SwmhPqWxdB5eFtXDPz"
	hecoAddress := "0xb4d6498f5574a18119b07089a95a9bb53fa3ade5"
	sig := "0xf327a06e10fabd18b7f92ded0fef053c0a382ed9cc7ee02f2b0b7707e5dd6acf7dfcfed209e447f602040c984d331e92a1166b2ca7ce98a9159cfac386e048dc1c"

	isValid, err := VerifySigForBind(context.Background(), hecoAddress, tronAddress, sig)
	if err != nil || !isValid {
		t.Fail()
	}
}

func TestEIP712SignBiz_SignFoT(t *testing.T) {
	type args struct {
		priv string
	}
	tests := []struct {
		name          string
		args          args
		wantSignature string
		wantErr       bool
	}{
		{
			name: "test",
			args: args{
				priv: "2bdd8a43b8a055632f9a8b38d7c9463bfbe8340abbd6634f51bbf000cbe0ca50",
			},
			wantSignature: "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, err := SignFoTest(tt.args.priv)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.SignFoT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("EIP712SignBiz.SignFoT() = %v, want %v", gotSignature, tt.wantSignature)
			}
		})
	}
}

func TestEIP712SignBiz_recover(t *testing.T) {
	tests := []struct {
		name              string
		wantRecoveredAddr string
		wantErr           bool
	}{
		{
			name:              "test",
			wantRecoveredAddr: "",
			wantErr:           false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecoveredAddr, err := recover()
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

func TestEIP712SignBiz_SignForZone(t *testing.T) {
	type args struct {
		priv string
	}
	tests := []struct {
		name          string
		args          args
		wantSignature string
		wantErr       bool
	}{
		{
			name: "test",
			args: args{
				priv: "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
			},
			wantSignature: "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, err := SignForZone(tt.args.priv)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.SignForZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("EIP712SignBiz.SignForZone() = %v, want %v", gotSignature, tt.wantSignature)
			}
		})
	}
}

func TestEIP712SignBiz_SignForOrder(t *testing.T) {
	type args struct {
		priv  string
		order string
	}
	tests := []struct {
		name          string
		args          args
		wantSignature string
		wantErr       bool
	}{
		{
			name: "test remix",
			args: args{
				priv:  "7e5bfb82febc4c2c8529167104271ceec190eafdca277314912eaabdb67c6e5f",
				order: `{"types":{"OrderComponents":[{"name":"offerer","type":"address"},{"name":"zone","type":"address"},{"name":"offer","type":"OfferItem[]"},{"name":"consideration","type":"ConsiderationItem[]"},{"name":"orderType","type":"uint8"},{"name":"startTime","type":"uint256"},{"name":"endTime","type":"uint256"},{"name":"zoneHash","type":"bytes32"},{"name":"salt","type":"uint256"},{"name":"conduitKey","type":"bytes32"},{"name":"counter","type":"uint256"}],"OfferItem":[{"name":"itemType","type":"uint8"},{"name":"token","type":"address"},{"name":"identifierOrCriteria","type":"uint256"},{"name":"startAmount","type":"uint256"},{"name":"endAmount","type":"uint256"}],"ConsiderationItem":[{"name":"itemType","type":"uint8"},{"name":"token","type":"address"},{"name":"identifierOrCriteria","type":"uint256"},{"name":"startAmount","type":"uint256"},{"name":"endAmount","type":"uint256"},{"name":"recipient","type":"address"}],"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}]},"domain":{"name":"Seaport","version":"1.5","chainId":"1","verifyingContract":"0xd8b934580fce35a11b58c6d73adee468a2833fa8"},"primaryType":"OrderComponents","message":{"offerer":"0xab8483f64d9c6d1ecf9b849ae677dd3315835cb2","zone":"0x0000000000000000000000000000000000000000","offer":[{"itemType":"2","token":"0xf8e81d47203a594245e36c48e151709f0c19fbe8","identifierOrCriteria":"1","startAmount":"1","endAmount":"1"}],"consideration":[{"itemType":"0","token":"0x0000000000000000000000000000000000000000","identifierOrCriteria":"0","startAmount":"20000000","endAmount":"20000000","recipient":"0xab8483f64d9c6d1ecf9b849ae677dd3315835cb2"},{"itemType":"0","token":"0x0000000000000000000000000000000000000000","identifierOrCriteria":"0","startAmount":"1000000","endAmount":"1000000","recipient":"0x5b38da6a701c568545dcfcb03fcb875f56beddc4"}],"orderType":"0","startTime":"1291123844","endTime":"1791123844","zoneHash":"0x0000000000000000000000000000000000000000000000000000000000000000","salt":"24446860302761739304752683030156737591518664810215442929816108075358245614181","conduitKey":"0x0000000000000000000000000000000000000000000000000000000000000000","counter":"0"}}`,
			},
			wantSignature: "",
			wantErr:       false,
		},
		{
			name: "test shate Order AS",
			args: args{
				priv:  "2bdd8a43b8a055632f9a8b38d7c9463bfbe8340abbd6634f51bbf000cbe0ca50",
				order: testOrderData,
			},
			wantSignature: "",
			wantErr:       false,
		},
		{
			name: "test shate Offer Qoe",
			args: args{
				priv:  "c52c5d5357b87fbbc73758d7f31e88ad80d7364101851749e501a80150da05ff",
				order: testOfferData,
			},
			wantSignature: "",
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSignature, err := SignForOrder(tt.args.priv, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("EIP712SignBiz.SignForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignature != tt.wantSignature {
				t.Errorf("EIP712SignBiz.SignForOrder() = %v, want %v", gotSignature, tt.wantSignature)
			}
		})
	}
}

// domain: 0xa6cfad1568db403a69402edef3705bea47f638d2dbb1387e8c851f476ec95bdf
// dataHash: 0x194b9090f745afbee7ba7bda60df9a2ecadb100d8daaec3ba457df77150b5d6b
// signDataDigest: 0x91c45bc9491d0f4d0e1e682bafc91ba30cefef5a8f0f6b9d5cc8832f92f4b845

var testOrderData = `
{
    "types": {
        "OrderComponents": [
            {
                "name": "offerer",
                "type": "address"
            },
            {
                "name": "zone",
                "type": "address"
            },
            {
                "name": "offer",
                "type": "OfferItem[]"
            },
            {
                "name": "consideration",
                "type": "ConsiderationItem[]"
            },
            {
                "name": "orderType",
                "type": "uint8"
            },
            {
                "name": "startTime",
                "type": "uint256"
            },
            {
                "name": "endTime",
                "type": "uint256"
            },
            {
                "name": "zoneHash",
                "type": "bytes32"
            },
            {
                "name": "salt",
                "type": "uint256"
            },
            {
                "name": "conduitKey",
                "type": "bytes32"
            },
            {
                "name": "counter",
                "type": "uint256"
            }
        ],
        "OfferItem": [
            {
                "name": "itemType",
                "type": "uint8"
            },
            {
                "name": "token",
                "type": "address"
            },
            {
                "name": "identifierOrCriteria",
                "type": "uint256"
            },
            {
                "name": "startAmount",
                "type": "uint256"
            },
            {
                "name": "endAmount",
                "type": "uint256"
            }
        ],
        "ConsiderationItem": [
            {
                "name": "itemType",
                "type": "uint8"
            },
            {
                "name": "token",
                "type": "address"
            },
            {
                "name": "identifierOrCriteria",
                "type": "uint256"
            },
            {
                "name": "startAmount",
                "type": "uint256"
            },
            {
                "name": "endAmount",
                "type": "uint256"
            },
            {
                "name": "recipient",
                "type": "address"
            }
        ],
        "EIP712Domain": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "version",
                "type": "string"
            },
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            }
        ]
    },
    "domain": {
        "name": "Seaport",
        "version": "1.5",
        "chainId": "2494104990",
        "verifyingContract": "0x74e80c9067f58873e626ff3966aa3abab7a1c6cb"
    },
    "primaryType": "OrderComponents",
    "message": {
        "offerer": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789",
        "zone": "0x23cafcac35dee6f487493f7eea284d8689b8c179",
        "offer": [
            {
                "itemType": "2",
                "token": "0xfb8d0f7d033268b76a5077a5462ae711d1e48a0b",
                "identifierOrCriteria": "13",
                "startAmount": "1",
                "endAmount": "1"
            }
        ],
        "consideration": [
            {
                "itemType": "1",
                "token": "0x42a1e39aefa49290f2b3f9ed688d7cecf86cd6e0",
                "identifierOrCriteria": "0",
                "startAmount": "2000000",
                "endAmount": "2000000",
                "recipient": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789"
            }
        ],
        "orderType": "0",
        "startTime": "1695614310",
        "endTime": "1727005827",
        "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "salt": "1315088691083223237152053736920091101095719954210894886471724671315003257441",
        "conduitKey": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "counter": "0"
    }
}
`

var testOfferData = `
{
    "types": {
        "OrderComponents": [
            {
                "name": "offerer",
                "type": "address"
            },
            {
                "name": "zone",
                "type": "address"
            },
            {
                "name": "offer",
                "type": "OfferItem[]"
            },
            {
                "name": "consideration",
                "type": "ConsiderationItem[]"
            },
            {
                "name": "orderType",
                "type": "uint8"
            },
            {
                "name": "startTime",
                "type": "uint256"
            },
            {
                "name": "endTime",
                "type": "uint256"
            },
            {
                "name": "zoneHash",
                "type": "bytes32"
            },
            {
                "name": "salt",
                "type": "uint256"
            },
            {
                "name": "conduitKey",
                "type": "bytes32"
            },
            {
                "name": "counter",
                "type": "uint256"
            }
        ],
        "OfferItem": [
            {
                "name": "itemType",
                "type": "uint8"
            },
            {
                "name": "token",
                "type": "address"
            },
            {
                "name": "identifierOrCriteria",
                "type": "uint256"
            },
            {
                "name": "startAmount",
                "type": "uint256"
            },
            {
                "name": "endAmount",
                "type": "uint256"
            }
        ],
        "ConsiderationItem": [
            {
                "name": "itemType",
                "type": "uint8"
            },
            {
                "name": "token",
                "type": "address"
            },
            {
                "name": "identifierOrCriteria",
                "type": "uint256"
            },
            {
                "name": "startAmount",
                "type": "uint256"
            },
            {
                "name": "endAmount",
                "type": "uint256"
            },
            {
                "name": "recipient",
                "type": "address"
            }
        ],
        "EIP712Domain": [
            {
                "name": "name",
                "type": "string"
            },
            {
                "name": "version",
                "type": "string"
            },
            {
                "name": "chainId",
                "type": "uint256"
            },
            {
                "name": "verifyingContract",
                "type": "address"
            }
        ]
    },
    "domain": {
        "name": "Seaport",
        "version": "1.5",
        "chainId": "2494104990",
        "verifyingContract": "0x74e80c9067f58873e626ff3966aa3abab7a1c6cb"
    },
    "primaryType": "OrderComponents",
    "message": {
		"offerer": "0x0efa12c664f53f568a9ef3cccdab363096877b6f",
        "zone": "0x23cafcac35dee6f487493f7eea284d8689b8c179",
        "offer": [
            {
                "itemType": "1",
                "token": "0x42a1e39aefa49290f2b3f9ed688d7cecf86cd6e0",
                "identifierOrCriteria": "0",
                "startAmount": "2100000",
                "endAmount": "2100000"
            }
        ],
        "consideration": [
            {
                "itemType": "2",
                "token": "0xfb8d0f7d033268b76a5077a5462ae711d1e48a0b",
                "identifierOrCriteria": "13",
                "startAmount": "1",
                "endAmount": "1",
                "recipient": "0x0efa12c664f53f568a9ef3cccdab363096877b6f"
            }
        ],
        "orderType": "0",
        "startTime": "1695614310",
        "endTime": "1727005827",
        "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "salt": "115336851138071314188779942529410797129178251872168203016652583271119391619243",
        "conduitKey": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "counter": "0"
    }
}
`
