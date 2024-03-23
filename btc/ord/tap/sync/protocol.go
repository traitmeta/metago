package sync

import (
	"github.com/traitmeta/metago/btc/ord/envelops"
)

type InscriptionData struct {
	ContentType string `json:"content_type"`
	Body        []byte `json:"body"`
	Destination string `json:"destination"`
}

func ConvertToInscriptionData(e envelops.Envelope) InscriptionData {
	return InscriptionData{
		ContentType: e.GetContentType(),
		Body:        e.GetContent(),
		Destination: "",
	}
}
