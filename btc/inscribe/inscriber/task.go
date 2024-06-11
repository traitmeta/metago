package inscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type SendTask struct {
	client  BTCBaseClient
	who     string
	runes   string
	orderId string
	count   int64
}

func (s *SendTask) LoopSendTxs(sr *SendResult) {
	for {
		if v, ok := sr.TxsStatus[sr.MiddleTx.WireTx.TxHash().String()]; !ok || !v {
			middleTxHash, err := s.client.SendRawTransaction(sr.MiddleTx.WireTx)
			if err != nil {
				log.Printf("send middle tx error %s \n", middleTxHash.String())
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("middleTxHash %s \n", middleTxHash.String())
			sr.TxsStatus[middleTxHash.String()] = true
			break
		}

		break
	}

	for i := 0; i < len(sr.RevealTxs); {
		if v, ok := sr.TxsStatus[sr.RevealTxs[i].WireTx.TxHash().String()]; !ok || !v {
			revealTxHash, err := s.client.SendRawTransaction(sr.RevealTxs[i].WireTx)
			if err != nil {
				log.Printf("revealTxHash %d %s , err: %v \n", i, revealTxHash.String(), err)
				time.Sleep(5 * time.Second)
				continue
			}

			sr.TxsStatus[revealTxHash.String()] = true
		}

		i++
	}

	s.SaveToFile(sr)
}

func (s *SendTask) SaveToFile(sr *SendResult) {
	key := fmt.Sprintf("%s:%s:%d:%s.json", s.who, s.runes, s.count, s.orderId)
	bytes, err := json.Marshal(sr)
	if err != nil {
		log.Printf("save to file failed %s, content : %s , err: %v \n", key, string(bytes), err)
		return
	}

	if err := os.WriteFile(key, bytes, 0666); err != nil {
		log.Printf("save to file failed %s, content : %s , err: %v \n", key, string(bytes), err)
	}
}
