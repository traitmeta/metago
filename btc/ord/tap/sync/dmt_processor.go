package sync

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/traitmeta/metago/btc/ord/common"
	"github.com/traitmeta/metago/btc/ord/tap/dal"
	"github.com/traitmeta/metago/btc/ord/tap/model"
)

type DMTProcessor struct {
	ctx   context.Context
	dao   *dal.Dal
	cache *Cache
}

func NewDMTProcessor(ctx context.Context, db *gorm.DB, cache *Cache) *DMTProcessor {
	return &DMTProcessor{
		ctx:   ctx,
		dao:   dal.NewDal(db),
		cache: cache,
	}
}

func (s *DMTProcessor) ProcessElement(blockHeight int64, txId string, envelopes []Envelope) []model.TapElement {
	var elements []model.TapElement
	for _, envelope := range envelopes {
		insData := envelope.ConvertToInscriptionData()
		if !strings.Contains(insData.ContentType, "text") {
			continue
		}

		content := string(insData.Body)
		element, err := ParseElementFromString(content)
		if err != nil {
			continue
		}

		elements = append(elements, model.TapElement{
			Element:              content,
			ElementInscriptionId: fmt.Sprintf("%si%d", txId, envelope.Offset),
			Name:                 element.Name,
			Pattern:              element.Pattern,
			Field:                element.Field,
			InscriptionHeight:    blockHeight,
		})
	}

	return elements
}

func (s *DMTProcessor) FilterValidElement(elements []model.TapElement) []model.TapElement {
	var validElements []model.TapElement
	var thisBlockInscribeElements = make(map[string]model.TapElement)
	for _, element := range elements {
		noName := strings.ToLower(ElementNoName(element.Pattern, element.Field))
		name := strings.ToLower(element.Name)
		if _, ok := thisBlockInscribeElements[name]; ok {
			continue
		}

		if _, ok := thisBlockInscribeElements[noName]; ok {
			continue
		}

		// 名字重复 = 非法
		if element, _ := s.cache.GetNameToElement(name); element != "" {
			continue
		}

		// 非名字部分重复 = 非法
		if element, _ := s.cache.GetNameToElement(noName); element != "" {
			continue
		}

		thisBlockInscribeElements[name] = element
		thisBlockInscribeElements[noName] = element
		validElements = append(validElements, element)
	}

	return validElements
}

func (s *DMTProcessor) ProcessDeploy(blockHeight, blockTime int64, txId string, envelopes []Envelope) ([]model.TapElementTick, []model.TapActivity, error) {
	var deployTicks []model.TapElementTick
	var deployActivities []model.TapActivity
	for _, envelope := range envelopes {
		dmt, err := EnvelopToDmtOpr(envelope)
		if err != nil && !errors.Is(err, ErrNotTapProtocol) {
			log.WithContext(s.ctx).WithField("inscription_id", fmt.Sprintf("%si%d", txId, envelope.Offset)).
				Debug("DMTIndexer ProcessDeploy EnvelopToDmtOpr fail")
			continue
		}
		if dmt == nil || dmt.Operation != common.DmtDeploy {
			continue
		}

		deployTicks = append(deployTicks, model.TapElementTick{
			ElementInscriptionId: dmt.Element,
			Tick:                 dmt.Ticker,
			TickInscriptionId:    fmt.Sprintf("%si%d", txId, envelope.Offset),
			InscriptionHeight:    blockHeight,
			DeployTime:           blockTime,
		})

		deployActivities = append(deployActivities, model.TapActivity{
			ElementInscriptionId: dmt.Element,
			Type:                 model.DeployType,
			Tick:                 dmt.Ticker,
			Body:                 string(envelope.GetContent()),
			InscriptionHeight:    blockHeight,
			InscriptionId:        fmt.Sprintf("%si%d", txId, envelope.Offset),
		})
	}

	return deployTicks, deployActivities, nil
}

// FilterValidDmtDeploy 缓存中没有，db中也没有就是有效的deploy
func (s *DMTProcessor) FilterValidDmtDeploy(deployTicks []model.TapElementTick, deployActivities []model.TapActivity) ([]model.TapElementTick, []model.TapActivity, error) {
	var validDeployTick []model.TapElementTick
	var validDeployActivities []model.TapActivity
	var thisBlockDeploys = make(map[string]model.TapElementTick)
	for i, deploy := range deployTicks {
		tick := strings.ToLower(deploy.Tick)
		if _, ok := thisBlockDeploys[tick]; ok {
			continue
		}

		detail, _ := s.cache.GetTickDeployDetail(tick)
		if detail != nil {
			// 无效 deploy
			continue
		}

		thisBlockDeploys[tick] = deploy
		dbElemTick, _ := s.dao.GetElementTickByTick(tick, deploy.InscriptionHeight)
		if dbElemTick != nil && strings.EqualFold(dbElemTick.Tick, tick) {
			continue
		}

		validDeployTick = append(validDeployTick, deploy)
		validDeployActivities = append(validDeployActivities, deployActivities[i])
	}

	return validDeployTick, validDeployActivities, nil
}

func (s *DMTProcessor) ProcessMint(blockHeight int64, txId string, envelopes []Envelope) ([]model.TapActivity, error) {
	var mintActivities []model.TapActivity
	for _, envelope := range envelopes {
		dmt, err := EnvelopToDmtOpr(envelope)
		if err != nil && !errors.Is(err, ErrNotTapProtocol) {
			log.WithContext(s.ctx).WithField("inscription_id", fmt.Sprintf("%si%d", txId, envelope.Offset)).
				Debug("DMTIndexer ProcessMint EnvelopToDmtOpr fail")
			continue
		}
		if dmt == nil || dmt.Operation != common.DmtMint {
			continue
		}

		if _, err := strconv.ParseInt(dmt.Block, 10, 64); err != nil {
			continue
		}

		mintActivities = append(mintActivities, model.TapActivity{
			Type:                model.MintType,
			DeployInscriptionId: dmt.Deploy,
			Tick:                dmt.Ticker,
			Body:                string(envelope.GetContent()),
			BlockNumber:         dmt.Block,
			InscriptionHeight:   blockHeight,
			InscriptionId:       fmt.Sprintf("%si%d", txId, envelope.Offset),
		})
	}

	return mintActivities, nil
}

func (s *DMTProcessor) FilterValidDmtMint(mintActivities []model.TapActivity, thisBlockValidDeploy []model.TapElementTick) ([]model.TapActivity, error) {
	var validMintActivities []model.TapActivity
	var thisBlockMintedActivities = make(map[string]model.TapActivity)
	for _, act := range mintActivities {
		tick := strings.ToLower(act.Tick)
		if _, ok := thisBlockMintedActivities[tick+act.BlockNumber]; ok {
			continue
		}

		mintedBlock, _ := s.cache.GetTickMintedBlock(tick, act.BlockNumber)
		if mintedBlock != "" {
			// 无效 mint
			continue
		}

		thisBlockMintedActivities[tick+act.BlockNumber] = act
		// 数据库存在 mint无效
		// TODO 用批量查询替代缓存
		count, _ := s.dao.CountMintActivityWithBlock(tick, act.BlockNumber, act.InscriptionHeight)
		if count > 0 {
			continue
		}

		detail, _ := s.cache.GetTickDeployDetail(tick)
		if detail != nil {
			act.ElementInscriptionId = detail.ElementInscriptionId
			validMintActivities = append(validMintActivities, act)
			continue
		}

		dbElemTick, _ := s.dao.GetElementTick(act.DeployInscriptionId, act.InscriptionHeight)
		if dbElemTick != nil {
			act.ElementInscriptionId = dbElemTick.ElementInscriptionId
			validMintActivities = append(validMintActivities, act)
			continue
		}

		for _, deploy := range thisBlockValidDeploy {
			if strings.EqualFold(deploy.TickInscriptionId, act.DeployInscriptionId) {
				act.ElementInscriptionId = deploy.ElementInscriptionId
				continue
			}
		}

		validMintActivities = append(validMintActivities, act)
	}

	return validMintActivities, nil
}

func EnvelopToDmtOpr(envelope Envelope) (*DmtOpr, error) {
	insData := envelope.ConvertToInscriptionData()
	dmtOpr := &DmtOpr{}
	err := json.Unmarshal(insData.Body, dmtOpr)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(dmtOpr.Protocol, common.TapProtocol) {
		return nil, ErrNotTapProtocol
	}

	return dmtOpr, nil
}
