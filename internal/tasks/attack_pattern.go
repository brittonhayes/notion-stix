package tasks

import (
	"context"
	"encoding/json"

	"github.com/TcM1911/stix2"
	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/dstotijn/go-notion"
	"github.com/hibiken/asynq"
)

const (
	TypeAttackPatternsPageCreate = "attack_patterns:create_page"
)

type CreateAttackPatternPagePayload struct {
	ParentPageID  string               `json:"parent_page_id,omitempty"`
	AttackPattern *stix2.AttackPattern `json:"attack_pattern,omitempty"`
	NotionClient  *notion.Client       `json:"notion_client,omitempty"`
}

func NewCreateAttackPatternsPageTask(ctx context.Context, client *notion.Client, parentPageId string, attackPattern *stix2.AttackPattern) (*asynq.Task, error) {
	payload, err := json.Marshal(CreateAttackPatternPagePayload{
		ParentPageID:  parentPageId,
		AttackPattern: attackPattern,
		NotionClient:  client,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAttackPatternsPageCreate, payload), nil
}

var _ asynq.Handler = (*AttackPatternProcessor)(nil)

type AttackPatternProcessor struct {
	repo notionstix.Repository
}

func NewAttackPatternProcessor(repo notionstix.Repository) *AttackPatternProcessor {
	return &AttackPatternProcessor{
		repo: repo,
	}
}

func (p *AttackPatternProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload CreateAttackPatternPagePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return asynq.SkipRetry
	}

	_, err := p.repo.CreateAttackPatternPage(ctx, payload.NotionClient, payload.ParentPageID, payload.AttackPattern)
	if err != nil {
		return asynq.SkipRetry
	}

	return nil
}
