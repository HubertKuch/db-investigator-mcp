package tools

import (
	"context"
	"fmt"
	"praca-db-tools-mcp/utils"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func CreateExecuteReadonlyStatementTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("execute_readonly_stmt",
		mcp.WithDescription("Wykonuje zapytanie akceptujac jedynie `SELECT`. Pelne `readonly`"),
		mcp.WithString("databaseName",
			mcp.Required(),
			mcp.Description("Nazwa bazy danych (np. backoffice, finanse)"),
		),
		mcp.WithString("stmt",
			mcp.Required(),
			mcp.Description("Zapytanie do wykonania np `SELECT * FROM users WHERE TRUE`"),
		),
	)

	return tool, executeReadonlyStatementToolHandler()
}

func executeReadonlyStatementToolHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := utils.ExtractArguments(request)

		dbname := args["databaseName"].(string)
		statement := args["stmt"].(string)

		if dbname == "" || statement == "" {
			return nil, fmt.Errorf("`databaseName` and `statement` arguments are required")
		}

		cleanStmt := strings.TrimSpace(statement)
		if !strings.HasPrefix(strings.ToUpper(cleanStmt), "SELECT") {
			return nil, fmt.Errorf("only SELECT statements are allowed for security reasons")
		}

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		jsonStatement := fmt.Sprintf("SELECT json_agg(t) FROM (%s) t;", strings.TrimSuffix(statement, ";"))

		result, stmtErr := utils.ExecuteStatement(dbname, jsonStatement)

		if stmtErr != nil {
			return nil, stmtErr
		}

		return mcp.NewToolResultText(result), nil
	}
}
