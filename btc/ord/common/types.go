package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type InscriptionBRC20Data struct {
	IsTransfer        bool
	BlockHash         string
	BlockHeight       int64 // Height of NFT show in block onCreate
	BlockTime         int64
	TxHash            string `json:"-"`
	TxIdx             uint32 `json:"-"`
	InputIndex        uint32
	UtxoIndex         uint32
	Operation         string
	Satoshi           int64 `json:"-"`
	From              string
	TO                string
	InscriptionNumber int64
	InscriptionID     string
	Inscription       BrcInscription
	Balance           decimal.Decimal
	Valid             bool
	InvalidReason     int
	GasFee            int64
}

type CachedBrcInscription struct {
	BrcInscription
	ToAddr  string `json:"to_addr"`
	Satoshi int64  `json:"satoshi"`
	Number  int64  `json:"number"`
	ID      string `json:"id"`
}

func (bi *CachedBrcInscription) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, value := range raw {
		switch strings.ToLower(key) {
		case "p":
			if key != "p" {
				return errors.New(fmt.Sprintf("p filed:%s", key))
			}
			if err := json.Unmarshal(value, &bi.Proto); err != nil {
				return err
			}
		case "op":
			if key != "op" {
				return errors.New(fmt.Sprintf("op filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.Operation); err != nil {
				return err
			}
		case "tick":
			if key != "tick" {
				return errors.New(fmt.Sprintf("tick filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Tick); err != nil {
				return err
			}
		case "amt":

			if key != "amt" {
				return errors.New(fmt.Sprintf("amt filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Amount); err != nil {
				return err
			}
		case "max":
			if key != "max" {
				return errors.New(fmt.Sprintf("max filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Max); err != nil {
				return err
			}
		case "lim":
			if key != "lim" {
				return errors.New(fmt.Sprintf("lim filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Limit); err != nil {
				return err
			}
		case "to":
			if key != "to" {
				return errors.New(fmt.Sprintf("to filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20To); err != nil {
				return err
			}
		case "dec":
			if key != "dec" {
				return errors.New(fmt.Sprintf("dec filed:%s", key))
			}
			if string(value) == "" {
				return errors.New(fmt.Sprintf("dec filed value:%s", value))
			}

			if err := json.Unmarshal(value, &bi.BRC20Decimal); err != nil {
				return err
			}
		case "fee":
			if err := json.Unmarshal(value, &bi.BRC20Fee); err != nil {
				return err
			}
		case "to_addr":
			if err := json.Unmarshal(value, &bi.ToAddr); err != nil {
				return err
			}
		case "satoshi":
			if err := json.Unmarshal(value, &bi.Satoshi); err != nil {
				return err
			}
		case "number":
			if err := json.Unmarshal(value, &bi.Number); err != nil {
				return err
			}
		case "id":
			if err := json.Unmarshal(value, &bi.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

type BrcInscription struct {
	Proto        string `json:"p,omitempty"`
	Operation    string `json:"op,omitempty"`
	BRC20Tick    string `json:"tick,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"`
	BRC20To      string `json:"to,omitempty"`
	BRC20Decimal string `json:"dec,omitempty"`
	BRC20Fee     string `json:"fee,omitempty"`
}

func (bi *BrcInscription) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for key, value := range raw {
		switch strings.ToLower(key) {
		case "p":
			if key != "p" {
				return errors.New(fmt.Sprintf("p filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.Proto); err != nil {
				return err
			}
		case "op":
			if key != "op" {
				return errors.New(fmt.Sprintf("op filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.Operation); err != nil {
				return err
			}
		case "tick":
			if key != "tick" {
				return errors.New(fmt.Sprintf("tick filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Tick); err != nil {
				return err
			}
		case "amt":
			if key != "amt" {
				return errors.New(fmt.Sprintf("amt filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Amount); err != nil {
				return err
			}
		case "max":
			if key != "max" {
				return errors.New(fmt.Sprintf("max filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Max); err != nil {
				return err
			}
		case "lim":
			if key != "lim" {
				return errors.New(fmt.Sprintf("lim filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Limit); err != nil {
				return err
			}
		case "to":
			if key != "to" {
				return errors.New(fmt.Sprintf("to filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20To); err != nil {
				return err
			}
		case "dec":
			if key != "dec" {
				return errors.New(fmt.Sprintf("dec filed:%s", key))
			}
			if string(value) == "" {
				return errors.New(fmt.Sprintf("dec filed value:%s", value))
			}
			if err := json.Unmarshal(value, &bi.BRC20Decimal); err != nil {
				return err
			}
		case "fee":
			if key != "fee" {
				return errors.New(fmt.Sprintf("fee filed:%s", key))
			}

			if err := json.Unmarshal(value, &bi.BRC20Fee); err != nil {
				return err
			}
		}
	}

	return nil
}
