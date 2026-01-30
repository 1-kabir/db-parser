package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// EmailValidator handles email validation
type EmailValidator struct {
	emailRegex *regexp.Regexp
}

// NewEmailValidator creates a new email validator
func NewEmailValidator() *EmailValidator {
	// Email regex that prevents leading/trailing dots in local part
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._%+\-]*[a-zA-Z0-9]@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$|^[a-zA-Z0-9]@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return &EmailValidator{emailRegex: emailRegex}
}

// IsValid checks if an email is valid
func (ev *EmailValidator) IsValid(email string) bool {
	email = strings.TrimSpace(email)
	return ev.emailRegex.MatchString(email)
}

// FileProcessor handles processing of email files
type FileProcessor struct {
	validator      *EmailValidator
	separator      string
	globalEmails   map[string]bool
	emailMutex     sync.RWMutex
	progressFunc   func(string)
	workerCount    int
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(separator string, workerCount int, progressFunc func(string)) *FileProcessor {
	return &FileProcessor{
		validator:    NewEmailValidator(),
		separator:    separator,
		globalEmails: make(map[string]bool),
		progressFunc: progressFunc,
		workerCount:  workerCount,
	}
}

// ProcessingResult contains the result of processing a file
type ProcessingResult struct {
	OriginalPath string
	OutputPath   string
	InvalidPath  string
	ValidCount   int
	InvalidCount int
	Error        error
}

// ProcessFiles processes multiple files in parallel
func (fp *FileProcessor) ProcessFiles(inputDir string, files []string) []ProcessingResult {
	results := make([]ProcessingResult, len(files))
	
	// Create output directory
	outputDir := filepath.Join(inputDir, "outputs")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fp.sendProgress(fmt.Sprintf("Error creating output directory: %v", err))
		return results
	}

	// Create work queue
	jobs := make(chan int, len(files))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < fp.workerCount; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := range jobs {
				filePath := filepath.Join(inputDir, files[i])
				fp.sendProgress(fmt.Sprintf("Worker %d: Processing %s", workerID, files[i]))
				results[i] = fp.processFile(filePath, outputDir)
			}
		}(w)
	}

	// Queue all jobs
	for i := range files {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	return results
}

// processFile processes a single file
func (fp *FileProcessor) processFile(inputPath string, outputDir string) ProcessingResult {
	result := ProcessingResult{
		OriginalPath: inputPath,
	}

	// Open input file
	inFile, err := os.Open(inputPath)
	if err != nil {
		result.Error = err
		return result
	}
	defer inFile.Close()

	// Extract filename
	filename := filepath.Base(inputPath)
	
	// Prepare output files
	outputPath := filepath.Join(outputDir, filename)
	invalidPath := filepath.Join(outputDir, strings.TrimSuffix(filename, ".txt")+" [INVALID].txt")
	
	result.OutputPath = outputPath
	result.InvalidPath = invalidPath

	// Create output files
	outFile, err := os.Create(outputPath)
	if err != nil {
		result.Error = err
		return result
	}
	defer outFile.Close()

	invalidFile, err := os.Create(invalidPath)
	if err != nil {
		result.Error = err
		return result
	}
	defer invalidFile.Close()

	// Process line by line
	scanner := bufio.NewScanner(inFile)
	validLines := 0
	invalidLines := 0

	// Set larger buffer for scanner to handle long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		email, _ := fp.extractEmail(line)
		
		if email == "" {
			// No email found, skip line
			invalidLines++
			fmt.Fprintf(invalidFile, "%s [NO_EMAIL]\n", line)
			continue
		}

		// Validate email
		if !fp.validator.IsValid(email) {
			invalidLines++
			fmt.Fprintf(invalidFile, "%s [INVALID_FORMAT]\n", line)
			continue
		}

		// Check for duplicates with write lock to avoid race condition
		fp.emailMutex.Lock()
		if fp.globalEmails[email] {
			// Duplicate found
			fp.emailMutex.Unlock()
			invalidLines++
			fmt.Fprintf(invalidFile, "%s [DUPLICATE]\n", line)
			continue
		}
		fp.globalEmails[email] = true
		fp.emailMutex.Unlock()

		// Valid email, write to output
		validLines++
		fmt.Fprintln(outFile, line)
	}

	if err := scanner.Err(); err != nil {
		result.Error = err
		return result
	}

	result.ValidCount = validLines
	result.InvalidCount = invalidLines

	// Update filename with line count
	if validLines > 0 {
		newFilename := updateFilenameWithCount(filename, validLines)
		newOutputPath := filepath.Join(outputDir, newFilename)
		
		// Close the file before renaming
		outFile.Close()
		
		if newFilename != filename {
			if err := os.Rename(outputPath, newOutputPath); err != nil {
				log.Printf("Error renaming output file: %v", err)
				result.OutputPath = outputPath // Keep original name
			} else {
				result.OutputPath = newOutputPath
			}
		}
	} else {
		// No valid lines, remove the output file
		outFile.Close()
		if err := os.Remove(outputPath); err != nil {
			log.Printf("Error removing empty output file: %v", err)
		}
		result.OutputPath = ""
	}

	// Remove invalid file if empty
	invalidFile.Close()
	if invalidLines == 0 {
		if err := os.Remove(invalidPath); err != nil {
			log.Printf("Error removing empty invalid file: %v", err)
		}
		result.InvalidPath = ""
	}

	return result
}

// extractEmail extracts email from a line based on separator
func (fp *FileProcessor) extractEmail(line string) (string, string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", ""
	}

	// Split by separator
	parts := strings.SplitN(line, fp.separator, 2)
	email := strings.TrimSpace(parts[0])
	
	remainder := ""
	if len(parts) > 1 {
		remainder = parts[1]
	}

	return email, remainder
}

// updateFilenameWithCount updates filename with line count in curly braces
func updateFilenameWithCount(filename string, count int) string {
	// Remove extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Find rightmost curly braces
	regex := regexp.MustCompile(`\{[^}]*\}`)
	matches := regex.FindAllStringIndex(nameWithoutExt, -1)

	if len(matches) > 0 {
		// Replace rightmost match
		lastMatch := matches[len(matches)-1]
		before := nameWithoutExt[:lastMatch[0]]
		after := nameWithoutExt[lastMatch[1]:]
		return fmt.Sprintf("%s{%d}%s%s", before, count, after, ext)
	}

	// No curly braces found, append count
	return fmt.Sprintf("%s {%d}%s", nameWithoutExt, count, ext)
}

// sendProgress sends progress update
func (fp *FileProcessor) sendProgress(message string) {
	if fp.progressFunc != nil {
		fp.progressFunc(message)
	}
}
