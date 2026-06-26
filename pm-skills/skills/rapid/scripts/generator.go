package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary string `json:"summary"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
}

type TaskEntry struct {
	Number  int
	Task    string
	Status  string
	Notes   string
	JiraKey string
}

type ReportData struct {
	LastWeekTasks    []TaskEntry
	ThisWeekTasks    []TaskEntry
	SupportNeeded    []TaskEntry
	Owner            string
	ReportDate       string
}

func main() {
	fmt.Println("🔄 Weekly Report Generator")
	fmt.Println(strings.Repeat("=", 50))

	today := time.Now().Format("2006/01/02")
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Fetch latest report from Confluence
	fmt.Println("\n📚 Step 1: Fetching latest report from Confluence...")
	confluencePath := "https://confluence.zalopay.vn/x/kSB6Eg"
	fmt.Printf("Searching for latest report in: %s\n", confluencePath)

	// Step 2: Pull Jira activities
	fmt.Println("\n⚙️  Step 2: Pulling Jira activities from past week...")
	jiraIssues := fetchJiraActivities()
	fmt.Printf("Found %d Jira issues from past week\n", len(jiraIssues))

	// Step 3: Populate last week's tasks
	fmt.Println("\n📝 Step 3: What I Did Last Week")
	fmt.Println("We'll review your Jira activities. Edit status and notes as needed.")
	lastWeekTasks := gatherLastWeekTasks(jiraIssues, reader)

	// Step 4: Gather this week's tasks
	fmt.Println("\n📋 Step 4: What I Will Do This Week")
	thisWeekTasks := gatherThisWeekTasks(reader)

	// Step 5: Gather support needed
	fmt.Println("\n🆘 Step 5: Where I Need Support")
	supportNeeded := gatherSupportNeeded(reader)

	// Step 6: Generate report
	fmt.Println("\n✨ Step 6: Generating report...")
	report := generateReport(ReportData{
		LastWeekTasks: lastWeekTasks,
		ThisWeekTasks: thisWeekTasks,
		SupportNeeded: supportNeeded,
		Owner:         "Anh Nguyễn Hà Minh (3)",
		ReportDate:    today,
	})

	// Save to sandbox
	filename := fmt.Sprintf("pm-skills/output/sandbox/weekly-report-%s.md", time.Now().Format("2006-01-02"))
	err := os.WriteFile(filename, []byte(report), 0644)
	if err != nil {
		fmt.Printf("❌ Error saving report: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Report saved to: %s\n", filename)
	fmt.Println("\n📋 Next: Copy to Confluence at Payment Platform > Meeting Notes/2026/PO Weekly Reports/Anhnhm3")
}

func fetchJiraActivities() []JiraIssue {
	// Calculate date from 7 days ago
	past := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	jql := fmt.Sprintf("assignee = currentUser() AND updated >= %s ORDER BY updated DESC", past)

	cmd := exec.Command("./pm-skills/tooling/jira/jira", "search", jql, "--limit", "50")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Could not fetch Jira issues: %v\n", err)
		return []JiraIssue{}
	}

	var result struct {
		Issues []JiraIssue `json:"issues"`
	}

	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Printf("⚠️  Could not parse Jira response: %v\n", err)
		return []JiraIssue{}
	}

	return result.Issues
}

func gatherLastWeekTasks(issues []JiraIssue, reader *bufio.Reader) []TaskEntry {
	tasks := []TaskEntry{}

	// Auto-populate from Jira
	for i, issue := range issues {
		task := TaskEntry{
			Number:  i + 1,
			Task:    issue.Fields.Summary,
			Status:  issue.Fields.Status.Name,
			JiraKey: issue.Key,
		}

		fmt.Printf("\n[%d] %s (%s)\n", task.Number, task.Task, issue.Key)
		fmt.Print("    Status (DONE/ONGOING/TBD/TBC)? [default: " + task.Status + "]: ")
		status, _ := reader.ReadString('\n')
		status = strings.TrimSpace(status)
		if status != "" {
			task.Status = status
		}

		fmt.Print("    Notes (press Enter to skip): ")
		notes, _ := reader.ReadString('\n')
		task.Notes = strings.TrimSpace(notes)

		tasks = append(tasks, task)
	}

	// Allow adding more tasks not in Jira
	fmt.Print("\nAdd more tasks not in Jira? (y/n): ")
	addMore, _ := reader.ReadString('\n')
	if strings.TrimSpace(addMore) == "y" {
		for {
			fmt.Print("\nTask name (or 'done'): ")
			taskName, _ := reader.ReadString('\n')
			taskName = strings.TrimSpace(taskName)
			if taskName == "done" {
				break
			}

			fmt.Print("Status: ")
			status, _ := reader.ReadString('\n')

			fmt.Print("Notes: ")
			notes, _ := reader.ReadString('\n')

			tasks = append(tasks, TaskEntry{
				Number: len(tasks) + 1,
				Task:   taskName,
				Status: strings.TrimSpace(status),
				Notes:  strings.TrimSpace(notes),
			})
		}
	}

	return tasks
}

func gatherThisWeekTasks(reader *bufio.Reader) []TaskEntry {
	tasks := []TaskEntry{}

	fmt.Println("\nEnter tasks for this week (or 'done' to finish):")
	for {
		fmt.Print("\nTask name (or 'done'): ")
		taskName, _ := reader.ReadString('\n')
		taskName = strings.TrimSpace(taskName)
		if taskName == "done" {
			break
		}

		fmt.Print("Status (TODO/ONGOING/TBD/TBC): ")
		status, _ := reader.ReadString('\n')

		fmt.Print("Notes: ")
		notes, _ := reader.ReadString('\n')

		tasks = append(tasks, TaskEntry{
			Number: len(tasks) + 1,
			Task:   taskName,
			Status: strings.TrimSpace(status),
			Notes:  strings.TrimSpace(notes),
		})
	}

	return tasks
}

func gatherSupportNeeded(reader *bufio.Reader) []TaskEntry {
	tasks := []TaskEntry{}

	fmt.Println("\nEnter blockers/support needed (or 'done' to finish):")
	for {
		fmt.Print("\nTask/Blocker (or 'done'): ")
		taskName, _ := reader.ReadString('\n')
		taskName = strings.TrimSpace(taskName)
		if taskName == "done" {
			break
		}

		fmt.Print("Status (TBD/TBC): ")
		status, _ := reader.ReadString('\n')

		fmt.Print("Details (who to ask, timeline, etc.): ")
		notes, _ := reader.ReadString('\n')

		tasks = append(tasks, TaskEntry{
			Number: len(tasks) + 1,
			Task:   taskName,
			Status: strings.TrimSpace(status),
			Notes:  strings.TrimSpace(notes),
		})
	}

	return tasks
}

func generateReport(data ReportData) string {
	var sb strings.Builder

	// Header
	sb.WriteString(fmt.Sprintf("# %s - Weekly Report - %s\n\n", data.Owner, data.ReportDate))
	sb.WriteString(fmt.Sprintf("**Owner:** %s  \n", data.Owner))
	sb.WriteString("**Status:** DONE\n\n")
	sb.WriteString("---\n\n")

	// Section 1: What I Did Last Week
	sb.WriteString("# 1. What I Did Last Week\n\n")
	sb.WriteString("| # | Task | Status | Notes |\n")
	sb.WriteString("|---|------|--------|--------|\n")
	for _, task := range data.LastWeekTasks {
		task.Task = strings.ReplaceAll(task.Task, "|", "\\|") // Escape pipe characters
		task.Notes = strings.ReplaceAll(task.Notes, "|", "\\|")
		task.Notes = strings.ReplaceAll(task.Notes, "\n", "<br>") // Convert newlines to br tags
		sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", task.Number, task.Task, task.Status, task.Notes))
	}

	sb.WriteString("\n---\n\n")

	// Section 2: What I Will Do This Week
	sb.WriteString("# 2. What I Will Do This Week\n\n")
	sb.WriteString("| # | Task | Status | Notes |\n")
	sb.WriteString("|---|------|--------|--------|\n")
	for _, task := range data.ThisWeekTasks {
		task.Task = strings.ReplaceAll(task.Task, "|", "\\|")
		task.Notes = strings.ReplaceAll(task.Notes, "|", "\\|")
		task.Notes = strings.ReplaceAll(task.Notes, "\n", "<br>")
		sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", task.Number, task.Task, task.Status, task.Notes))
	}

	sb.WriteString("\n---\n\n")

	// Section 3: Where I Need Support
	sb.WriteString("# 3. Where I Need Support\n\n")
	sb.WriteString("| # | Task | Status | Notes |\n")
	sb.WriteString("|---|------|--------|--------|\n")
	for _, task := range data.SupportNeeded {
		task.Task = strings.ReplaceAll(task.Task, "|", "\\|")
		task.Notes = strings.ReplaceAll(task.Notes, "|", "\\|")
		task.Notes = strings.ReplaceAll(task.Notes, "\n", "<br>")
		sb.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", task.Number, task.Task, task.Status, task.Notes))
	}

	return sb.String()
}
