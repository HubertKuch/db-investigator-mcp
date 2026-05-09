# 🔍 DB Investigator MCP Server

This is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that allows AI agents (like Claude Desktop) to safely explore, schema-map, and query your databases.

Built with Go for maximum performance and zero-runtime overhead.

## 🚀 Features

- **Multi-DB Support**: Works with **PostgreSQL**, **MySQL**, and **SQLite** out of the box.
- **Smart Schema Mapping**: Dedicated tool to extract DDL and table definitions so the AI understands your data model.
- **Safety First**: The execution tool is strictly `readonly`. It enforces `SELECT` statements only to prevent accidental `DROP TABLE` or data loss.
- **Local Caching**: Automatically caches schema definitions locally to reduce database load and speed up context retrieval.
- **Zero Dependencies**: Compiled Go binary—no need for Node.js or Python on your host.

### Supported databases 

| Database | Supported |
| :--- | :---: |
| SQLite | ✅ |
| PostgreSQL | ✅ |
| MySQL | ✅ |

## 🛠️ Tools Included

1. `list_dbs`: Lists all available databases in the cluster or directory.
2. `refresh_ddl_schema`: Fetches full DDL/table structures for a specific database.
3. `execute_readonly_stmt`: Executes a SQL `SELECT` statement and returns results as JSON.

## 📦 Installation

```bash
# Clone the repo
git https://github.com/HubertKuch/db-investigator-mcp
cd db-investigator-mcp

# Build the binary
go build -o db-mcp # or `just build` if you use just
```

## ⚙️ Configuration

The server is configured via environment variables.

| Variable | Description | Example |
| :--- | :--- | :--- |
| `DB_TYPE` | `postgres`, `mysql`, or `sqlite` | `postgres` |
| `DB_HOST` | Database host (PG/MySQL) | `localhost` |
| `DB_PORT` | Database port (PG/MySQL) | `5432` or `3306` |
| `DB_USER` | Database user | `admin` |
| `DB_PASSWORD`| Database password | `secret` |
| `DB_PATH` | Directory containing `.db` files (SQLite only) | `/home/user/sqlitedbs` |

## 🤖 Sample configuration

```json
{
  "mcpServers": {
    "my-databases": {
      "command": "/absolute/path/to/mcp-binary",
      "env": {
        "DB_TYPE": "postgres",
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "your_user",
        "DB_PASSWORD": "your_password"
      }
    }
  }
}
```

## 🏗️ Architecture

The project follows a clean, modular structure:
- `tools/`: MCP tool definitions and handlers.
- `utils/`: Core logic and database drivers.
- `utils/db_*.go`: Specific implementations for each database type.

Adding a new database (like Oracle or SQL Server) is as easy as implementing the `DBDriver` interface.

## 🧪 Running Tests

```bash
go test ./...
```

## 📄 License

MIT © 2026 Hubert Kuch
