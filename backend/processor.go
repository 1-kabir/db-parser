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
	validateEmail  bool
	outputInvalid  bool
	deleteAfterSep bool
}

// NewFileProcessor creates a new file processor
func NewFileProcessor(separator string, workerCount int, validateEmail bool, outputInvalid bool, deleteAfterSep bool, progressFunc func(string)) *FileProcessor {
	return &FileProcessor{
		validator:      NewEmailValidator(),
		separator:      separator,
		globalEmails:   make(map[string]bool),
		progressFunc:   progressFunc,
		workerCount:    workerCount,
		validateEmail:  validateEmail,
		outputInvalid:  outputInvalid,
		deleteAfterSep: deleteAfterSep,
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
func (fp *FileProcessor) ProcessFiles(inputDir string, files []string, outputDirName string) []ProcessingResult {
	results := make([]ProcessingResult, len(files))
	
	// Create output directory
	var outputDir string
	if filepath.IsAbs(outputDirName) {
		outputDir = outputDirName
	} else {
		outputDir = filepath.Join(inputDir, outputDirName)
	}
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fp.sendProgress(fmt.Sprintf("Error creating output directory '%s': %v. Please check write permissions.", outputDir, err))
		// Set error for all files
		for i := range results {
			results[i].Error = fmt.Errorf("failed to create output directory: %v", err)
			results[i].OriginalPath = filepath.Join(inputDir, files[i])
		}
		return results
	}
	fp.sendProgress(fmt.Sprintf("Output directory created: %s", outputDir))

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
		result.Error = fmt.Errorf("failed to open file: %v. Check if the file exists and you have read permissions", err)
		return result
	}
	defer inFile.Close()

	// Extract filename
	filename := filepath.Base(inputPath)
	
	// Prepare output files
	outputPath := filepath.Join(outputDir, filename)
	invalidPath := filepath.Join(outputDir, strings.TrimSuffix(filename, filepath.Ext(filename))+" [INVALID]"+filepath.Ext(filename))
	
	result.OutputPath = outputPath
	result.InvalidPath = invalidPath

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		result.Error = fmt.Errorf("failed to create output file: %v. Check if you have write permissions", err)
		return result
	}
	defer outFile.Close()

	// Create invalid file only if outputInvalid is true
	var invalidFile *os.File
	if fp.outputInvalid {
		invalidFile, err = os.Create(invalidPath)
		if err != nil {
			result.Error = fmt.Errorf("failed to create invalid file: %v. Check if you have write permissions", err)
			return result
		}
		defer invalidFile.Close()
	}

	// Process line by line
	scanner := bufio.NewScanner(inFile)
	validLines := 0
	invalidLines := 0

	// Set larger buffer for scanner to handle long lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		originalLine := line
		
		// If deleteAfterSep mode, remove separator and everything after
		if fp.deleteAfterSep {
			if idx := strings.Index(line, fp.separator); idx != -1 {
				line = line[:idx]
			}
		}
		
		email, _ := fp.extractEmail(originalLine)
		
		if email == "" {
			// No email found
			invalidLines++
			if fp.outputInvalid && invalidFile != nil {
				fmt.Fprintf(invalidFile, "%s [NO_EMAIL]\n", originalLine)
			}
			continue
		}

		// Validate email if validation is enabled
		if fp.validateEmail && !fp.validator.IsValid(email) {
			invalidLines++
			if fp.outputInvalid && invalidFile != nil {
				fmt.Fprintf(invalidFile, "%s [INVALID_FORMAT]\n", originalLine)
			}
			continue
		}

		// Check for duplicates with write lock to avoid race condition
		fp.emailMutex.Lock()
		if fp.globalEmails[email] {
			// Duplicate found
			fp.emailMutex.Unlock()
			invalidLines++
			if fp.outputInvalid && invalidFile != nil {
				fmt.Fprintf(invalidFile, "%s [DUPLICATE]\n", originalLine)
			}
			continue
		}
		fp.globalEmails[email] = true
		fp.emailMutex.Unlock()

		// Valid email, write to output
		validLines++
		if fp.deleteAfterSep {
			// Write only the part before separator
			fmt.Fprintln(outFile, line)
		} else {
			// Write the full line
			fmt.Fprintln(outFile, originalLine)
		}
	}

	if err := scanner.Err(); err != nil {
		result.Error = fmt.Errorf("error reading file: %v", err)
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
				log.Printf("Warning: Could not rename output file: %v", err)
				result.OutputPath = outputPath // Keep original name
			} else {
				result.OutputPath = newOutputPath
			}
		}
	} else {
		// No valid lines, remove the output file
		outFile.Close()
		if err := os.Remove(outputPath); err != nil {
			log.Printf("Warning: Could not remove empty output file: %v", err)
		}
		result.OutputPath = ""
	}

	// Remove invalid file if empty or if not outputting invalid files
	if invalidFile != nil {
		invalidFile.Close()
		if invalidLines == 0 {
			if err := os.Remove(invalidPath); err != nil {
				log.Printf("Warning: Could not remove empty invalid file: %v", err)
			}
			result.InvalidPath = ""
		}
	} else {
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

// formatNumberWithCommas formats a number with comma separators (international format)
func formatNumberWithCommas(n int) string {
	str := fmt.Sprintf("%d", n)
	if n < 1000 {
		return str
	}
	
	// Insert commas from right to left
	var result strings.Builder
	for i, digit := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}
	return result.String()
}

// updateFilenameWithCount updates filename with line count in curly braces
// Removes ALL existing curly bracket content and adds formatted count at the end
func updateFilenameWithCount(filename string, count int) string {
	// Remove extension
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Remove ALL curly brackets and their content
	regex := regexp.MustCompile(`\{[^}]*\}`)
	cleanedName := regex.ReplaceAllString(nameWithoutExt, "")
	
	// Clean up multiple consecutive spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	cleanedName = spaceRegex.ReplaceAllString(cleanedName, " ")
	
	// Trim any leading/trailing spaces
	cleanedName = strings.TrimSpace(cleanedName)
	
	// Add the formatted count in curly braces
	formattedCount := formatNumberWithCommas(count)
	return fmt.Sprintf("%s {%s}%s", cleanedName, formattedCount, ext)
}

// sendProgress sends progress update
func (fp *FileProcessor) sendProgress(message string) {
	if fp.progressFunc != nil {
		fp.progressFunc(message)
	}
}
