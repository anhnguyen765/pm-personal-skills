package main

// Unified Jira CLI tool — replaces MCP jira_get_issue and jira_search_issues.
//
// Usage:
//   jira get <issue-key>
//   jira search "<jql>" [--limit N]
//
// Credentials loaded from nearest .env:
//   JIRA_URL, JIRA_API_TOKEN

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ---- structs ----------------------------------------------------------------

type IssueDetail struct {
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

type SearchResult struct {
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

// ---- helpers ----------------------------------------------------------------

func loadEnv() {
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			f, err := os.Open(envPath)
			if err == nil {
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
						parts := strings.SplitN(line, "=", 2)
						k := strings.TrimSpace(parts[0])
						v := strings.Trim(strings.TrimSpace(parts[1]), "\"'")
						os.Setenv(k, v)
					}
				}
				f.Close()
				return
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "error: %s not set in .env\n", key)
		os.Exit(1)
	}
	return v
}

func jiraRequest(method, urlStr, token string) ([]byte, int, error) {
	req, err := http.NewRequest(method, urlStr, nil)
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

// ---- commands ---------------------------------------------------------------

func cmdGet(baseURL, token, issueKey string) {
	fields := "summary,description,issuetype,status,priority,assignee,reporter,labels,created,updated,comment"
	endpoint := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=%s", baseURL, url.PathEscape(issueKey), fields)

	body, status, err := jiraRequest("GET", endpoint, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if status != 200 {
		fmt.Fprintf(os.Stderr, "error: Jira returned HTTP %d\n%s\n", status, string(body))
		os.Exit(1)
	}

	var issue IssueDetail
	json.Unmarshal(body, &issue)
	f := issue.Fields

	fmt.Printf("# %s — %s\n\n", issue.Key, f.Summary)
	fmt.Printf("**Type:** %s | **Status:** %s | **Priority:** %s\n", f.IssueType.Name, f.Status.Name, f.Priority.Name)
	fmt.Printf("**Assignee:** %s | **Reporter:** %s\n", f.Assignee.DisplayName, f.Reporter.DisplayName)
	if len(f.Labels) > 0 {
		fmt.Printf("**Labels:** %s\n", strings.Join(f.Labels, ", "))
	}
	fmt.Printf("**Created:** %s | **Updated:** %s\n\n", f.Created[:10], f.Updated[:10])

	if f.Description != "" {
		fmt.Printf("## Description\n\n%s\n\n", f.Description)
	}

	if len(f.Comment.Comments) > 0 {
		fmt.Printf("## Comments (%d)\n\n", len(f.Comment.Comments))
		for _, c := range f.Comment.Comments {
			ts := c.Created
			if len(ts) >= 10 {
				ts = ts[:10]
			}
			fmt.Printf("**%s** (%s):\n%s\n\n", c.Author.DisplayName, ts, c.Body)
		}
	}
}

func cmdSearch(baseURL, token, jql string, limit int) {
	params := url.Values{}
	params.Set("jql", jql)
	params.Set("maxResults", fmt.Sprintf("%d", limit))
	params.Set("fields", "summary,issuetype,status,priority,assignee,updated")
	endpoint := fmt.Sprintf("%s/rest/api/2/search?%s", baseURL, params.Encode())

	body, status, err := jiraRequest("GET", endpoint, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if status != 200 {
		fmt.Fprintf(os.Stderr, "error: Jira returned HTTP %d\n%s\n", status, string(body))
		os.Exit(1)
	}

	var result SearchResult
	json.Unmarshal(body, &result)

	fmt.Printf("**JQL:** `%s`\n", jql)
	fmt.Printf("**Total:** %d (showing %d)\n\n", result.Total, len(result.Issues))

	if len(result.Issues) == 0 {
		fmt.Println("No issues found.")
		return
	}

	fmt.Printf("| Key | Type | Status | Priority | Assignee | Updated | Summary |\n")
	fmt.Printf("|---|---|---|---|---|---|---|\n")
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
		fmt.Printf("| %s | %s | %s | %s | %s | %s | %s |\n",
			issue.Key, f.IssueType.Name, f.Status.Name, f.Priority.Name,
			f.Assignee.DisplayName, updated, summary)
	}
}

// ---- main -------------------------------------------------------------------

func main() {
	loadEnv()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n  jira get <issue-key>\n  jira search \"<jql>\" [--limit N]\n")
		os.Exit(1)
	}

	baseURL := strings.TrimRight(mustEnv("JIRA_URL"), "/")
	token := mustEnv("JIRA_API_TOKEN")
	command := os.Args[1]

	switch command {
	case "get":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: jira get <issue-key>\n")
			os.Exit(1)
		}
		cmdGet(baseURL, token, os.Args[2])

	case "search":
		fs := flag.NewFlagSet("search", flag.ExitOnError)
		limit := fs.Int("limit", 20, "max results")
		fs.Parse(os.Args[2:])
		args := fs.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Usage: jira search \"<jql>\" [--limit N]\n")
			os.Exit(1)
		}
		cmdSearch(baseURL, token, args[0], *limit)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nCommands: get, search\n", command)
		os.Exit(1)
	}
}
