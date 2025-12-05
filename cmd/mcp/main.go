package main

import (
	"database/sql"
	"net/http"

	"github.com/connor-davis/zingfibre-core/cmd/mcp/trino"
	"github.com/gofiber/fiber/v2/log"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "github.com/trinodb/trino-go-client/trino"
)

func main() {
	trinoDsn := "http://user@localhost:8081"
	trinoDb, err := sql.Open("trino", trinoDsn)

	if err != nil {
		log.Fatalf("🔥 Failed to connect to Trino database: %s", err.Error())

		return
	}

	defer trinoDb.Close()

	server := mcp.NewServer(&mcp.Implementation{Name: "zing-mcp", Version: "v1.0.0"}, nil)

	// Register Trino tool
	trino := trino.New(trinoDb)

	mcp.AddTool(server, &mcp.Tool{Name: "list-catalogs", Description: "Get a list of catalogs using TrinoDB."}, trino.ListCatalogs)
	mcp.AddTool(server, &mcp.Tool{Name: "list-schemas", Description: "Get a list of schemas for a given catalog using TrinoDB."}, trino.ListSchemas)
	mcp.AddTool(server, &mcp.Tool{Name: "list-tables", Description: "Get a list of tables for a given catalog and schema using TrinoDB."}, trino.ListTables)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	log.Info("MCP server listening on http://localhost:6174")

	if err := http.ListenAndServe("0.0.0.0:6174", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
