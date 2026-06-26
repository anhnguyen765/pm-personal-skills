package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type ConfluenceSearchResult struct {
	Results []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"results"`
}

func loadEnv() {
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
	params.Add("limit", "15")

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

func main() {
	loadEnv()
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run search_titles.go <query>")
		os.Exit(1)
	}

	query := os.Args[1]
	url := strings.TrimRight(getEnvVar("CONFLUENCE_URL"), "/")
	email := getEnvVar("CONFLUENCE_EMAIL")
	token := getEnvVar("CONFLUENCE_API_TOKEN")

	results, err := searchConfluence(query, url, email, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Search Results for: %s\n", query)
	for i, res := range results {
		fmt.Printf("[%d] ID: %s | Title: %s\n", i+1, res.ID, res.Title)
	}
}
