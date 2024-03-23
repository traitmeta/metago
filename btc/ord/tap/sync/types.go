package sync

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const ElementId = ".element"
const TapCacheId = "tap_cache_id:"

var (
	ErrNotMatchElementPattern = errors.New("not match element pattern")
	ErrNotEndWithElement      = errors.New("not end with element")
	ErrNotTapProtocol         = errors.New("not tap protocol")
)

type TickDetail struct {
	ElementInscriptionId string `json:"element_inscription_id" gorm:"column:element_inscription_id;default:NULL"`
	TickInscriptionId    string `json:"tick_inscription_id" gorm:"column:tick_inscription_id;default:NULL"`
	InscriptionHeight    int64  `json:"inscription_height" gorm:"column:inscription_height;"`
}

type DmtOpr struct {
	Protocol  string `json:"p"`
	Operation string `json:"op"`
	Element   string `json:"elem"`
	Deploy    string `json:"dep"` // mint
	Ticker    string `json:"tick"`
	Project   string `json:"prj"`
	Dimension string `json:"dim"`
	DataTypes string `json:"dt"`
	Id        string `json:"id"`
	Block     string `json:"blk"` // mint
	Amount    string `json:"amt"` // transfer
}

type Element struct {
	Name    string
	Pattern string
	Field   string
}

func (e *Element) NoName() string {
	if e.Pattern != "" {
		return fmt.Sprintf("%s.%s.element", e.Pattern, e.Field)
	}
	return fmt.Sprintf("%s.element", e.Field)
}

func ElementNoName(pattern, field string) string {
	if pattern != "" {
		return fmt.Sprintf("%s.%s.element", pattern, field)
	}
	return fmt.Sprintf("%s.element", field)
}

func ParseElementFromString(content string) (*Element, error) {
	var element = &Element{}
	if !MatchElementPattern(content) {
		return nil, ErrNotMatchElementPattern
	}

	elementSplits := strings.Split(content, ".")
	if !strings.EqualFold(elementSplits[len(elementSplits)-1], "element") {
		return nil, ErrNotEndWithElement
	}

	if len(elementSplits) == 4 {
		element.Pattern = elementSplits[1]
		element.Field = elementSplits[2]
	} else {
		element.Field = elementSplits[1]
	}
	element.Name = elementSplits[0]

	return element, nil
}

type CachedKeys struct {
	NameKeys   []string
	NoNameKeys []string
	DeployKeys []string
	MintKeys   []string
}
