package tools

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

type MockDriver struct {
	ExecuteFunc    func(dbname string, statement string) (string, error)
	ExtractDDLFunc func(dbname string) (string, error)
	ListDBsFunc    func() (string, error)
}

func (m *MockDriver) ExecuteStatement(dbname string, statement string) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(dbname, statement)
	}
	return "", nil
}

func (m *MockDriver) ExtractDDL(dbname string) (string, error) {
	if m.ExtractDDLFunc != nil {
		return m.ExtractDDLFunc(dbname)
	}
	return "", nil
}

func (m *MockDriver) ListDatabases() (string, error) {
	if m.ListDBsFunc != nil {
		return m.ListDBsFunc()
	}
	return "", nil
}

func TestExecuteReadonlyStatementTool(t *testing.T) {
	mock := &MockDriver{
		ExecuteFunc: func(dbname string, statement string) (string, error) {
			if dbname == "testdb" && statement == "SELECT json_agg(t) FROM (SELECT * FROM users) t;" {
				return `[{"id":1, "name":"test"}]`, nil
			}
			return "", fmt.Errorf("unexpected call")
		},
	}

	_, handler := CreateExecuteReadonlyStatementTool(mock)

	req := mcp.CallToolRequest{}
	req.Params.Name = "execute_readonly_stmt"
	req.Params.Arguments = map[string]any{
		"databaseName": "testdb",
		"stmt":         "SELECT * FROM users",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("expected content in result")
	}

	text := result.Content[0].(mcp.TextContent).Text
	if text != `[{"id":1, "name":"test"}]` {
		t.Errorf("expected result text, got %s", text)
	}
}

func TestExecuteReadonlyStatementTool_Security(t *testing.T) {
	mock := &MockDriver{}
	_, handler := CreateExecuteReadonlyStatementTool(mock)

	req := mcp.CallToolRequest{}
	req.Params.Name = "execute_readonly_stmt"
	req.Params.Arguments = map[string]any{
		"databaseName": "testdb",
		"stmt":         "DROP TABLE users",
	}

	_, err := handler(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for non-SELECT statement, got nil")
	}
}

func TestListDatabasesTool(t *testing.T) {
	mock := &MockDriver{
		ListDBsFunc: func() (string, error) {
			return `["db1", "db2"]`, nil
		},
	}

	_, handler := CreateListDatabasesTool(mock)

	req := mcp.CallToolRequest{}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	if text != `["db1", "db2"]` {
		t.Errorf("expected result text, got %s", text)
	}
}

func TestRefreshDDLSchemaTool(t *testing.T) {
	mock := &MockDriver{
		ExtractDDLFunc: func(dbname string) (string, error) {
			if dbname == "testdb" {
				return "CREATE TABLE users (...);", nil
			}
			return "", fmt.Errorf("unexpected call")
		},
	}

	_, handler := CreateRefreshDDLSchemaTool(mock)

	req := mcp.CallToolRequest{}
	req.Params.Name = "refresh_ddl_schema"
	req.Params.Arguments = map[string]any{
		"databaseName": "testdb",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := result.Content[0].(mcp.TextContent).Text
	if text != "CREATE TABLE users (...);" {
		t.Errorf("expected DDL result, got %s", text)
	}
}
