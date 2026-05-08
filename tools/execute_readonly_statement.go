package tools

import (
	"context"
	"fmt"
	"log/slog"
	"praca-db-tools-mcp/utils"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func NewExecuteReadonlyStatementTool(driver utils.DBDriver) (mcp.Tool, server.ToolHandlerFunc) {
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

	return tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := utils.ExtractArguments(request)
		dbname, ok1 := args["databaseName"].(string)
		statement, ok2 := args["stmt"].(string)

		if !ok1 || !ok2 || dbname == "" || statement == "" {
			return nil, fmt.Errorf("`databaseName` and `stmt` arguments are required")
		}

		cleanStmt := strings.TrimSpace(statement)
		if !strings.HasPrefix(strings.ToUpper(cleanStmt), "SELECT") {
			return nil, fmt.Errorf("only SELECT statements are allowed for security reasons")
		}

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		jsonStatement := fmt.Sprintf("SELECT json_agg(t) FROM (%s) t;", strings.TrimSuffix(statement, ";"))

		result, err := driver.ExecuteStatement(dbname, jsonStatement)
		if err != nil {
			slog.Error("Failed to execute readonly statement", "db", dbname, "stmt", statement, "error", err)
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	}
}
