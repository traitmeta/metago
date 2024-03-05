package tap

/*
DMTDeploy

	{
		"p": "tap",
		"op": "dmt-deploy",
		"elem": "<inscriptionID>",
		"tick": "<name>",
		"prj": "<0.bitmap inscriptionID>",
		"dim": "h | v | d | a",
		"dt": "h | n | x | s | b",
		"id": "<content inscriptionID>"
	}
*/
type DMTDeploy struct {
	Protocol  string `json:"p"`    // Protocol: TAP
	Operation string `json:"op"`   // Operation: Event (dmt-deploy, token-mint, token-transfer)
	Element   string `json:"elem"` // Element: Reference to the .element inscriptionID
	Ticker    string `json:"tick"` // Ticker: 3 and 5 to 32 (UTF16)
	Project   string `json:"prj"`  // Project: Reference to existing project's (i.e. 0.bitmap) inscriptionID
	Dimension string `json:"dim"`  // Dimension: Only required if you're recognizing a pattern in a given field. (horizontal, vertical, diagonal, add) 'add' represents the sum of 'h', 'v', and 'd', counts of the given pattern.
	DataTypes string `json:"dt"`   // Data Types: "h" hex, "n" numeric, "x" unix time, "s" string, "b" boolean. If the "dt" field is left out, the data type set within the block data is used by default. ⚠️ If the element you choose has a pattern, "dt" in your deploy inscription is required. ⚠️
	Id        string `json:"id"`   // ID: this is the inscription id of your UNAT content. (generative script art, 3d model, app)
}

func (d *DMTDeploy) Type() string {
	return d.Operation
}

/*
DMTMint

	{
		"p": "tap",
		"op": "dmt-mint",
		"dep": "<inscriptionID>",
		"tick": "<name>",
		"blk": "#"
	}
*/
type DMTMint struct {
	Protocol  string `json:"p"`    // Protocol: TAP
	Operation string `json:"op"`   // Operation: Event (dmt-deploy, token-mint, token-transfer)
	Deploy    string `json:"dep"`  // Deploy: Reference to the dmt-deploy inscriptionID. Note: "dep" will be optional from a block height TBA)
	Ticker    string `json:"tick"` // Ticker: Matches the ticker in the dmt-deploy inscription
	Block     string `json:"blk"`  // Block: Inscribe a block number of your choosing. ⚠️ First is first.
}

/*
DMTTransfer

	{
		"p": "tap",
		"op": "token-transfer",
		"tick": "dmt-<name>",
		"amt": "#"
	}
*/
type DMTTransfer struct {
	Protocol  string `json:"p"`    // Protocol: TAP
	Operation string `json:"op"`   // Operation: Event (dmt-deploy, token-mint, token-transfer)
	Ticker    string `json:"tick"` // Ticker: Matches the ticker in the token-mint inscription
	Amount    string `json:"amt"`  // Amount: Chose an amount of NATs to transfer.
}
