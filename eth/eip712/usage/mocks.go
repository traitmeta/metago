package usage

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

var testZoneSigner = `
	{
		"types": {
		  "SignedOrder": [
			{ "name": "fulfiller", "type": "address" },
			{ "name": "expiration", "type": "uint64" },
			{ "name": "orderHash", "type": "bytes32" },
			{ "name": "context", "type": "bytes" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "SignedZone",
		  "version": "1.0",
		  "chainId": "1",
		  "verifyingContract": "0xd182d0a388f4923c478395dfb4ea889e55013967"
		},
		"primaryType": "SignedOrder",
		"message": {
		  "fulfiller": "0x4b20993bc481177ec7e8f571cecae8a9e22c02db",
		  "expiration": "1751641513",
		  "orderHash": "0x846c2aa6277c50980556cccc77d2c9bcde1258b00228ce062da733268802fa01",
		  "context": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}
	  }
	`

var testOrderNft = `
	{
		"types": {
		  "OrderComponents": [
			{ "name": "offerer", "type": "address" },
			{ "name": "zone", "type": "address" },
			{ "name": "offer", "type": "OfferItem[]" },
			{ "name": "consideration", "type": "ConsiderationItem[]" },
			{ "name": "orderType", "type": "uint8" },
			{ "name": "startTime", "type": "uint256" },
			{ "name": "endTime", "type": "uint256" },
			{ "name": "zoneHash", "type": "bytes32" },
			{ "name": "salt", "type": "uint256" },
			{ "name": "conduitKey", "type": "bytes32" },
			{ "name": "counter", "type": "uint256" }
		  ],
		  "OfferItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" }
		  ],
		  "ConsiderationItem": [
			{ "name": "itemType", "type": "uint8" },
			{ "name": "token", "type": "address" },
			{ "name": "identifierOrCriteria", "type": "uint256" },
			{ "name": "startAmount", "type": "uint256" },
			{ "name": "endAmount", "type": "uint256" },
			{ "name": "recipient", "type": "address" }
		  ],
		  "EIP712Domain": [
			{ "name": "name", "type": "string" },
			{ "name": "version", "type": "string" },
			{ "name": "chainId", "type": "uint256" },
			{ "name": "verifyingContract", "type": "address" }
		  ]
		},
		"domain": {
		  "name": "Seaport",
		  "version": "1.5",
		  "chainId": "2494104990",
		  "verifyingContract": "0x9afa6139c383e7a3796131a43dcb86caa9178170"
		},
		"primaryType": "OrderComponents",
		"message": {
		  "offerer": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789",
		  "zone": "0x0000000000000000000000000000000000000000",
		  "offer": [
			{
			  "itemType": "2",
			  "token": "0xfb8d0f7d033268b76a5077a5462ae711d1e48a0b",
			  "identifierOrCriteria": "4",
			  "startAmount": "1",
			  "endAmount": "1"
			}
		  ],
		  "consideration": [
			{
			  "itemType": "0",
			  "token": "0x0000000000000000000000000000000000000000",
			  "identifierOrCriteria": "0",
			  "startAmount": "20000000",
			  "endAmount": "20000000",
			  "recipient": "0x33a4f229bd34ea7783302c99ffd6e26324bd2789"
			},
			{
			  "itemType": "0",
			  "token": "0x0000000000000000000000000000000000000000",
			  "identifierOrCriteria": "0",
			  "startAmount": "1000000",
			  "endAmount": "1000000",
			  "recipient": "0x95b467b0d33c34d5bc2ab3fb005cf9aca4033f00"
			}
		  ],
		  "orderType": "0",
		  "startTime": "1591123844",
		  "endTime": "1791123844",
		  "zoneHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "salt": "24446860302761739304752683030156737591518664810215442929816108075358245610000",
		  "conduitKey": "0x0000000000000000000000000000000000000000000000000000000000000000",
		  "counter": "0"
		}
	  }
	`
