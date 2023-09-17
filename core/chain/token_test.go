package chain

import (
	"math/big"
	"reflect"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/traitmeta/metago/core/common"
	"github.com/traitmeta/metago/core/models"
)

func Test_doParse(t *testing.T) {
	type args struct {
		logs []models.Event
		acc  TokenTransfers
	}
	tests := []struct {
		name    string
		args    args
		want    TokenTransfers
		wantErr bool
	}{
		{
			name: "test parse tokens of batch 1155",
			args: args{
				logs: []models.Event{
					{
						Address:     "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
						FirstTopic:  "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb",
						SecondTopic: "0x0000000000000000000000008b736035bbda71825e0219f5fe4dfb22c35fbddc",
						ThirdTopic:  "0x0000000000000000000000000000000000000000000000000000000000000000",
						FourthTopic: "0x0000000000000000000000004D8834907CEE521D08A3fC77478A0fbc7c94FD0a",
						Data:        "0x00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000000f00000000000000000000000000000000000000000000000000000000000000840000000000000000000000000000000000000000000000000000000000000083000000000000000000000000000000000000000000000000000000000000007300000000000000000000000000000000000000000000000000000000000000740000000000000000000000000000000000000000000000000000000000000077000000000000000000000000000000000000000000000000000000000000006b000000000000000000000000000000000000000000000000000000000000007c00000000000000000000000000000000000000000000000000000000000000700000000000000000000000000000000000000000000000000000000000000072000000000000000000000000000000000000000000000000000000000000006a00000000000000000000000000000000000000000000000000000000000000750000000000000000000000000000000000000000000000000000000000000069000000000000000000000000000000000000000000000000000000000000006d00000000000000000000000000000000000000000000000000000000000000870000000000000000000000000000000000000000000000000000000000000086000000000000000000000000000000000000000000000000000000000000000f000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},
				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC1155,
						ContractAddress: "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0x0000000000000000000000000000000000000000",
						ToAddress:            "0x4D8834907CEE521D08A3fC77478A0fbc7c94FD0a",
						TokenContractAddress: "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						Amounts: []*big.Int{
							big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(2), big.NewInt(2),
							big.NewInt(1), big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1),
						},
						TokenIds: []*big.Int{
							big.NewInt(132), big.NewInt(131), big.NewInt(115), big.NewInt(116), big.NewInt(119), big.NewInt(107), big.NewInt(124),
							big.NewInt(112), big.NewInt(114), big.NewInt(106), big.NewInt(117), big.NewInt(105), big.NewInt(109), big.NewInt(135), big.NewInt(134),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse tokens signle 1155",
			args: args{
				logs: []models.Event{
					{
						Address:     "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
						FirstTopic:  common.ERC1155SingleTransferSignature,
						SecondTopic: "0x0000000000000000000000008b736035bbda71825e0219f5fe4dfb22c35fbddc",
						ThirdTopic:  "0x0000000000000000000000000000000000000000000000000000000000000000",
						FourthTopic: "0x0000000000000000000000004D8834907CEE521D08A3fC77478A0fbc7c94FD0a",
						Data:        "0x00000000000000000000000000000000000000000000000000000000000000460000000000000000000000000000000000000000000000000000000000000001",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},
				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC1155,
						ContractAddress: "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0x0000000000000000000000000000000000000000",
						ToAddress:            "0x4D8834907CEE521D08A3fC77478A0fbc7c94FD0a",
						TokenContractAddress: "0xbCa5858dfd00cEa2eb85e2AB678a5867a18A24c4",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						Amount:               big.NewInt(1),
						TokenId:              big.NewInt(70),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse token 721",
			args: args{
				logs: []models.Event{
					{
						Address:     "0xC6CA7be41Ba10a3645988B77a523231666540b82",
						FirstTopic:  "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						SecondTopic: "0x000000000000000000000000c904a40ed8656ea828f13d3720bcb1f8aff46098",
						ThirdTopic:  "0x0000000000000000000000006c8f9a46294f7e279f23cf2d7900e07f98be0aee",
						FourthTopic: "0x0000000000000000000000000000000000000000000000000000000000045ece",
						Data:        "0x",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},
				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC721,
						ContractAddress: "0xC6CA7be41Ba10a3645988B77a523231666540b82",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0xc904a40eD8656EA828F13D3720bcB1F8Aff46098",
						ToAddress:            "0x6C8f9A46294f7e279F23cF2D7900E07F98be0Aee",
						TokenContractAddress: "0xC6CA7be41Ba10a3645988B77a523231666540b82",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						TokenId:              big.NewInt(286414),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse token 20",
			args: args{
				logs: []models.Event{
					{
						Address:     "0xC6CA7be41Ba10a3645988B77a523231666540b82",
						FirstTopic:  "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
						SecondTopic: "0x00000000000000000000000024d1ddca687cb784572533a97575044444444444",
						ThirdTopic:  "0x000000000000000000000000a09bd67169286f91adf433e68c63a096cf70ad98",
						FourthTopic: "",
						Data:        "0x00000000000000000000000000000000000000000000000000000000000a2c2a",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},

				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC20,
						ContractAddress: "0xC6CA7be41Ba10a3645988B77a523231666540b82",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0x24d1DDca687Cb784572533A97575044444444444",
						ToAddress:            "0xa09bd67169286F91Adf433e68C63A096cf70Ad98",
						TokenContractAddress: "0xC6CA7be41Ba10a3645988B77a523231666540b82",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						Amount:               big.NewInt(666666),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse token WETH Despoit",
			args: args{
				logs: []models.Event{
					{
						Address:     "0x4200000000000000000000000000000000000006",
						FirstTopic:  "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
						SecondTopic: "0x0000000000000000000000002d6937030cc4f1df9c04848554e73be898e8098b",
						ThirdTopic:  "",
						FourthTopic: "",
						Data:        "0x00000000000000000000000000000000000000000000000000038d7ea4c68000",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},

				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC20,
						ContractAddress: "0x4200000000000000000000000000000000000006",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0x0000000000000000000000000000000000000000",
						ToAddress:            "0x2D6937030Cc4F1Df9c04848554e73be898E8098b",
						TokenContractAddress: "0x4200000000000000000000000000000000000006",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						Amount:               big.NewInt(1000000000000000),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse token WETH Withdraw",
			args: args{
				logs: []models.Event{
					{
						Address:     "0x4200000000000000000000000000000000000006",
						FirstTopic:  "0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65",
						SecondTopic: "0x000000000000000000000000e71273af9e573d68323adf96cccb90a47c68d3c7",
						ThirdTopic:  "",
						FourthTopic: "",
						Data:        "0x00000000000000000000000000000000000000000000000008a491333094b012",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},
				acc: TokenTransfers{},
			},
			want: TokenTransfers{
				Tokens: []models.Token{
					{
						Type:            common.ERC20,
						ContractAddress: "0x4200000000000000000000000000000000000006",
					},
				},
				TokenTransfers: []models.TokenTransfer{
					{
						TransactionHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						LogIndex:             0,
						FromAddress:          "0xE71273Af9e573D68323AdF96cCcb90a47C68d3C7",
						ToAddress:            "0x0000000000000000000000000000000000000000",
						TokenContractAddress: "0x4200000000000000000000000000000000000006",
						BlockNumber:          10010,
						BlockHash:            "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						Amount:               big.NewInt(0).SetBytes(ethcommon.FromHex("0x00000000000000000000000000000000000000000000000008a491333094b012")),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test parse token unkown",
			args: args{
				logs: []models.Event{
					{
						Address:     "0x4200000000000000000000000000000000000006",
						FirstTopic:  "0x6fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65",
						SecondTopic: "0x000000000000000000000000e71273af9e573d68323adf96cccb90a47c68d3c7",
						ThirdTopic:  "",
						FourthTopic: "",
						Data:        "0x00000000000000000000000000000000000000000000000008a491333094b012",
						BlockNumber: 10010,
						TxHash:      "0xdc10d88baa62afce2005b3423875a7c6c5a73003e0513e505025a56258870020",
						TxIndex:     0,
						BlockHash:   "0x350479050cc11e6cf26a65d3b43dfd1a68194eeb1f6128d74901447b83770ad3",
						LogIndex:    0,
						Removed:     false,
					},
				},
				acc: TokenTransfers{},
			},
			want:    TokenTransfers{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := doParse(tt.args.logs, tt.args.acc)
			if (err != nil) != tt.wantErr {
				t.Errorf("doParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("doParse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
