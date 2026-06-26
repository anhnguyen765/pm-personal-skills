package main

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
	"strings"
)

type ConfluenceSearchResult struct {
	Results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"results"`
}

type ConfluencePage struct {
	Title string `json:"title"`
	Body  struct {
		Storage struct {
			Value string `json:"value"`
		} `json:"storage"`
	} `json:"body"`
}

func loadEnv() {
	// Try to find .env file by walking up from current directory
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			file, err := os.Open(envPath)
			if err == nil {
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
						parts := strings.SplitN(line, "=", 2)
						key := strings.TrimSpace(parts[0])
						val := strings.TrimSpace(parts[1])
						// Remove quotes if present
						val = strings.Trim(val, "\"'")
						os.Setenv(key, val)
					}
				}
				file.Close()
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

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "Error: %s environment variable is not set.\n", name)
		os.Exit(1)
	}
	return val
}

func searchConfluence(query, urlStr, email, token string) ([]struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}, error) {
	cql := fmt.Sprintf("text ~ \"%s\"", query)
	params := url.Values{}
	params.Add("cql", cql)
	params.Add("limit", "5")

	searchURL := fmt.Sprintf("%s/rest/api/content/search?%s", urlStr, params.Encode())
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.SetBasicAuth(email, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("search failed with status %d", resp.StatusCode)
	}

	var result ConfluenceSearchResult
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	return result.Results, nil
}

func getPageContent(pageID, url, email, token string) (string, string, error) {
	pageURL := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage", url, pageID)
	req, _ := http.NewRequest("GET", pageURL, nil)
	req.SetBasicAuth(email, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("fetch page failed with status %d", resp.StatusCode)
	}

	var page ConfluencePage
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &page)

	// Simple HTML tag removal
	re := regexp.MustCompile("<.*?>")
	cleanText := re.ReplaceAllString(page.Body.Storage.Value, "")
	cleanText = html.UnescapeString(cleanText)

	return page.Title, cleanText, nil
}

func main() {
	loadEnv()
	if len(os.Args) < 2 {
		fmt.Println("Usage: fetch_confluence <search_query>")
		os.Exit(1)
	}

	query := os.Args[1]
	url := strings.TrimRight(getEnvVar("CONFLUENCE_URL"), "/")
	email := getEnvVar("CONFLUENCE_EMAIL")
	token := getEnvVar("CONFLUENCE_API_TOKEN")

	var results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	var err error

	isNumeric := true
	for _, char := range query {
		if char < '0' || char > '9' {
			isNumeric = false
			break
		}
	}

	if isNumeric {
		results = append(results, struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}{ID: query, Title: "Direct Fetch"})
	} else {
		results, err = searchConfluence(query, url, email, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching Confluence: %v\n", err)
			os.Exit(1)
		}

		if len(results) == 0 {
			fmt.Printf("No results found for '%s'.\n", query)
			return
		}
	}

	fmt.Printf("--- Top %d Results ---\n", len(results))
	for _, res := range results {
		title, content, err := getPageContent(res.ID, url, email, token)
		if err != nil {
			fmt.Printf("\nERROR fetching %s: %v\n", res.Title, err)
			continue
		}

		fmt.Printf("\nSOURCE: %s (ID: %s)\n", title, res.ID)
		fmt.Println(strings.Repeat("-", 20))
		fmt.Println(content)
		fmt.Println("\n" + strings.Repeat("=", 40))
	}
}
