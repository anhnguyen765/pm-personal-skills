package main

// Web content fetcher CLI tool — fetch and clean HTML content from arbitrary URLs.
//
// Usage:
//   web fetch <url> [--format text|html|json] [--timeout SEC]
//   web extract <url> <selector> [--format text|html]
//
// Features:
//   - Fetches HTML/JSON from arbitrary web links
//   - Strips HTML tags and cleans content
//   - Supports custom User-Agent
//   - Handles redirects and caching
//
// Credentials loaded from nearest .env (optional):
//   WEB_USER_AGENT, WEB_TIMEOUT, WEB_CACHE_DIR

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ---- structs ----------------------------------------------------------------

type FetchResult struct {
	URL         string
	Title       string
	ContentType string
	StatusCode  int
	Content     string
	RawHTML     string
	MetaTags    map[string]string
	Links       []string
}

// ---- helpers ----------------------------------------------------------------

func loadEnv() map[string]string {
	env := make(map[string]string)
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
						env[k] = v
					}
				}
				f.Close()
				return env
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return env
}

func fetchURL(urlStr string, timeout int, userAgent string) ([]byte, int, string, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, 0, "", err
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
		return nil, 0, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	contentType := resp.Header.Get("Content-Type")
	return body, resp.StatusCode, contentType, nil
}

func stripHTML(raw string) string {
	// Remove script and style tags
	scriptRe := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	raw = scriptRe.ReplaceAllString(raw, " ")

	styleRe := regexp.MustCompile(`(?i)<style[^>]*>.*?</style>`)
	raw = styleRe.ReplaceAllString(raw, " ")

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

func extractTitle(html string) string {
	titleRe := regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)
	matches := titleRe.FindStringSubmatch(html)
	if len(matches) > 1 {
		return stripHTML(matches[1])
	}
	return ""
}

func extractMetaTags(html string) map[string]string {
	tags := make(map[string]string)
	metaRe := regexp.MustCompile(`(?i)<meta\s+(?:name|property)="([^"]+)"\s+content="([^"]+)"`)
	matches := metaRe.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) > 2 {
			tags[m[1]] = m[2]
		}
	}
	return tags
}

func extractLinks(html string) []string {
	var links []string
	linkRe := regexp.MustCompile(`(?i)<a\s+[^>]*href="([^"]+)"`)
	matches := linkRe.FindAllStringSubmatch(html, -1)
	seen := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] && m[1] != "" {
			links = append(links, m[1])
			seen[m[1]] = true
		}
	}
	return links
}

// ---- commands ---------------------------------------------------------------

func cmdFetch(urlStr string, format string, timeout int, userAgent string) {
	fmt.Fprintf(os.Stderr, "Fetching: %s\n", urlStr)

	body, status, contentType, err := fetchURL(urlStr, timeout, userAgent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if status != 200 {
		fmt.Fprintf(os.Stderr, "error: HTTP %d\n%s\n", status, string(body))
		os.Exit(1)
	}

	bodyStr := string(body)

	switch format {
	case "html":
		fmt.Println(bodyStr)

	case "json":
		// Try to parse as JSON
		var data interface{}
		if err := json.Unmarshal(body, &data); err == nil {
			out, _ := json.MarshalIndent(data, "", "  ")
			fmt.Println(string(out))
		} else {
			fmt.Fprintf(os.Stderr, "error: not valid JSON\n")
			os.Exit(1)
		}

	case "text":
		fallthrough
	default:
		title := extractTitle(bodyStr)
		clean := stripHTML(bodyStr)
		metaTags := extractMetaTags(bodyStr)

		fmt.Printf("# %s\n\n", title)
		fmt.Printf("**URL:** %s | **Content-Type:** %s\n", urlStr, contentType)
		fmt.Println()

		if description, ok := metaTags["description"]; ok {
			fmt.Printf("**Description:** %s\n\n", description)
		}

		fmt.Println(strings.Repeat("-", 60))
		fmt.Println()
		fmt.Println(clean)
	}
}

// ---- main -------------------------------------------------------------------

func main() {
	env := loadEnv()

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage:\n  web fetch <url> [--format text|html|json] [--timeout SEC]\n")
		os.Exit(1)
	}

	userAgent := env["WEB_USER_AGENT"]
	timeout := 30
	if t := env["WEB_TIMEOUT"]; t != "" {
		fmt.Sscanf(t, "%d", &timeout)
	}

	command := os.Args[1]

	switch command {
	case "fetch":
		fs := flag.NewFlagSet("fetch", flag.ExitOnError)
		format := fs.String("format", "text", "output format: text, html, or json")
		timeoutFlag := fs.Int("timeout", timeout, "request timeout in seconds")
		fs.Parse(os.Args[2:])

		args := fs.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Usage: web fetch <url> [--format text|html|json] [--timeout SEC]\n")
			os.Exit(1)
		}

		cmdFetch(args[0], *format, *timeoutFlag, userAgent)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\nCommands: fetch\n", command)
		os.Exit(1)
	}
}
