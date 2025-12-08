package dynamicQueries

import (
	"database/sql"

	"github.com/connor-davis/zingfibre-core/cmd/api/http/middleware"
	"github.com/connor-davis/zingfibre-core/common"
	"github.com/connor-davis/zingfibre-core/internal/models/system"
	"github.com/connor-davis/zingfibre-core/internal/mysql/radius"
	"github.com/connor-davis/zingfibre-core/internal/mysql/zing"
	"github.com/connor-davis/zingfibre-core/internal/postgres"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type DynamicQueriesRouter struct {
	Postgres   *postgres.Queries
	Zing       *zing.Queries
	Radius     *radius.Queries
	Middleware *middleware.Middleware
	Sessions   *session.Store
	OpenAI     openai.Client
	Trino      *sql.DB
}

func NewDynamicQueriesRouter(
	postgres *postgres.Queries,
	zing *zing.Queries,
	radius *radius.Queries,
	middleware *middleware.Middleware,
	sessions *session.Store,
	trino *sql.DB,
) *DynamicQueriesRouter {
	openai := openai.NewClient(option.WithAPIKey(common.EnvString("OPENAI_API_KEY", "")))

	return &DynamicQueriesRouter{
		Postgres:   postgres,
		Zing:       zing,
		Radius:     radius,
		Middleware: middleware,
		Sessions:   sessions,
		OpenAI:     openai,
		Trino:      trino,
	}
}

func (r *DynamicQueriesRouter) RegisterRoutes() []system.Route {
	return []system.Route{
		r.GetDynamicQueriesRoute(),
		r.GetDynamicQueryResultsRoute(),
		r.GetDynamicQueryExportRoute(),
		r.GenerateDynamicQueryRoute(),
		r.GetDynamicQueryRoute(),
		r.CreateDynamicQueryRoute(),
		r.UpdateDynamicQueryRoute(),
		r.DeleteDynamicQueryRoute(),
	}
}
