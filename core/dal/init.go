package dal

func Init() {
	InitBlockDal()
	InitTransactionDal()
	InitEventDal()
}
