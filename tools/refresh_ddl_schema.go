package tools

import (
	"context"
	"fmt"
	"log/slog"
	"praca-db-tools-mcp/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewRefreshDDLSchemaTool(driver utils.DBDriver) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("refresh_ddl_schema",
		mcp.WithDescription("Pobiera strukturę DDL (definicje tabel) dla podanej bazy danych. Używaj tego, gdy użytkownik pyta o strukturę, klucze lub tabele."),
		mcp.WithString("databaseName",
			mcp.Required(),
			mcp.Description("Nazwa bazy danych do sprawdzenia (np. backoffice, finanse)"),
		),
	)

	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := utils.ExtractArguments(request)
		dbName, ok := args["databaseName"].(string)
		if !ok || dbName == "" {
			return nil, fmt.Errorf("databaseName argument is required")
		}

		cacheFile := fmt.Sprintf("schema_%s.sql", dbName)

		if cached, err := readFromCache(cacheFile); err == nil {
			slog.Info("Using cached DDL schema", "db", dbName)
			return mcp.NewToolResultText(cached), nil
		}

		ddlResult, err := driver.ExtractDDL(dbName)
		if err != nil {
			slog.Error("Failed to extract DDL", "db", dbName, "error", err)
			return nil, err
		}

		if _, err := saveToCache(cacheFile, ddlResult); err != nil {
			slog.Warn("Failed to save DDL to cache", "db", dbName, "error", err)
		}

		return mcp.NewToolResultText(ddlResult), nil
	}
}
