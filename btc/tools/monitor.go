package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const url = "https://mempool.space/api/address/%s/txs/summary"

type TxSummary struct {
	TxId   string `json:"txid"`
	Height int64  `json:"height"`
	Value  int64  `json:"value"`
	Time   int64  `json:"time"`
}

func GetBalanceToday(address string) (int64, error) {
	resp, err := http.Get(fmt.Sprintf(url, address))
	if err != nil {
		return 0, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return 0, errors.New("request err")
	}

	var txsSummary []TxSummary
	err = json.Unmarshal(body, &txsSummary)
	if err != nil {
		return 0, err
	}

	var amount int64
	for _, v := range txsSummary {
		ts := time.Now().AddDate(0, 0, -1)
		timestamp := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, ts.Location()).Unix()
		if v.Time > timestamp {
			amount += v.Value
		}
	}
	return amount, nil
}
