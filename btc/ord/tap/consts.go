package tap

// Data Types
const (
	Hex      = "h"
	Numeric  = "n"
	UnixTime = "x"
	String   = "s"
	Boolean  = "b"
)

// Dimension
const (
	Horizontal = "h"
	Vertical   = "v"
	Diagonal   = "d"
	Add        = "a"
)

// InvalidChar character in name
var InvalidChar = []rune{'/', '.', '[', ']', '{', '}', ':', ';', '"', '\'', ' '}

type Field struct {
	FieldNum    int    `json:"field_num"`
	FieldName   string `json:"field_name"`
	FieldGroup  string `json:"field_group"`
	DefaultType string `json:"default_type"`
}

var FieldMap = map[int]Field{
	0: {
		FieldNum:    0,
		FieldName:   "block_hash",
		FieldGroup:  "block",
		DefaultType: Hex,
	},
	1: {
		FieldNum:    1,
		FieldName:   "size",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	2: {
		FieldNum:    2,
		FieldName:   "strippedsize",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	3: {
		FieldNum:    3,
		FieldName:   "weight",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	4: {
		FieldNum:    4,
		FieldName:   "height",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	5: {
		FieldNum:    5,
		FieldName:   "version",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	6: {
		FieldNum:    6,
		FieldName:   "versionHex",
		FieldGroup:  "block",
		DefaultType: Hex,
	},
	7: {
		FieldNum:    7,
		FieldName:   "merkleroot",
		FieldGroup:  "block",
		DefaultType: Hex,
	},
	8: {
		FieldNum:    8,
		FieldName:   "time",
		FieldGroup:  "block",
		DefaultType: UnixTime,
	},
	9: {
		FieldNum:    9,
		FieldName:   "mediantime",
		FieldGroup:  "block",
		DefaultType: UnixTime,
	},
	10: {
		FieldNum:    10,
		FieldName:   "nonce",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	11: {
		FieldNum:    11,
		FieldName:   "bits",
		FieldGroup:  "block",
		DefaultType: Hex,
	},
	12: {
		FieldNum:    12,
		FieldName:   "difficulty",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	13: {
		FieldNum:    13,
		FieldName:   "chainwork",
		FieldGroup:  "block",
		DefaultType: Hex,
	},
	14: {
		FieldNum:    14,
		FieldName:   "nTx",
		FieldGroup:  "block",
		DefaultType: Numeric,
	},
	15: {
		FieldNum:    15,
		FieldName:   Hex,
		FieldGroup:  "transaction",
		DefaultType: Hex,
	},

	16: { // "txid" : Hex,
		FieldNum:    16,
		FieldName:   "txid",
		FieldGroup:  "transaction",
		DefaultType: Hex,
	},
	17: { //"tx_hash" : Hex,
		FieldNum:    17,
		FieldName:   "tx_hash",
		FieldGroup:  "transaction",
		DefaultType: Hex,
	},
	18: { //"size" : n,
		FieldNum:    18,
		FieldName:   "size",
		FieldGroup:  "transaction",
		DefaultType: Numeric,
	},
	19: { //"vsize" : n,
		FieldNum:    19,
		FieldName:   "vsize",
		FieldGroup:  "transaction",
		DefaultType: Numeric,
	},

	20: { //"weight" : n,
		FieldNum:    20,
		FieldName:   "weight",
		FieldGroup:  "transaction",
		DefaultType: Numeric,
	},
	21: { //"version" : n,
		FieldNum:    21,
		FieldName:   "version",
		FieldGroup:  "transaction",
		DefaultType: Numeric,
	},
	22: { //"locktime" : xxx,
		FieldNum:    22,
		FieldName:   "locktime",
		FieldGroup:  "transaction",
		DefaultType: UnixTime,
	},
	23: { //"blocktime" : xxx,
		FieldNum:    23,
		FieldName:   "blocktime",
		FieldGroup:  "transaction",
		DefaultType: UnixTime,
	},
	24: { //"asm" : String,
		FieldNum:    24,
		FieldName:   "asm",
		FieldGroup:  "input",
		DefaultType: String,
	},
	25: { //"hex" : "hex"
		FieldNum:    25,
		FieldName:   Hex,
		FieldGroup:  "input",
		DefaultType: Hex,
	},
	26: { // "sequence" : n,
		FieldNum:    26,
		FieldName:   "sequence",
		FieldGroup:  "input",
		DefaultType: Numeric,
	},
	27: { //"txinwitness" : Hex,
		FieldNum:    27,
		FieldName:   "txinwitness",
		FieldGroup:  "input",
		DefaultType: Hex,
	},

	28: { //"value" : n,
		FieldNum:    28,
		FieldName:   "value",
		FieldGroup:  "input", // output?
		DefaultType: Numeric,
	},
	29: { //"n" : n,
		FieldNum:    29,
		FieldName:   Numeric,
		FieldGroup:  "input", // output?
		DefaultType: Numeric,
	},
	30: { //"asm" : String,
		FieldNum:    30,
		FieldName:   "asm",
		FieldGroup:  "output",
		DefaultType: String,
	},
	31: { //"hex" : String,
		FieldNum:    31,
		FieldName:   Hex,
		FieldGroup:  "output",
		DefaultType: String,
	},
	32: { //"reqSigs" : n,
		FieldNum:    32,
		FieldName:   "reqSigs",
		FieldGroup:  "output",
		DefaultType: Numeric,
	},
	33: { //"type" : String,
		FieldNum:    33,
		FieldName:   "type",
		FieldGroup:  "output",
		DefaultType: String,
	},
	34: { //"witness":boolean,
		FieldNum:    34,
		FieldName:   "witness",
		FieldGroup:  "extras",
		DefaultType: Boolean,
	},
	35: { //"btc_fee": n,
		FieldNum:    35,
		FieldName:   "btc_fee",
		FieldGroup:  "extras",
		DefaultType: Numeric,
	},
	36: { // "is_coinbase": boolean,
		FieldNum:    36,
		FieldName:   "is_coinbase",
		FieldGroup:  "extras",
		DefaultType: Boolean,
	},
	37: { //"coinbase":"hex"
		FieldNum:    37,
		FieldName:   "coinbase",
		FieldGroup:  "extras",
		DefaultType: Hex,
	},
}
