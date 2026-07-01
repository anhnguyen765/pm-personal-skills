package main

// PM Tools MCP Server — exposes Jira, Confluence, and Web tools via MCP
// (Model Context Protocol) over stdio. Compatible with Claude Desktop, Claude
// Code, and any MCP-compliant client.
//
// Tools exposed:
//   jira_get          - Get a Jira issue by key (full details + comments)
//   jira_search       - JQL search returning a Markdown table
//   confluence_search - Full-text search across Confluence spaces
//   confluence_fetch  - Fetch a Confluence page by ID (cleaned text)
//   web_fetch         - Fetch & clean arbitrary web URLs (text / html / json)
//
// Credentials are loaded from the nearest .env file (walks up from cwd), or
// from the environment directly:
//   JIRA_URL, JIRA_API_TOKEN
//   CONFLUENCE_URL, CONFLUENCE_EMAIL, CONFLUENCE_API_TOKEN
//   WEB_USER_AGENT (optional), WEB_TIMEOUT (optional)
//
// Build:
//   cd pm-skills/tooling/mcp && go build -o pm-tools-mcp server.go
//
// Claude Desktop config  (~/.config/claude/claude_desktop_config.json):
//   {
//     "mcpServers": {
//       "pm-tools": {
//         "command": "/absolute/path/to/pm-tools-mcp",
//         "env": {
//           "JIRA_URL": "https://jira.example.com",
//           "JIRA_API_TOKEN": "...",
//           "CONFLUENCE_URL": "https://confluence.example.com",
//           "CONFLUENCE_EMAIL": "you@example.com",
//           "CONFLUENCE_API_TOKEN": "..."
//         }
//       }
//     }
//   }

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ─── MCP / JSON-RPC types ────────────────────────────────────────────────────

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type initResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    capabilities `json:"capabilities"`
	ServerInfo      serverInfo   `json:"serverInfo"`
}

type capabilities struct {
	Tools map[string]interface{} `json:"tools"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type mcpTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type toolsListResult struct {
	Tools []mcpTool `json:"tools"`
}

type toolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolCallResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ─── Jira types ──────────────────────────────────────────────────────────────

type issueDetail struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Description string `json:"description"`
		IssueType   struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Priority struct {
			Name string `json:"name"`
		} `json:"priority"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Reporter struct {
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Labels  []string `json:"labels"`
		Created string   `json:"created"`
		Updated string   `json:"updated"`
		Comment struct {
			Comments []struct {
				Author struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Body    string `json:"body"`
				Created string `json:"created"`
			} `json:"comments"`
		} `json:"comment"`
	} `json:"fields"`
}

type searchResult struct {
	Total  int `json:"total"`
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Summary   string `json:"summary"`
			IssueType struct {
				Name string `json:"name"`
			} `json:"issuetype"`
			Status struct {
				Name string `json:"name"`
			} `json:"status"`
			Priority struct {
				Name string `json:"name"`
			} `json:"priority"`
			Assignee struct {
				DisplayName string `json:"displayName"`
			} `json:"assignee"`
			Updated string `json:"updated"`
		} `json:"fields"`
	} `json:"issues"`
}

// ─── Confluence types ────────────────────────────────────────────────────────

type confluenceSearchResults struct {
	Results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"`
		Space struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		} `json:"space"`
		Links struct {
			WebUI string `json:"webui"`
		} `json:"_links"`
	} `json:"results"`
	TotalSize int `json:"totalSize"`
}

type pageContent struct {
	Title string `json:"title"`
	Space struct {
		Key  string `json:"key"`
		Name string `json:"name"`
	} `json:"space"`
	Version struct {
		Number int    `json:"number"`
		When   string `json:"when"`
	} `json:"version"`
	Body struct {
		Storage struct {
			Value string `json:"value"`
		} `json:"storage"`
	} `json:"body"`
	Links struct {
		WebUI string `json:"webui"`
	} `json:"_links"`
}

// ─── Env loading ─────────────────────────────────────────────────────────────

func readEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			k := strings.TrimSpace(parts[0])
			v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
			if os.Getenv(k) == "" { // env vars already set take priority
				os.Setenv(k, v)
			}
		}
	}
}

func loadEnv() {
	// 1. Explicit path wins — set by Claude Desktop config via PM_TOOLS_ENV_FILE
	if explicit := os.Getenv("PM_TOOLS_ENV_FILE"); explicit != "" {
		readEnvFile(explicit)
		return
	}
	// 2. Walk up from cwd — works for Claude Code and direct CLI use
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			readEnvFile(envPath)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func env(key string) string { return os.Getenv(key) }

func requireEnv(keys ...string) error {
	for _, k := range keys {
		if os.Getenv(k) == "" {
			return fmt.Errorf("required env var %s is not set (add to .env or environment)", k)
		}
	}
	return nil
}

// ─── Jira functions ──────────────────────────────────────────────────────────

func jiraHTTP(urlStr, token string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

func jiraGet(issueKey string) (string, error) {
	if err := requireEnv("JIRA_URL", "JIRA_API_TOKEN"); err != nil {
		return "", err
	}
	baseURL := strings.TrimRight(env("JIRA_URL"), "/")
	token := env("JIRA_API_TOKEN")

	fields := "summary,description,issuetype,status,priority,assignee,reporter,labels,created,updated,comment"
	endpoint := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=%s", baseURL, url.PathEscape(issueKey), fields)

	body, status, err := jiraHTTP(endpoint, token)
	if err != nil {
		return "", err
	}
	if status != 200 {
		return "", fmt.Errorf("Jira returned HTTP %d: %s", status, string(body))
	}

	var issue issueDetail
	json.Unmarshal(body, &issue)
	f := issue.Fields

	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s — %s\n\n", issue.Key, f.Summary)
	fmt.Fprintf(&sb, "**Type:** %s | **Status:** %s | **Priority:** %s\n", f.IssueType.Name, f.Status.Name, f.Priority.Name)
	fmt.Fprintf(&sb, "**Assignee:** %s | **Reporter:** %s\n", f.Assignee.DisplayName, f.Reporter.DisplayName)
	if len(f.Labels) > 0 {
		fmt.Fprintf(&sb, "**Labels:** %s\n", strings.Join(f.Labels, ", "))
	}
	if len(f.Created) >= 10 && len(f.Updated) >= 10 {
		fmt.Fprintf(&sb, "**Created:** %s | **Updated:** %s\n\n", f.Created[:10], f.Updated[:10])
	}
	if f.Description != "" {
		fmt.Fprintf(&sb, "## Description\n\n%s\n\n", f.Description)
	}
	if len(f.Comment.Comments) > 0 {
		fmt.Fprintf(&sb, "## Comments (%d)\n\n", len(f.Comment.Comments))
		for _, c := range f.Comment.Comments {
			ts := c.Created
			if len(ts) >= 10 {
				ts = ts[:10]
			}
			fmt.Fprintf(&sb, "**%s** (%s):\n%s\n\n", c.Author.DisplayName, ts, c.Body)
		}
	}
	return sb.String(), nil
}

func jiraSearch(jql string, limit int) (string, error) {
	if err := requireEnv("JIRA_URL", "JIRA_API_TOKEN"); err != nil {
		return "", err
	}
	baseURL := strings.TrimRight(env("JIRA_URL"), "/")
	token := env("JIRA_API_TOKEN")

	params := url.Values{}
	params.Set("jql", jql)
	params.Set("maxResults", strconv.Itoa(limit))
	params.Set("fields", "summary,issuetype,status,priority,assignee,updated")
	endpoint := fmt.Sprintf("%s/rest/api/2/search?%s", baseURL, params.Encode())

	body, status, err := jiraHTTP(endpoint, token)
	if err != nil {
		return "", err
	}
	if status != 200 {
		return "", fmt.Errorf("Jira returned HTTP %d: %s", status, string(body))
	}

	var result searchResult
	json.Unmarshal(body, &result)

	var sb strings.Builder
	fmt.Fprintf(&sb, "**JQL:** `%s`\n", jql)
	fmt.Fprintf(&sb, "**Total:** %d (showing %d)\n\n", result.Total, len(result.Issues))

	if len(result.Issues) == 0 {
		sb.WriteString("No issues found.")
		return sb.String(), nil
	}

	sb.WriteString("| Key | Type | Status | Priority | Assignee | Updated | Summary |\n")
	sb.WriteString("|---|---|---|---|---|---|---|\n")
	for _, issue := range result.Issues {
		f := issue.Fields
		updated := f.Updated
		if len(updated) >= 10 {
			updated = updated[:10]
		}
		summary := f.Summary
		if len(summary) > 60 {
			summary = summary[:57] + "..."
		}
		fmt.Fprintf(&sb, "| %s | %s | %s | %s | %s | %s | %s |\n",
			issue.Key, f.IssueType.Name, f.Status.Name, f.Priority.Name,
			f.Assignee.DisplayName, updated, summary)
	}
	return sb.String(), nil
}

// ─── Confluence functions ────────────────────────────────────────────────────

func confluenceHTTP(urlStr, email, token string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, 0, err
	}
	req.SetBasicAuth(email, token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}

func stripHTMLTags(raw string) string {
	tagRe := regexp.MustCompile(`<[^>]+>`)
	clean := tagRe.ReplaceAllString(raw, " ")
	clean = html.UnescapeString(clean)
	spaceRe := regexp.MustCompile(`[ \t]+`)
	clean = spaceRe.ReplaceAllString(clean, " ")
	lineRe := regexp.MustCompile(`\n{3,}`)
	clean = lineRe.ReplaceAllString(clean, "\n\n")
	return strings.TrimSpace(clean)
}

func confluenceSearch(query, space string, limit int) (string, error) {
	if err := requireEnv("CONFLUENCE_URL", "CONFLUENCE_EMAIL", "CONFLUENCE_API_TOKEN"); err != nil {
		return "", err
	}
	baseURL := strings.TrimRight(env("CONFLUENCE_URL"), "/")
	email := env("CONFLUENCE_EMAIL")
	token := env("CONFLUENCE_API_TOKEN")

	cql := fmt.Sprintf(`text ~ "%s" AND type in (page,blogpost)`, query)
	if space != "" {
		cql += fmt.Sprintf(` AND space.key = "%s"`, space)
	}
	params := url.Values{}
	params.Set("cql", cql)
	params.Set("limit", strconv.Itoa(limit))
	params.Set("expand", "space")
	endpoint := fmt.Sprintf("%s/rest/api/content/search?%s", baseURL, params.Encode())

	body, status, err := confluenceHTTP(endpoint, email, token)
	if err != nil {
		return "", err
	}
	if status != 200 {
		return "", fmt.Errorf("Confluence returned HTTP %d: %s", status, string(body))
	}

	var results confluenceSearchResults
	json.Unmarshal(body, &results)

	var sb strings.Builder
	fmt.Fprintf(&sb, "**Query:** `%s`\n", query)
	fmt.Fprintf(&sb, "**Total:** %d (showing %d)\n\n", results.TotalSize, len(results.Results))

	if len(results.Results) == 0 {
		sb.WriteString("No results found.")
		return sb.String(), nil
	}
	for i, r := range results.Results {
		webURL := ""
		if r.Links.WebUI != "" {
			webURL = baseURL + r.Links.WebUI
		}
		fmt.Fprintf(&sb, "%d. **%s** (ID: `%s`)\n", i+1, r.Title, r.ID)
		fmt.Fprintf(&sb, "   Space: %s (%s)", r.Space.Name, r.Space.Key)
		if webURL != "" {
			fmt.Fprintf(&sb, " | %s", webURL)
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func confluenceFetch(pageID string) (string, error) {
	if err := requireEnv("CONFLUENCE_URL", "CONFLUENCE_EMAIL", "CONFLUENCE_API_TOKEN"); err != nil {
		return "", err
	}
	baseURL := strings.TrimRight(env("CONFLUENCE_URL"), "/")
	email := env("CONFLUENCE_EMAIL")
	token := env("CONFLUENCE_API_TOKEN")

	endpoint := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage,space,version", baseURL, url.PathEscape(pageID))
	body, status, err := confluenceHTTP(endpoint, email, token)
	if err != nil {
		return "", err
	}
	if status != 200 {
		return "", fmt.Errorf("Confluence returned HTTP %d: %s", status, string(body))
	}

	var page pageContent
	json.Unmarshal(body, &page)

	webURL := ""
	if page.Links.WebUI != "" {
		webURL = baseURL + page.Links.WebUI
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", page.Title)
	when := page.Version.When
	if len(when) >= 10 {
		when = when[:10]
	}
	fmt.Fprintf(&sb, "**Space:** %s (%s) | **Version:** %d | **Updated:** %s\n",
		page.Space.Name, page.Space.Key, page.Version.Number, when)
	if webURL != "" {
		fmt.Fprintf(&sb, "**URL:** %s\n", webURL)
	}
	sb.WriteString("\n" + strings.Repeat("-", 60) + "\n\n")
	sb.WriteString(stripHTMLTags(page.Body.Storage.Value))
	return sb.String(), nil
}

// ─── Web fetch functions ─────────────────────────────────────────────────────

func stripHTMLFull(raw string) string {
	scriptRe := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	raw = scriptRe.ReplaceAllString(raw, " ")
	styleRe := regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`)
	raw = styleRe.ReplaceAllString(raw, " ")
	tagRe := regexp.MustCompile(`<[^>]+>`)
	clean := tagRe.ReplaceAllString(raw, " ")
	clean = html.UnescapeString(clean)
	spaceRe := regexp.MustCompile(`[ \t]+`)
	clean = spaceRe.ReplaceAllString(clean, " ")
	lineRe := regexp.MustCompile(`\n{3,}`)
	clean = lineRe.ReplaceAllString(clean, "\n\n")
	return strings.TrimSpace(clean)
}

func extractTitle(rawHTML string) string {
	re := regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)
	m := re.FindStringSubmatch(rawHTML)
	if len(m) > 1 {
		return stripHTMLTags(m[1])
	}
	return ""
}

func extractMeta(rawHTML string) map[string]string {
	tags := make(map[string]string)
	re := regexp.MustCompile(`(?i)<meta\s+(?:name|property)="([^"]+)"\s+content="([^"]+)"`)
	for _, m := range re.FindAllStringSubmatch(rawHTML, -1) {
		if len(m) > 2 {
			tags[m[1]] = m[2]
		}
	}
	return tags
}

func webFetch(urlStr, format string, timeout int, userAgent string) (string, error) {
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, urlStr)
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	contentType := resp.Header.Get("Content-Type")
	bodyStr := string(bodyBytes)

	switch format {
	case "html":
		return bodyStr, nil
	case "json":
		var data interface{}
		if err := json.Unmarshal(bodyBytes, &data); err != nil {
			return "", fmt.Errorf("response is not valid JSON: %v", err)
		}
		out, _ := json.MarshalIndent(data, "", "  ")
		return string(out), nil
	default: // text
		title := extractTitle(bodyStr)
		clean := stripHTMLFull(bodyStr)
		metaTags := extractMeta(bodyStr)
		var sb strings.Builder
		fmt.Fprintf(&sb, "# %s\n\n", title)
		fmt.Fprintf(&sb, "**URL:** %s | **Content-Type:** %s\n\n", urlStr, contentType)
		if desc, ok := metaTags["description"]; ok {
			fmt.Fprintf(&sb, "**Description:** %s\n\n", desc)
		}
		sb.WriteString(strings.Repeat("-", 60) + "\n\n")
		sb.WriteString(clean)
		return sb.String(), nil
	}
}

// ─── Tool definitions ────────────────────────────────────────────────────────

func toolList() []mcpTool {
	return []mcpTool{
		{
			Name:        "jira_get",
			Description: "Get full details of a Jira issue including description and comments.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"issue_key": map[string]interface{}{
						"type":        "string",
						"description": "Jira issue key, e.g. PCDPC-123",
					},
				},
				"required": []string{"issue_key"},
			},
		},
		{
			Name:        "jira_search",
			Description: "Search Jira issues using JQL. Returns a Markdown table.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"jql": map[string]interface{}{
						"type":        "string",
						"description": "JQL query string, e.g. \"project = PCDPC AND status = 'In Progress'\"",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (default: 20)",
					},
				},
				"required": []string{"jql"},
			},
		},
		{
			Name:        "confluence_search",
			Description: "Full-text search across Confluence pages and blog posts.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query text",
					},
					"space": map[string]interface{}{
						"type":        "string",
						"description": "Optional Confluence space key to limit the search (e.g. \"PAYMENTS\")",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results (default: 10)",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "confluence_fetch",
			Description: "Fetch the full text content of a Confluence page by its numeric page ID.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"page_id": map[string]interface{}{
						"type":        "string",
						"description": "Numeric Confluence page ID (from confluence_search results)",
					},
				},
				"required": []string{"page_id"},
			},
		},
		{
			Name:        "web_fetch",
			Description: "Fetch content from any public URL. Supports clean text extraction, raw HTML, or JSON parsing.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url": map[string]interface{}{
						"type":        "string",
						"description": "Full URL to fetch (https://...)",
					},
					"format": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"text", "html", "json"},
						"description": "Output format: text (default, cleaned), html (raw), or json (parsed)",
					},
					"timeout": map[string]interface{}{
						"type":        "integer",
						"description": "Request timeout in seconds (default: 30)",
					},
				},
				"required": []string{"url"},
			},
		},
	}
}

// ─── MCP request handling ────────────────────────────────────────────────────

var stdout = bufio.NewWriter(os.Stdout)

func send(v interface{}) {
	b, _ := json.Marshal(v)
	stdout.Write(b)
	stdout.WriteByte('\n')
	stdout.Flush()
}

func respond(id interface{}, result interface{}) {
	send(rpcResponse{JSONRPC: "2.0", ID: id, Result: result})
}

func respondErr(id interface{}, code int, msg string) {
	send(rpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}})
}

func toolOK(text string) toolCallResult {
	return toolCallResult{Content: []contentItem{{Type: "text", Text: text}}}
}

func toolErr(err error) toolCallResult {
	return toolCallResult{Content: []contentItem{{Type: "text", Text: err.Error()}}, IsError: true}
}

func handleToolCall(req rpcRequest) {
	var params toolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		respondErr(req.ID, -32602, "invalid params")
		return
	}

	var args map[string]interface{}
	json.Unmarshal(params.Arguments, &args)

	str := func(k string) string {
		if v, ok := args[k]; ok {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
	intArg := func(k string, def int) int {
		if v, ok := args[k]; ok {
			switch n := v.(type) {
			case float64:
				return int(n)
			case string:
				i, _ := strconv.Atoi(n)
				return i
			}
		}
		return def
	}

	switch params.Name {
	case "jira_get":
		key := str("issue_key")
		if key == "" {
			respond(req.ID, toolErr(fmt.Errorf("issue_key is required")))
			return
		}
		out, err := jiraGet(key)
		if err != nil {
			respond(req.ID, toolErr(err))
			return
		}
		respond(req.ID, toolOK(out))

	case "jira_search":
		jql := str("jql")
		if jql == "" {
			respond(req.ID, toolErr(fmt.Errorf("jql is required")))
			return
		}
		out, err := jiraSearch(jql, intArg("limit", 20))
		if err != nil {
			respond(req.ID, toolErr(err))
			return
		}
		respond(req.ID, toolOK(out))

	case "confluence_search":
		query := str("query")
		if query == "" {
			respond(req.ID, toolErr(fmt.Errorf("query is required")))
			return
		}
		out, err := confluenceSearch(query, str("space"), intArg("limit", 10))
		if err != nil {
			respond(req.ID, toolErr(err))
			return
		}
		respond(req.ID, toolOK(out))

	case "confluence_fetch":
		pageID := str("page_id")
		if pageID == "" {
			respond(req.ID, toolErr(fmt.Errorf("page_id is required")))
			return
		}
		out, err := confluenceFetch(pageID)
		if err != nil {
			respond(req.ID, toolErr(err))
			return
		}
		respond(req.ID, toolOK(out))

	case "web_fetch":
		rawURL := str("url")
		if rawURL == "" {
			respond(req.ID, toolErr(fmt.Errorf("url is required")))
			return
		}
		format := str("format")
		if format == "" {
			format = "text"
		}
		timeout := intArg("timeout", 30)
		out, err := webFetch(rawURL, format, timeout, env("WEB_USER_AGENT"))
		if err != nil {
			respond(req.ID, toolErr(err))
			return
		}
		respond(req.ID, toolOK(out))

	default:
		respond(req.ID, toolErr(fmt.Errorf("unknown tool: %s", params.Name)))
	}
}

func handleRequest(req rpcRequest) {
	switch req.Method {
	case "initialize":
		respond(req.ID, initResult{
			ProtocolVersion: "2024-11-05",
			Capabilities:    capabilities{Tools: map[string]interface{}{}},
			ServerInfo:      serverInfo{Name: "pm-tools", Version: "1.0.0"},
		})

	case "notifications/initialized":
		// no response for notifications

	case "ping":
		respond(req.ID, map[string]interface{}{})

	case "tools/list":
		respond(req.ID, toolsListResult{Tools: toolList()})

	case "tools/call":
		handleToolCall(req)

	default:
		if req.ID != nil {
			respondErr(req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
		}
	}
}

// ─── Entry point ─────────────────────────────────────────────────────────────

func main() {
	loadEnv()

	scanner := bufio.NewScanner(os.Stdin)
	// 4 MB buffer — handles large Confluence pages
	scanner.Buffer(make([]byte, 4*1024*1024), 4*1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			// malformed JSON — log to stderr only (stdout is reserved for protocol)
			fmt.Fprintf(os.Stderr, "mcp parse error: %v\n", err)
			continue
		}
		handleRequest(req)
	}
}
