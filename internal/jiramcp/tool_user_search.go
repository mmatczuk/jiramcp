package jiramcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type UserSearchArgs struct {
	Query string `json:"query" jsonschema:"Search query: display name, email address, or partial match. Examples: \"John Parker\", \"john.parker@\", \"parker\"."`
}

var userSearchTool = &mcp.Tool{
	Name: "jira_user_search",
	Description: `Search for Jira users by name or email. Returns display name and account ID.

Use this to find the account ID needed for jira_write assignee field. Accepts partial matches — search by first name, last name, or email prefix.`,
}

func (h *handlers) handleUserSearch(ctx context.Context, _ *mcp.CallToolRequest, args UserSearchArgs) (*mcp.CallToolResult, any, error) {
	if args.Query == "" {
		return textResult("query is required. Provide a display name, email, or partial match.", true), nil, nil
	}

	users, err := h.client.SearchUsers(ctx, args.Query)
	if err != nil {
		return textResult(fmt.Sprintf("User search failed: %v", err), true), nil, nil
	}

	if len(users) == 0 {
		return textResult(fmt.Sprintf("No users found for %q. Try a different spelling or use an email address.", args.Query), false), nil, nil
	}

	var results []map[string]any
	for _, u := range users {
		entry := map[string]any{
			"accountId":   u.AccountID,
			"displayName": u.DisplayName,
		}
		if u.EmailAddress != "" {
			entry["emailAddress"] = u.EmailAddress
		}
		results = append(results, entry)
	}

	data, _ := json.Marshal(results)
	out := fmt.Sprintf("Found %d user(s) matching %q. Use the accountId with jira_write assignee field.\n\n%s", len(results), args.Query, string(data))

	return textResult(out, false), nil, nil
}
