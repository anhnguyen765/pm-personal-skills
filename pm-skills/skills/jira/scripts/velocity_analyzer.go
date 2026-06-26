package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type JiraIssue struct {
	Fields struct {
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Summary string `json:"summary"`
	} `json:"fields"`
}

func loadEnv() {
	current, err := os.Getwd()
	if err != nil {
		return
	}

	for {
		envFile := filepath.Join(current, ".env")
		if _, err := os.Stat(envFile); err == nil {
			file, err := os.Open(envFile)
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
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
}

func getLatestReport(reportDir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(reportDir, "daily_report_*.md"))
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no daily reports found in %s", reportDir)
	}

	sort.Slice(files, func(i, j int) bool {
		fi, _ := os.Stat(files[i])
		fj, _ := os.Stat(files[j])
		return fi.ModTime().After(fj.ModTime())
	})

	return files[0], nil
}

func analyzeVelocity() string {
	loadEnv()
	jiraURL := os.Getenv("JIRA_URL")
	jiraToken := os.Getenv("JIRA_API_TOKEN")
	reportDir := "/Users/lap14569/Library/CloudStorage/OneDrive-VNGGroupJSC/vng/Daily"

	latestFile, err := getLatestReport(reportDir)
	if err != nil {
		return err.Error()
	}

	content, err := os.ReadFile(latestFile)
	if err != nil {
		return fmt.Sprintf("Error reading file %s: %v", latestFile, err)
	}

	lines := strings.Split(string(content), "\n")
	var stuckTasks []string
	inStuckSection := false

	for _, line := range lines {
		if strings.Contains(line, "## 3. Stuck Tasks") {
			inStuckSection = true
			continue
		}
		if inStuckSection && strings.HasPrefix(line, "| PC-") {
			parts := strings.Split(line, "|")
			if len(parts) > 5 {
				stuckTasks = append(stuckTasks, strings.TrimSpace(parts[1]))
			}
		}
	}

	if len(stuckTasks) == 0 {
		return fmt.Sprintf("No stuck tasks found in the latest report: %s", filepath.Base(latestFile))
	}

	results := []string{fmt.Sprintf("# Velocity Analysis based on %s\n", filepath.Base(latestFile))}

	client := &http.Client{}
	headers := map[string]string{
		"Authorization": "Bearer " + jiraToken,
		"Accept":        "application/json",
	}

	for _, key := range stuckTasks {
		url := fmt.Sprintf("%s/rest/api/2/issue/%s?fields=status,assignee,summary", jiraURL, key)
		req, _ := http.NewRequest("GET", url, nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			results = append(results, fmt.Sprintf("### %s: Error connecting to Jira: %v", key, err))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var issue JiraIssue
			body, _ := io.ReadAll(resp.Body)
			json.Unmarshal(body, &issue)

			status := issue.Fields.Status.Name
			assignee := issue.Fields.Assignee.DisplayName
			if assignee == "" {
				assignee = "Unassigned"
			}
			summary := issue.Fields.Summary

			comment := fmt.Sprintf("**Suggested Comment for %s:**\n", key)
			switch status {
			case "In Dev":
				comment += fmt.Sprintf("> \"@%s, this task has been in 'In Dev' for several days. Are there any technical blockers or dependencies I can help resolve?\"\n", assignee)
			case "New":
				comment += fmt.Sprintf("> \"@%s, this task is still in 'New' status. Let's discuss if we should start it or if it needs to be reprioritized.\"\n", assignee)
			default:
				comment += fmt.Sprintf("> \"@%s, I see this is now in '%s'. Great progress! Let's ensure the documentation is updated.\"\n", assignee, status)
			}

			results = append(results, fmt.Sprintf("### %s: %s", key, summary))
			results = append(results, fmt.Sprintf("- **Current Status:** %s", status))
			results = append(results, fmt.Sprintf("- **Assignee:** %s", assignee))
			results = append(results, comment)
		} else {
			results = append(results, fmt.Sprintf("### %s: Error fetching data from Jira (Status %d)", key, resp.StatusCode))
		}
	}

	return strings.Join(results, "\n")
}

func main() {
	fmt.Println(analyzeVelocity())
}
