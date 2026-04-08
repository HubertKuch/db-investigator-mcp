publish:
    go build -p 8 -x main.go
    mkdir -p ~/mcp-servers
    mv ./main ~/mcp-servers/db-investigator