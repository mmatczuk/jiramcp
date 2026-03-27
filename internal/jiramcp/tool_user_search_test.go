package jiramcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/mmatczuk/jira-mcp/internal/jira"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func callUserSearch(t *testing.T, h *handlers, args UserSearchArgs) (string, bool) {
	t.Helper()
	result, _, err := h.handleUserSearch(context.Background(), nil, args)
	require.NoError(t, err)
	text := result.Content[0].(*mcp.TextContent).Text
	return text, result.IsError
}

func TestUserSearch_Success(t *testing.T) {
	mc := &mockClient{
		SearchUsersFn: func(_ context.Context, query string) ([]jira.User, error) {
			assert.Equal(t, "chaput", query)
			return []jira.User{
				{DisplayName: "Jon Chaput", AccountID: "abc123", EmailAddress: "jon.chaput@example.com"},
				{DisplayName: "Jane Chaput", AccountID: "def456"},
			}, nil
		},
	}
	h := &handlers{client: mc}
	text, isErr := callUserSearch(t, h, UserSearchArgs{Query: "chaput"})
	assert.False(t, isErr)
	assert.Contains(t, text, "Found 2 user(s)")
	assert.Contains(t, text, "abc123")
	assert.Contains(t, text, "Jon Chaput")
	assert.Contains(t, text, "jon.chaput@example.com")
	assert.Contains(t, text, "def456")
	assert.Contains(t, text, "accountId")
}

func TestUserSearch_NoResults(t *testing.T) {
	mc := &mockClient{
		SearchUsersFn: func(context.Context, string) ([]jira.User, error) {
			return nil, nil
		},
	}
	h := &handlers{client: mc}
	text, isErr := callUserSearch(t, h, UserSearchArgs{Query: "nonexistent"})
	assert.False(t, isErr)
	assert.Contains(t, text, "No users found")
}

func TestUserSearch_EmptyQuery(t *testing.T) {
	h := &handlers{client: &mockClient{}}
	text, isErr := callUserSearch(t, h, UserSearchArgs{Query: ""})
	assert.True(t, isErr)
	assert.Contains(t, text, "query is required")
}

func TestUserSearch_ClientError(t *testing.T) {
	mc := &mockClient{
		SearchUsersFn: func(context.Context, string) ([]jira.User, error) {
			return nil, fmt.Errorf("connection refused")
		},
	}
	h := &handlers{client: mc}
	text, isErr := callUserSearch(t, h, UserSearchArgs{Query: "test"})
	assert.True(t, isErr)
	assert.Contains(t, text, "User search failed")
	assert.Contains(t, text, "connection refused")
}

func TestUserSearch_OmitsEmptyEmail(t *testing.T) {
	mc := &mockClient{
		SearchUsersFn: func(context.Context, string) ([]jira.User, error) {
			return []jira.User{
				{DisplayName: "Bot User", AccountID: "bot1"},
			}, nil
		},
	}
	h := &handlers{client: mc}
	text, isErr := callUserSearch(t, h, UserSearchArgs{Query: "bot"})
	assert.False(t, isErr)
	assert.Contains(t, text, "Found 1 user(s)")
	assert.NotContains(t, text, "emailAddress")
}
