package common

const (
	BRC20    = "brc20"
	Tap      = "tap"
	TapBlock = "block"
)

const (
	RollBackBlockNumber = 3
)

const (
	Brc20OpDeploy           = "deploy"
	Brc20OpMint             = "mint"
	Brc20OpInscribeTransfer = "transfer"
	Brc20OpBalanceTransfer  = "balance-transfer"
)

const (
	ProcessOrdInscriptionNumKey    = "Process_Ord_Inscription_Num"
	OrdIndexBlockNumKey            = "Ord_Index_Block_Num"
	OrdInscriptionDetailUpdatedNum = "Ord_Inscription_Detail_Updated_Num_20230712"
	ProcessOrdIndexBlockNumKey     = "Process_Ord_Index_Block_Num"
	UTXOInscriptionLocation        = "UTXO_Inscription_Location"
	UTXOThisOutputAddress          = "UTXO_This_Output_Address"
)

const (
	OpDeploy           = 1
	OpMint             = 2
	OpInscribeTransfer = 3
	OpBalanceTransfer  = 4
)
const PushForwardBlock = 10
const SearchBlockRange = 5
const HotMintBlockNum = 5
const (
	TextType = "text"
	HTML     = "html"
	Image    = "image"
)

const (
	ARC20  = "arc-20"
	BRC420 = "brc-420"
	TapDmt = "dmt-tap"
)

const (
	DeployOp   = "deploy"
	MintOp     = "mint"
	TransferOp = "transfer"
)

const (
	Atomicals      = "Atomicals"
	OrdinalsType   = "Ordinals"
	DomainProtocol = "Domain"
)

const BtcSupply = 21000000

const (
	Brc420MintPattern          = `^/content/[0-9a-z]{64}i\d$`
	Brc420MintContentPrefixLen = 9
	DomainPattern              = `^[a-zA-Z0-9]+\.[a-zA-Z0-9]+$`
)

const (
	TapProtocol    = "tap"
	ElementPattern = `^[^/\.\[\]{}:;"']+\..*\.*\d+.element$`
)

const (
	DmtDeploy             = "dmt-deploy"
	DmtMint               = "dmt-mint"
	DmtTransfer           = "protocol-transfer"
	DmtTransferTickPrefix = "dmt-"
)

const HotMintBlockChannel = "hotmint-block"
