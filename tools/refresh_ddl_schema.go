package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"praca-db-tools-mcp/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CreateRefreshDDLSchemaTool(driver utils.DBDriver) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("refresh_ddl_schema",
		mcp.WithDescription("Pobiera strukturę DDL (definicje tabel) dla podanej bazy danych. Używaj tego, gdy użytkownik pyta o strukturę, klucze lub tabele."),
		mcp.WithString("databaseName",
			mcp.Required(),
			mcp.Description("Nazwa bazy danych do sprawdzenia (np. backoffice, finanse)"),
		),
	)

	return tool, createRefreshDDLSchemaHandler(driver)
}

func createRefreshDDLSchemaHandler(driver utils.DBDriver) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var saveToGlobalCache = func(dbname string, content string) (string, error) {
		cacheDir, _ := utils.GetCacheDir()

		fileName := fmt.Sprintf("schema_%s.sql", dbname)
		fullPath := filepath.Join(cacheDir, fileName)

		err := os.WriteFile(fullPath, []byte(content), 0644)

		if err != nil {
			return "", fmt.Errorf("błąd podczas zapisu pliku: %w", err)
		}

		return fullPath, nil
	}

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := utils.ExtractArguments(request)

		dbName := args["databaseName"].(string)

		if dbName == "" {
			return nil, fmt.Errorf("databaseName argument is required")
		}

		ddlResult, err := driver.ExtractDDL(dbName)

		if err != nil {
			println(err.Error())

			return nil, err
		}

		_, cacheErr := saveToGlobalCache(dbName, ddlResult)

		if cacheErr != nil {
			return nil, cacheErr
		}

		return mcp.NewToolResultText(ddlResult), nil
	}
}
