package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func findProjectRoot() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(current, ".env")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("project root not found")
}

func saveMemo(content, title string, sandbox bool) (string, error) {
	dateStr := time.Now().Format("2006-01-02")
	cleanTitle := strings.ReplaceAll(strings.ToLower(title), " ", "_")
	filename := fmt.Sprintf("ZLP_%s_memo_%s.md", dateStr, cleanTitle)

	var outputDir string
	if sandbox {
		root, err := findProjectRoot()
		if err != nil {
			// Fallback logic similar to Python script
			cwd, _ := os.Getwd()
			outputDir = filepath.Join(cwd, "pm-skills", "output", "sandbox")
		} else {
			outputDir = filepath.Join(root, "pm-skills", "output", "sandbox")
		}
	} else {
		outputDir = "/Users/lap14569/Library/CloudStorage/OneDrive-VNGGroupJSC/vng/Notes"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}

	filepath := filepath.Join(outputDir, filename)
	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

func main() {
	loadEnv()
	sandbox := flag.Bool("sandbox", false, "Save to the local sandbox instead of OneDrive")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: save_memo [--sandbox] <title> <content>")
		os.Exit(1)
	}

	title := args[0]
	content := args[1]

	path, err := saveMemo(content, title, *sandbox)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving memo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Memo saved successfully: %s\n", path)
}
