package ord

type InscribeTool interface {
	GetPayAddrAndFee(request *InscriptionRequest) (payAddr, payAddrPK string, inscFee, serviceFee, minerFee int64, err error)
	Inscribe(commitTxHash string, actualMiddlePrevOutputFee int64, payAddrPK string, request *InscriptionRequest) (ctxTxData *CtxTxData, err error)
}
