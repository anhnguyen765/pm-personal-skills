package main

// Unified Confluence CLI tool — replaces MCP confluence_search and
// confluence_get_page_content.
//
// Usage:
//   confluence search "<query>" [--space SPACE] [--limit N]
//   confluence fetch <page_id>
//
// Credentials loaded from nearest .env:
//   CONFLUENCE_URL, CONFLUENCE_EMAIL, CONFLUENCE_API_TOKEN

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ---- structs ----------------------------------------------------------------

type SearchResults struct {
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

type PageContent struct {
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

func confluenceRequest(urlStr, email, token string) ([]byte, int, error) {
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

func stripHTML(raw string) string {
	// Remove HTML tags
	tagRe := regexp.MustCompile(`<[^>]+>`)
	clean := tagRe.ReplaceAllString(raw, " ")
	// Decode HTML entities
	clean = html.UnescapeString(clean)
	// Collapse whitespace
	spaceRe := regexp.MustCompile(`[ \t]+`)
	clean = spaceRe.ReplaceAllString(clean, " ")
	lineRe := regexp.MustCompile(`\n{3,}`)
	clean = lineRe.ReplaceAllString(clean, "\n\n")
	return strings.TrimSpace(clean)
}

// ---- commands ---------------------------------------------------------------

func cmdSearch(baseURL, email, token, query, space string, limit int) {
	cql := fmt.Sprintf("text ~ \"%s\" AND type in (page,blogpost)", query)
	if space != "" {
		cql += fmt.Sprintf(` AND space.key = "%s"`, space)
	}

	params := url.Values{}
	params.Set("cql", cql)
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("expand", "space")
	endpoint := fmt.Sprintf("%s/rest/api/content/search?%s", baseURL, params.Encode())

	body, status, err := confluenceRequest(endpoint, email, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if status != 200 {
		fmt.Fprintf(os.Stderr, "error: Confluence returned HTTP %d\n%s\n", status, string(body))
		os.Exit(1)
	}

	var results SearchResults
	json.Unmarshal(body, &results)

	fmt.Printf("**Query:** `%s`\n", query)
	fmt.Printf("**Total:** %d (showing %d)\n\n", results.TotalSize, len(results.Results))

	if len(results.Results) == 0 {
		fmt.Println("No results found.")
		return
	}

	for i, r := range results.Results {
		webURL := ""
		if r.Links.WebUI != "" {
			webURL = strings.TrimRight(baseURL, "/") + r.Links.WebUI
		}
		fmt.Printf("%d. **%s** (ID: `%s`)\n", i+1, r.Title, r.ID)
		fmt.Printf("   Space: %s (%s)", r.Space.Name, r.Space.Key)
		if webURL != "" {
			fmt.Printf(" | %s", webURL)
		}
		fmt.Println()
	}
}

func cmdFetch(baseURL, email, token, pageID string) {
	endpoint := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage,space,version", baseURL, url.PathEscape(pageID))

	body, status, err := confluenceRequest(endpoint, email, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if status != 200 {
		fmt.Fprintf(os.Stderr, "error: Confluence returned HTTP %d\n%s\n", status, string(body))
		os.Exit(1)
	}

	var page PageContent
	json.Unmarshal(body, &page)

	webURL := ""
	if page.Links.WebUI != "" {
		webURL = strings.TrimRight(baseURL, "/") + page.Links.WebUI
	}

	fmt.Printf("# %s\n\n", page.Title)
	fmt.Printf("**Space:** %s (%s) | **Version:** %d | **Updated:** %s\n",
		page.Space.Name, page.Space.Key, page.Version.Number, page.Version.When[:10])
	if webURL != "" {
		fmt.Printf("**URL:** %s\n", webURL)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println()
	fmt.Println(stripHTML(page.Body.Storage.Value))
}

// ---- main -------------------------------------------------------------------

func main() {
	loadEnv()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n  confluence search \"<query>\" [--space KEY] [--limit N]\n  confluence fetch <page_id>\n")
		os.Exit(1)
	}

	baseURL := strings.TrimRight(mustEnv("CONFLUENCE_URL"), "/")
	email := mustEnv("CONFLUENCE_EMAIL")
	token := mustEnv("CONFLUENCE_API_TOKEN")
	command := os.Args[1]

	switch command {
	case "search":
		fs := flag.NewFlagSet("search", flag.ExitOnError)
		space := fs.String("space", "", "limit to a Confluence space key")
		limit := fs.Int("limit", 10, "max results")
		fs.Parse(os.Args[2:])
		args := fs.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Usage: confluence search \"<query>\" [--space KEY] [--limit N]\n")
			os.Exit(1)
		}
		cmdSearch(baseURL, email, token, args[0], *space, *limit)

	case "fetch":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: confluence fetch <page_id>\n")
			os.Exit(1)
		}
		cmdFetch(baseURL, email, token, os.Args[2])

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nCommands: search, fetch\n", command)
		os.Exit(1)
	}
}
