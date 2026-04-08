package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"praca-db-tools-mcp/utils"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CreateListDatabasesTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list_dbs",
		mcp.WithDescription("Zwraca listę wszystkich dostępnych baz danych w klastrze."),
	)

	return tool, listDatabasesHandler()
}

func listDatabasesHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var saveToGlobalCache = func(content string) (string, error) {
		cacheDir, _ := utils.GetCacheDir()

		fullPath := filepath.Join(cacheDir, "available_dbs.txt")

		err := os.WriteFile(fullPath, []byte(content), 0644)

		if err != nil {
			return "", fmt.Errorf("błąd podczas zapisu pliku: %w", err)
		}

		return fullPath, nil
	}

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

		query := "SELECT json_agg(datname) FROM pg_database WHERE datistemplate = false AND datname != 'postgres';"

		result, err := utils.ExecuteStatement("postgres", query)

		if err != nil {
			return nil, fmt.Errorf("nie udało się pobrać listy baz: %w", err)
		}

		if result == "" || strings.TrimSpace(result) == "null" {
			return mcp.NewToolResultText("[]"), nil
		}

		_, err = saveToGlobalCache(result)

		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	}
}
