package ai

import (
	"github.com/connor-davis/zingfibre-core/common"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AI interface {
	CreateDynamicQuery(queryId uuid.UUID, prompt string)
}

type ai struct {
	postgres *postgres.Queries
	openai   openai.Client
}

func New(postgres *postgres.Queries) AI {
	openai := openai.NewClient(option.WithAPIKey(common.EnvString("OPENAI_KEY", "")))

	return &ai{
		postgres: postgres,
		openai:   openai,
	}
}
