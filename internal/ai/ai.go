package ai

import (
	"github.com/connor-davis/zingfibre-core/common"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type AI interface {
	GenerateDynamicQuery(queryId uuid.UUID, prompt string, ctx *fiber.Ctx) error
}

type ai struct {
	postgres *postgres.Queries
	openai   openai.Client
}

func New(postgres *postgres.Queries) AI {
	openai := openai.NewClient(option.WithAPIKey(common.EnvString("OPENAI_API_KEY", "")))

	return &ai{
		postgres: postgres,
		openai:   openai,
	}
}
