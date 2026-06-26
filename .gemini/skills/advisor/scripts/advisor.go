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

type SearchResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
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

func searchConfluence(query, urlStr, email, token string) ([]SearchResult, error) {
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

	var data SearchResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	return data.Results, nil
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

	var page ConfluencePage
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &page)

	// Clean HTML
	re := regexp.MustCompile("<.*?>")
	cleanText := re.ReplaceAllString(page.Body.Storage.Value, " ")
	cleanText = html.UnescapeString(cleanText)
	
	// Remove excessive whitespace
	spaceRe := regexp.MustCompile(`\s+`)
	cleanText = spaceRe.ReplaceAllString(cleanText, " ")

	return page.Title, cleanText, nil
}

func main() {
	loadEnv()
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  advisor search <query>")
		fmt.Println("  advisor fetch <page_id>")
		os.Exit(1)
	}

	command := os.Args[1]
	arg := os.Args[2]

	url := strings.TrimRight(getEnvVar("CONFLUENCE_URL"), "/")
	email := getEnvVar("CONFLUENCE_EMAIL")
	token := getEnvVar("CONFLUENCE_API_TOKEN")

	switch command {
	case "search":
		results, err := searchConfluence(arg, url, email, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for _, res := range results {
			fmt.Printf("ID: %s | Title: %s\n", res.ID, res.Title)
		}
	case "fetch":
		title, content, err := getPageContent(arg, url, email, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("TITLE: %s\n", title)
		fmt.Println(strings.Repeat("-", 20))
		fmt.Println(content)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
