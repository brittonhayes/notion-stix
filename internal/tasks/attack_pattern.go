package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TcM1911/stix2"
	notionstix "github.com/brittonhayes/notion-stix"
	"github.com/dstotijn/go-notion"
	"github.com/hibiken/asynq"
)

const (
	TypeAttackPatternsPageCreate = "attack_patterns:create_page"
)

var _ asynq.Handler = (*AttackPattern)(nil)

type AttackPattern struct {
	client *notion.Client
	repo   notionstix.Repository
}

type CreateAttackPatternPagePayload struct {
	ParentPageID  string
	AttackPattern *stix2.AttackPattern
}

func NewCreateAttackPatternsPageTask(ctx context.Context, client *notion.Client, properties CreateAttackPatternPagePayload) (*asynq.Task, error) {
	payload, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeAttackPatternsPageCreate, payload), nil
}

func (p *AttackPattern) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload CreateAttackPatternPagePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	_, err := p.repo.CreateAttackPatternPage(ctx, p.client, payload.ParentPageID, payload.AttackPattern)
	if err != nil {
		return fmt.Errorf("CreateAttackPatternPage failed: %v: %w", err, asynq.SkipRetry)
	}

	return nil
}
