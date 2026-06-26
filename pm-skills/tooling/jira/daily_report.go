package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Sprint struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SprintResponse struct {
	Values []Sprint `json:"values"`
}

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name           string `json:"name"`
			StatusCategory struct {
				Name string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Created string `json:"created"`
	} `json:"fields"`
	Changelog struct {
		Histories []struct {
			Created string `json:"created"`
			Items   []struct {
				Field string `json:"field"`
			} `json:"items"`
		} `json:"histories"`
	} `json:"changelog"`
}

type SearchResponse struct {
	Issues []Issue `json:"issues"`
}

func loadEnv() {
	current, _ := os.Getwd()
	for {
		envFile := filepath.Join(current, ".env")
		if _, err := os.Stat(envFile); err == nil {
			file, _ := os.Open(envFile)
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
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
}

func findProjectRoot() string {
	current, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(current, ".env")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	cwd, _ := os.Getwd()
	return cwd
}

func getActiveSprint(jiraURL, jiraToken string, boardID int) (*Sprint, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", jiraURL, boardID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+jiraToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data SprintResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	if len(data.Values) > 0 {
		return &data.Values[0], nil
	}
	return nil, fmt.Errorf("no active sprint found")
}

func parseJiraTime(tStr string) time.Time {
	// Jira format: 2026-05-12T09:44:03.123+0700
	// We'll simplify and use a common subset
	layout := "2006-01-02T15:04:05.000-0700"
	t, err := time.Parse(layout, tStr)
	if err != nil {
		// Fallback to simpler format if needed
		layoutSimple := "2006-01-02T15:04:05.000Z0700"
		t, _ = time.Parse(layoutSimple, tStr)
	}
	return t
}

func generateReport() (string, error) {
	loadEnv()
	jiraURL := os.Getenv("JIRA_URL")
	jiraToken := os.Getenv("JIRA_API_TOKEN")
	boardID := 1131

	sprint, err := getActiveSprint(jiraURL, jiraToken, boardID)
	if err != nil {
		return "", err
	}

	searchURL := fmt.Sprintf("%s/rest/api/2/search?jql=sprint=%d&expand=changelog&fields=summary,status,assignee,created&maxResults=100", jiraURL, sprint.ID)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("Authorization", "Bearer "+jiraToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data SearchResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &data)

	now := time.Now()
	reportLines := []string{
		fmt.Sprintf("# Daily Standup Report: %s", sprint.Name),
		fmt.Sprintf("Generated at: %s\n", now.Format("2006-01-02 15:04:05")),
	}

	assigneeTasks := make(map[string][]string)
	type RiskTask struct {
		Key      string
		Assignee string
		Status   string
		Age      int
		Summary  string
	}
	var riskTasks []RiskTask
	doneCount := 0

	for _, issue := range data.Issues {
		status := issue.Fields.Status.Name
		cat := issue.Fields.Status.StatusCategory.Name
		assignee := issue.Fields.Assignee.DisplayName
		if assignee == "" {
			assignee = "Unassigned"
		}

		assigneeTasks[assignee] = append(assigneeTasks[assignee], fmt.Sprintf("**%s** [%s]", issue.Key, status))

		if cat == "Done" || strings.ToUpper(status) == "DONE" || strings.ToUpper(status) == "CLOSED" || strings.ToUpper(status) == "RESOLVED" {
			doneCount++
		}

		// Age in status
		lastChangeStr := issue.Fields.Created
		for i := len(issue.Changelog.Histories) - 1; i >= 0; i-- {
			history := issue.Changelog.Histories[i]
			statusChanged := false
			for _, item := range history.Items {
				if item.Field == "status" {
					statusChanged = true
					break
				}
			}
			if statusChanged {
				lastChangeStr = history.Created
				break
			}
		}

		// Simplified age calculation
		var age int
		lastChange, err := time.Parse("2006-01-02T15:04:05.000-0700", lastChangeStr)
		if err != nil {
			lastChange, _ = time.Parse("2006-01-02T15:04:05.000Z", lastChangeStr[:strings.Index(lastChangeStr, ".")+4]+"Z")
		}
		
		age = int(now.Sub(lastChange).Hours() / 24)

		if age > 3 && cat != "Done" && strings.ToUpper(status) != "NEW" {
			riskTasks = append(riskTasks, RiskTask{
				Key:      issue.Key,
				Assignee: assignee,
				Status:   status,
				Age:      age,
				Summary:  issue.Fields.Summary,
			})
		}
	}

	total := len(data.Issues)
	percent := 0.0
	if total > 0 {
		percent = float64(doneCount) / float64(total) * 100
	}

	reportLines = append(reportLines, fmt.Sprintf("## 1. Summary\n- **Total Tasks:** %d\n- **Completed:** %d (%.1f%%)\n", total, doneCount, percent))
	reportLines = append(reportLines, "## 2. Tasks by Assignee")
	
	var assignees []string
	for a := range assigneeTasks {
		assignees = append(assignees, a)
	}
	sort.Strings(assignees)

	for _, a := range assignees {
		tasks := assigneeTasks[a]
		sort.Strings(tasks)
		reportLines = append(reportLines, fmt.Sprintf("### %s (%d)", a, len(tasks)))
		for _, t := range tasks {
			reportLines = append(reportLines, "- "+t)
		}
	}

	reportLines = append(reportLines, "\n## 3. Stuck Tasks (> 3 Days in Status, excluding NEW)")
	if len(riskTasks) == 0 {
		reportLines = append(reportLines, "No tasks stuck for more than 3 days.")
	} else {
		sort.Slice(riskTasks, func(i, j int) bool {
			return riskTasks[i].Age > riskTasks[j].Age
		})
		reportLines = append(reportLines, "| Key | Assignee | Status | Days | Summary |")
		reportLines = append(reportLines, "| :--- | :--- | :--- | :--- | :--- |")
		for _, rt := range riskTasks {
			summ := rt.Summary
			if len(summ) > 50 {
				summ = summ[:50]
			}
			reportLines = append(reportLines, fmt.Sprintf("| %s | %s | %s | %d | %s |", rt.Key, rt.Assignee, rt.Status, rt.Age, summ))
		}
	}

	return strings.Join(reportLines, "\n"), nil
}

func main() {
	sandbox := flag.Bool("sandbox", false, "Save to the local sandbox instead of OneDrive")
	flag.Parse()

	content, err := generateReport()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("daily_report_%s.md", time.Now().Format("20060102"))
	var outputDir string
	if *sandbox {
		outputDir = filepath.Join(findProjectRoot(), "pm-skills", "output", "sandbox")
	} else {
		outputDir = "/Users/lap14569/Library/CloudStorage/OneDrive-VNGGroupJSC/vng/Daily"
	}

	os.MkdirAll(outputDir, 0755)
	filepath := filepath.Join(outputDir, filename)
	os.WriteFile(filepath, []byte(content), 0644)

	fmt.Printf("Report generated successfully: %s\n", filepath)
}
