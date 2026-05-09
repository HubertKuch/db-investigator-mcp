package tools

import (
	"context"
	"log/slog"
	"praca-db-tools-mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewListDatabasesTool(driver utils.DBDriver) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list_dbs",
		mcp.WithDescription("Zwraca listę wszystkich dostępnych baz danych w klastrze."),
	)

	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cacheFile := "available_dbs.txt"

		// Try to read from cache first
		if cached, err := readFromCache(cacheFile); err == nil {
			slog.Info("Using cached database list")
			return mcp.NewToolResultText(cached), nil
		}

		result, err := driver.ListDatabases()
		if err != nil {
			slog.Error("Failed to list databases", "error", err)
			return nil, err
		}

		if result == "" || strings.TrimSpace(result) == "null" {
			return mcp.NewToolResultText("[]"), nil
		}

		if _, err := saveToCache(cacheFile, result); err != nil {
			slog.Warn("Failed to save database list to cache", "error", err)
		}

		return mcp.NewToolResultText(result), nil
	}
}
