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

type Issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Assignee struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		StoryPoints float64 `json:"customfield_10801"`
		Status      struct {
			Name           string `json:"name"`
			StatusCategory struct {
				Name string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
	} `json:"fields"`
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

func main() {
	loadEnv()
	jiraURL := os.Getenv("JIRA_URL")
	jiraToken := os.Getenv("JIRA_API_TOKEN")

	sprints := []struct {
		ID   int
		Name string
	}{
		{4862, "26.04.A"},
		{4767, "26.04.B"},
		{4768, "26.04.C"},
		{4825, "26.05.A"},
		{4826, "26.05.B"},
		{4827, "26.06.A"},
	}

	// Track issue presence across sprints
	issueMap := make(map[string][]string) // Key -> []SprintNames
	issueData := make(map[string]Issue)
	
	// Velocity Tracking
	donePoints := make(map[string]map[string]float64) // member -> sprint -> points
	carryPoints := make(map[string]map[string]float64) // member -> sprint -> points
	members := make(map[string]bool)

	client := &http.Client{}

	for _, s := range sprints {
		url := fmt.Sprintf("%s/rest/api/2/search?jql=sprint=%d&fields=summary,assignee,customfield_10801,status&maxResults=200", jiraURL, s.ID)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+jiraToken)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var data SearchResponse
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &data)

		for _, issue := range data.Issues {
			issueMap[issue.Key] = append(issueMap[issue.Key], s.Name)
			issueData[issue.Key] = issue

			assignee := issue.Fields.Assignee.DisplayName
			if assignee == "" {
				assignee = "Unassigned"
			}
			members[assignee] = true

			if issue.Fields.Status.StatusCategory.Name == "Done" {
				if donePoints[assignee] == nil {
					donePoints[assignee] = make(map[string]float64)
				}
				donePoints[assignee][s.Name] += issue.Fields.StoryPoints
			} else {
				// Task is in the sprint but NOT done
				if carryPoints[assignee] == nil {
					carryPoints[assignee] = make(map[string]float64)
				}
				carryPoints[assignee][s.Name] += issue.Fields.StoryPoints
			}
		}
	}

	// 1. Velocity Analysis (Done Points)
	fmt.Println("### Velocity Analysis (Completed Story Points)")
	var memberList []string
	for m := range members {
		memberList = append(memberList, m)
	}
	sort.Strings(memberList)

	fmt.Printf("| %-25s |", "Member")
	for _, s := range sprints {
		fmt.Printf(" %-7s |", s.Name)
	}
	fmt.Println("\n|---------------------------|" + strings.Repeat("---------|", len(sprints)))

	for _, m := range memberList {
		fmt.Printf("| %-25s |", m)
		for _, s := range sprints {
			fmt.Printf(" %-7.1f |", donePoints[m][s.Name])
		}
		fmt.Println()
	}

	// 2. Carry-Over Analysis (Not Done Points)
	fmt.Println("\n### Carry-Over Analysis (In-Progress/Spillover Points)")
	fmt.Printf("| %-25s |", "Member")
	for _, s := range sprints {
		fmt.Printf(" %-7s |", s.Name)
	}
	fmt.Println("\n|---------------------------|" + strings.Repeat("---------|", len(sprints)))

	for _, m := range memberList {
		fmt.Printf("| %-25s |", m)
		for _, s := range sprints {
			fmt.Printf(" %-7.1f |", carryPoints[m][s.Name])
		}
		fmt.Println()
	}

	// 3. Persistent Carry-Over Tasks (Tasks that appeared in 3+ sprints)
	fmt.Println("\n### Persistent Carry-Over Tasks (3+ Sprints)")
	fmt.Printf("| %-10s | %-15s | %-10s | %-20s | %-40s |\n", "Key", "Assignee", "Status", "Sprints Seen", "Summary")
	fmt.Println("|------------|-----------------|------------|----------------------|------------------------------------------|")
	
	var persistentKeys []string
	for key, sprintsSeen := range issueMap {
		if len(sprintsSeen) >= 3 {
			persistentKeys = append(persistentKeys, key)
		}
	}
	sort.Strings(persistentKeys)

	for _, key := range persistentKeys {
		issue := issueData[key]
		sprintsSeen := strings.Join(issueMap[key], ", ")
		assignee := issue.Fields.Assignee.DisplayName
		if assignee == "" {
			assignee = "Unassigned"
		}
		summary := issue.Fields.Summary
		if len(summary) > 40 {
			summary = summary[:37] + "..."
		}
		fmt.Printf("| %-10s | %-15s | %-10s | %-20s | %-40s |\n", key, assignee, issue.Fields.Status.Name, sprintsSeen, summary)
	}
}
