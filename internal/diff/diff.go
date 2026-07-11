package diff

import (
	"fmt"
	"os"
	"strings"
)

type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeModified ChangeType = "modified"
	ChangeUnchanged ChangeType = "unchanged"
)

type FileDiff struct {
	Path     string
	Type     ChangeType
	OldSize  int
	NewSize  int
	Lines    []DiffLine
}

type DiffLine struct {
	OldNum int
	NewNum int
	Type   ChangeType
	Content string
}

func ComputeDiff(oldContent, newContent string, path string) *FileDiff {
	if oldContent == "" && newContent != "" {
		return &FileDiff{
			Path:    path,
			Type:    ChangeAdded,
			NewSize: len(newContent),
			Lines:   addedLines(newContent),
		}
	}
	if oldContent != "" && newContent == "" {
		return &FileDiff{
			Path:    path,
			Type:    ChangeRemoved,
			OldSize: len(oldContent),
			Lines:   removedLines(oldContent),
		}
	}
	if oldContent == newContent {
		return &FileDiff{
			Path:     path,
			Type:     ChangeUnchanged,
			OldSize:  len(oldContent),
			NewSize:  len(newContent),
		}
	}

	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")
	diffLines := computeLineDiff(oldLines, newLines)

	return &FileDiff{
		Path:    path,
		Type:    ChangeModified,
		OldSize: len(oldContent),
		NewSize: len(newContent),
		Lines:   diffLines,
	}
}

func addedLines(content string) []DiffLine {
	var lines []DiffLine
	for i, line := range strings.Split(content, "\n") {
		lines = append(lines, DiffLine{
			OldNum:  0,
			NewNum:  i + 1,
			Type:    ChangeAdded,
			Content: line,
		})
	}
	return lines
}

func removedLines(content string) []DiffLine {
	var lines []DiffLine
	for i, line := range strings.Split(content, "\n") {
		lines = append(lines, DiffLine{
			OldNum:  i + 1,
			NewNum:  0,
			Type:    ChangeRemoved,
			Content: line,
		})
	}
	return lines
}

func computeLineDiff(oldLines, newLines []string) []DiffLine {
	var result []DiffLine
	maxLen := len(oldLines)
	if len(newLines) > maxLen {
		maxLen = len(newLines)
	}

	for i := 0; i < maxLen; i++ {
		var oldLine, newLine string
		oldNum := 0
		newNum := 0

		if i < len(oldLines) {
			oldLine = oldLines[i]
			oldNum = i + 1
		}
		if i < len(newLines) {
			newLine = newLines[i]
			newNum = i + 1
		}

		if oldLine == newLine {
			result = append(result, DiffLine{
				OldNum:  oldNum,
				NewNum:  newNum,
				Type:    ChangeUnchanged,
				Content: oldLine,
			})
		} else {
			if oldNum > 0 {
				result = append(result, DiffLine{
					OldNum:  oldNum,
					NewNum:  0,
					Type:    ChangeRemoved,
					Content: oldLine,
				})
			}
			if newNum > 0 {
				result = append(result, DiffLine{
					OldNum:  0,
					NewNum:  newNum,
					Type:    ChangeAdded,
					Content: newLine,
				})
			}
		}
	}
	return result
}

func ComputeDirectoryDiff(oldDir, newDir string, paths []string) []*FileDiff {
	var diffs []*FileDiff
	for _, path := range paths {
		oldContent := readFileIfExists(oldDir, path)
		newContent := readFileIfExists(newDir, path)
		diffs = append(diffs, ComputeDiff(oldContent, newContent, path))
	}
	return diffs
}

func readFileIfExists(dir, path string) string {
	fullPath := fmt.Sprintf("%s/%s", strings.TrimRight(dir, "/"), strings.TrimLeft(path, "/"))
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}
	return string(data)
}

func FormatDiff(diff *FileDiff) string {
	if diff == nil {
		return ""
	}
	var sb strings.Builder

	switch diff.Type {
	case ChangeAdded:
		sb.WriteString(fmt.Sprintf("\033[32m+++ %s (added)\033[0m\n", diff.Path))
	case ChangeRemoved:
		sb.WriteString(fmt.Sprintf("\033[31m--- %s (removed)\033[0m\n", diff.Path))
	case ChangeUnchanged:
		sb.WriteString(fmt.Sprintf("    %s (unchanged)\n", diff.Path))
		return sb.String()
	case ChangeModified:
		sb.WriteString(fmt.Sprintf("\033[33m~~~ %s (modified: %d -> %d bytes)\033[0m\n", diff.Path, diff.OldSize, diff.NewSize))
	}

	for _, line := range diff.Lines {
		switch line.Type {
		case ChangeAdded:
			sb.WriteString(fmt.Sprintf("\033[32m+ %s\033[0m\n", line.Content))
		case ChangeRemoved:
			sb.WriteString(fmt.Sprintf("\033[31m- %s\033[0m\n", line.Content))
		default:
			sb.WriteString(fmt.Sprintf("  %s\n", line.Content))
		}
	}
	return sb.String()
}

func Summary(diffs []*FileDiff) (added, removed, modified, unchanged int) {
	for _, d := range diffs {
		switch d.Type {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			modified++
		case ChangeUnchanged:
			unchanged++
		}
	}
	return
}
