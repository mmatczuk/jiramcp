# JIRA MCP

Give your AI assistant full JIRA access without burning half its context window on tool selection.

Most JIRA MCPs expose too many tools. Every time the model picks one, it spends tokens deciding between `jira_get_issue`, `jira_fetch_issue`, `jira_issue_get`...

`jira-mcp` has three:

| Tool | What it does |
|---|---|
| `jira_read` | Fetch issues by key, search by JQL, list projects/boards/sprints |
| `jira_write` | Create, update, delete, transition, comment — accepts Markdown, converts to Atlassian format; supports `dry_run` to preview changes before committing |
| `jira_schema` | Discover fields, transitions, and allowed values |

Less tool surface area means more of the context window goes to your actual work. The model makes fewer wrong choices, calls fewer redundant tools, and gets to the answer faster. The three tools compose naturally: schema to discover, read to find, write to change.

Your API token never leaves your machine. `jira-mcp` runs as a local process over stdio — no server to host, no proxy in the middle, no credentials sent anywhere except directly to Atlassian.

## Quick start

**Prerequisites:** A JIRA Cloud account with an API token.

### 1. Get an API token

1. Go to [API token management](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Click "Create API token"
3. Give it a label, click "Create"
4. Copy the token immediately — you can't see it again

### 2. Install and add to Claude Code

Pick one path:

**Homebrew**

```bash
brew tap mmatczuk/jira-mcp https://github.com/mmatczuk/jira-mcp
brew install jira-mcp
```

```bash
claude mcp add-json jira '{
  "command": "jira-mcp",
  "env": {
    "JIRA_URL": "https://yourcompany.atlassian.net",
    "JIRA_EMAIL": "you@company.com",
    "JIRA_API_TOKEN": "your-api-token"
  }
}'
```

**Docker**

No install needed. The `-e VAR` flags (without a value) forward each variable from `env` into the container:

```bash
claude mcp add-json jira '{
  "command": "docker",
  "args": [
    "run", "-i", "--rm",
    "-e", "JIRA_URL",
    "-e", "JIRA_EMAIL",
    "-e", "JIRA_API_TOKEN",
    "mmatczuk/jira-mcp"
  ],
  "env": {
    "JIRA_URL": "https://yourcompany.atlassian.net",
    "JIRA_EMAIL": "you@company.com",
    "JIRA_API_TOKEN": "your-api-token"
  }
}'
```

**Binary**

Download the binary for your platform from the [releases page](https://github.com/mmatczuk/jira-mcp/releases) and put it on your `PATH`, then:

```bash
claude mcp add-json jira '{
  "command": "jira-mcp",
  "env": {
    "JIRA_URL": "https://yourcompany.atlassian.net",
    "JIRA_EMAIL": "you@company.com",
    "JIRA_API_TOKEN": "your-api-token"
  }
}'
```

### 3. Verify it works

Ask Claude: *"List my JIRA projects"* — if you see your projects, you're good.

### Other MCP clients

Use the same binary and env vars. The server speaks standard MCP over stdio.

## Compared to [mcp-atlassian](https://github.com/sooperset/mcp-atlassian)

| | jira-mcp | mcp-atlassian |
|---|---|---|
| Tools | 3 | 72 |
| Runtime | Go binary | Python |
| Scope | Jira | Jira + Confluence |
| Auth | API token | API token, OAuth 2.0, PAT |
| Transport | stdio (local only) | stdio, SSE |

mcp-atlassian is the right choice if you need Confluence or OAuth. jira-mcp is for people who want Jira to just work, with minimal overhead on the model.

## Compared to [acli](https://developer.atlassian.com/cloud/acli/guides/introduction/)

acli is Atlassian's official CLI for humans running terminal commands. For AI-driven Jira automation, jira-mcp is safer: structured tool calls instead of shell execution, no risk of shell injection, and blast radius limited to Jira — the model can never touch user management or org settings.

| | jira-mcp | acli via bash |
|---|---|---|
| Integration | MCP tool (structured calls) | Shell command execution |
| Shell injection risk | None | Possible |
| Read-only mode | Instruct model to use `jira_read` only | Requires separate read-only API token |
| Blast radius | Jira only | Jira + Admin (user deletion, org settings) + Rovodev |
| Markdown input | Native (converts to Atlassian format) | None |
| Dry run | Built-in | None |
| Output format | Structured JSON | Text — requires parsing |
| Context cost | 3 tools | Full CLI surface (commands + flags) |
| Maintained by | Community | Atlassian |

acli is the right choice for human-driven terminal workflows or when you need Admin/Rovodev operations. jira-mcp is for AI-driven Jira automation.

## License

MIT
