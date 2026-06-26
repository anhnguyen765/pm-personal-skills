package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func main() {
	loadEnv()
	if len(os.Args) < 2 {
		fmt.Println("Usage: fetch_page <page_id>")
		os.Exit(1)
	}
	pageID := os.Args[1]
	urlStr := strings.TrimRight(getEnvVar("CONFLUENCE_URL"), "/")
	email := getEnvVar("CONFLUENCE_EMAIL")
	token := getEnvVar("CONFLUENCE_API_TOKEN")

	pageURL := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage", urlStr, pageID)
	req, _ := http.NewRequest("GET", pageURL, nil)
	req.SetBasicAuth(email, token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching page: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Error fetching page status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var page ConfluencePage
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &page)

	re := regexp.MustCompile("<.*?>")
	cleanText := re.ReplaceAllString(page.Body.Storage.Value, "\n")
	cleanText = html.UnescapeString(cleanText)

	fmt.Printf("TITLE: %s\n", page.Title)
	fmt.Println(cleanText)
}
