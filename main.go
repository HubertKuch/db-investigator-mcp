package main

import (
	"praca-db-tools-mcp/tools"
	"praca-db-tools-mcp/utils"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	_, err := utils.EnsureCacheDir()

	if err != nil {
		panic(err)
	}

	driver, err := utils.GetDriver()
	if err != nil {
		panic(err)
	}

	s := server.NewMCPServer("db-tools-server", "1.0.0")

	t1, h1 := tools.CreateRefreshDDLSchemaTool(driver)
	s.AddTool(t1, h1)

	t2, h2 := tools.CreateExecuteReadonlyStatementTool(driver)
	s.AddTool(t2, h2)

	t3, h3 := tools.CreateListDatabasesTool(driver)
	s.AddTool(t3, h3)

	if err := server.ServeStdio(s); err != nil {
		panic(err)
	}
}
