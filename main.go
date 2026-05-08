package main

import (
	"log/slog"
	"os"

	"praca-db-tools-mcp/tools"
	"praca-db-tools-mcp/utils"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	if _, err := utils.GetCacheDir(); err != nil {
		slog.Error("Failed to initialize cache directory", "error", err)
		os.Exit(1)
	}

	driver, err := utils.GetDriver()
	if err != nil {
		slog.Error("Failed to get database driver", "error", err)
		os.Exit(1)
	}

	s := server.NewMCPServer("db-tools-server", "1.0.0")

	s.AddTool(tools.NewRefreshDDLSchemaTool(driver))
	s.AddTool(tools.NewExecuteReadonlyStatementTool(driver))
	s.AddTool(tools.NewListDatabasesTool(driver))

	slog.Info("Starting MCP server")
	if err := server.ServeStdio(s); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}
