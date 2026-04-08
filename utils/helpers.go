package utils

import "github.com/mark3labs/mcp-go/mcp"

func ExtractArguments(request mcp.CallToolRequest) map[string]interface{} {
	args, ok := request.Params.Arguments.(map[string]interface{})

	if !ok {
		return nil
	}

	return args
}
