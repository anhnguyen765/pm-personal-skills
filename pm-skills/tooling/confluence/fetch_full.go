package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ConfluencePage struct {
	Title string `json:"title"`
	Body  struct {
		Storage struct {
			Value string `json:"value"`
		} `json:"storage"`
	} `json:"body"`
}

func getEnvVar(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Fprintf(os.Stderr, "Error: %s environment variable is not set.\n", name)
		os.Exit(1)
	}
	return val
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: fetch_full <page_id>")
		os.Exit(1)
	}

	pageID := os.Args[1]
	url := strings.TrimRight(getEnvVar("CONFLUENCE_URL"), "/")
	email := getEnvVar("CONFLUENCE_EMAIL")
	token := getEnvVar("CONFLUENCE_API_TOKEN")

	title, content, err := getPageContent(pageID, url, email, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("TITLE: %s\n", title)
	fmt.Println(strings.Repeat("-", 20))
	fmt.Println(content)
}
