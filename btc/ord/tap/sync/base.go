package sync

type BaseSync struct {
	ch chan bool
}

func (b *BaseSync) Receive() <-chan bool {
	return b.ch
}

func (b *BaseSync) Send() {
	b.ch <- true
}

type Sync interface {
	Start()
	SyncBlock() error
	RollBack(int64, int64) error
}
