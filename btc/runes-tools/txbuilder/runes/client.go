package runes

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

var (
	encipherPath   = "%s/rune/encipher"
	commitmentPath = "%s/rune/commitment/%s"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint: endpoint}
}

type EncipherResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (c *Client) Encipher(stone RuneStone) (string, error) {
	url := fmt.Sprintf(encipherPath, c.endpoint)
	stoneBytes, err := json.Marshal(stone)
	if err != nil {
		return "", errors.Wrap(err, "cannot marshal runestone")
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(stoneBytes))
	if err != nil {
		return "", errors.Wrap(err, "post encipher from rune service failed")
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("status err:" + response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("read body:" + err.Error())
	}

	var scriptKey EncipherResponse
	err = json.Unmarshal(body, &scriptKey)
	if err != nil {
		return "", errors.New("unmarshal body:" + err.Error())
	}

	if scriptKey.Code != http.StatusOK {
		return "", errors.New("encipher failed:" + scriptKey.Message)
	}

	return scriptKey.Data, nil
}

type CommitmentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Rune       string `json:"rune"`
		Spacers    int    `json:"spacers"`
		Commitment string `json:"commitment"`
	} `json:"data"`
}

func (c *Client) Commitment(name string) (string, error) {
	url := fmt.Sprintf(commitmentPath, c.endpoint, name)
	response, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "post encipher from rune service failed")
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("status err:" + response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("read body:" + err.Error())
	}

	var resp CommitmentResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", errors.New("unmarshal body:" + err.Error())
	}

	if resp.Code != http.StatusOK {
		return "", errors.New("encipher failed:" + resp.Message)
	}

	return resp.Data.Commitment, nil
}

func CreateRuneStoneOutput(cli *Client, data RuneStone) (*wire.TxOut, error) {
	encipher, err := cli.Encipher(data)
	if err != nil {
		return nil, err
	}

	script, err := hex.DecodeString(encipher)
	if err != nil {
		return nil, err
	}

	output := wire.NewTxOut(0, script)

	return output, nil
}
