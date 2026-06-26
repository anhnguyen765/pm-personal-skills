package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

type AuditData struct {
	HowToUse  []string     `json:"how_to_use"`
	Checklist []AuditRow   `json:"checklist"`
	Summary   []SummaryRow `json:"summary"`
}

type AuditRow struct {
	RowIndex int    `json:"row_index"`
	Desc     string `json:"description"`
	What     string `json:"what_to_check"`
	Failure  string `json:"failure_example"`
}

type SummaryRow struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type WriteRequest struct {
	InputPath  string         `json:"input_path"`
	OutputPath string         `json:"output_path"`
	Results    []AuditResult `json:"results"`
	Summary    *AuditSummary `json:"summary"`
}

type AuditSummary struct {
	SquadProduct string `json:"squad_product"`
	Reviewer     string `json:"reviewer"`
	Date         string `json:"date"`
	Scope        string `json:"scope"`
	KeyRisks     string `json:"key_risks"`
	HighSeverity string `json:"high_severity"`
	MedSeverity  string `json:"med_severity"`
	LowSeverity  string `json:"low_severity"`
	NextSteps    string `json:"next_steps"`
	Questions    string `json:"questions"`
}

type AuditResult struct {
	RowIndex int    `json:"row_index"`
	Status   string `json:"status"`   // Col I: Pass, Flag, Fail, N/A
	Evidence string `json:"evidence"` // Col J: Evidence and Finding details
	Severity string `json:"severity"` // Col K: Low, Medium, High
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

func readExcel(path string) (*AuditData, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := &AuditData{}

	// Read How to Use
	rows, _ := f.GetRows("📋 How to Use")
	for _, row := range rows {
		if len(row) > 0 {
			data.HowToUse = append(data.HowToUse, row[0])
		}
	}

	// Read Checklist
	checklistRows, _ := f.GetRows("✅ Checklist")
	for i, row := range checklistRows {
		if i < 3 { // Skip headers (Rows 1-3, data starts from Row 4)
			continue
		}
		
		desc := ""
		what := ""
		failure := ""
		
		if len(row) > 5 { desc = row[5] }
		if len(row) > 6 { what = row[6] }
		if len(row) > 7 { failure = row[7] }

		data.Checklist = append(data.Checklist, AuditRow{
			RowIndex: i + 1,
			Desc:     desc,
			What:     what,
			Failure:  failure,
		})
	}

	// Read Summary
	summaryRows, _ := f.GetRows("📝 Summary")
	for _, row := range summaryRows {
		if len(row) > 1 {
			data.Summary = append(data.Summary, SummaryRow{
				Label: row[0],
				Value: row[1],
			})
		}
	}

	return data, nil
}

func writeExcel(req WriteRequest) error {
	f, err := excelize.OpenFile(req.InputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a map of provided results for quick lookup
	resultsMap := make(map[int]AuditResult)
	for _, res := range req.Results {
		resultsMap[res.RowIndex] = res
	}

	// Iterate through all possible audit rows (4 to 78)
	for i := 4; i <= 78; i++ {
		if res, ok := resultsMap[i]; ok {
			// Write provided result
			f.SetCellValue("✅ Checklist", fmt.Sprintf("I%d", i), res.Status)
			f.SetCellValue("✅ Checklist", fmt.Sprintf("J%d", i), res.Evidence)
			f.SetCellValue("✅ Checklist", fmt.Sprintf("K%d", i), res.Severity)
		} else {
			// Write default N/A
			f.SetCellValue("✅ Checklist", fmt.Sprintf("I%d", i), "N/A")
			f.SetCellValue("✅ Checklist", fmt.Sprintf("J%d", i), "Criteria not applicable for this specific flow/domain.")
			f.SetCellValue("✅ Checklist", fmt.Sprintf("K%d", i), "Low")
		}
	}

	// Write Summary if provided
	if req.Summary != nil {
		s := req.Summary
		f.SetCellValue("📝 Summary", "B2", fmt.Sprintf("Squad / Product: %s     Reviewer: %s     Date: %s", s.SquadProduct, s.Reviewer, s.Date))
		f.SetCellValue("📝 Summary", "B4", s.Scope)
		f.SetCellValue("📝 Summary", "B6", s.KeyRisks)
		f.SetCellValue("📝 Summary", "B8", s.HighSeverity)
		f.SetCellValue("📝 Summary", "B10", s.MedSeverity)
		f.SetCellValue("📝 Summary", "B12", s.LowSeverity)
		f.SetCellValue("📝 Summary", "B14", s.NextSteps)
		f.SetCellValue("📝 Summary", "B16", s.Questions)
	}

	return f.SaveAs(req.OutputPath)
}

func main() {
	loadEnv()
	mode := flag.String("mode", "read", "read or write")
	path := flag.String("path", "", "input file path")
	writeReq := flag.String("data", "", "JSON data for write mode")
	jsonFile := flag.String("json-file", "", "Path to JSON data file")
	flag.Parse()

	if *mode == "read" {
		data, err := readExcel(*path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		jsonBytes, _ := json.MarshalIndent(data, "", "  ")
		fmt.Println(string(jsonBytes))
	} else if *mode == "write" {
		var req WriteRequest
		var jsonBytes []byte
		var err error

		if *jsonFile != "" {
			jsonBytes, err = os.ReadFile(*jsonFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading JSON file: %v\n", err)
				os.Exit(1)
			}
		} else {
			jsonBytes = []byte(*writeReq)
		}

		if err := json.Unmarshal(jsonBytes, &req); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing write data: %v\n", err)
			os.Exit(1)
		}
		if err := writeExcel(req); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Success")
	}
}
