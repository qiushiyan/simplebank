// This programs formats the output of `kubectl get` commands into a table
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var headers []string
	var rows [][]string
	var currentComponentType string // Track the current component type

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")

		if strings.HasPrefix(line, "kubectl get") {
			// Process the current component type before resetting for the new command
			if len(rows) > 0 {
				printTable(headers, rows, currentComponentType)
				rows = [][]string{} // Reset rows for the next section
			}

			// Extract and print the component type in blue
			re := regexp.MustCompile(`kubectl get\s+(\S+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				componentType := matches[1]
				fmt.Println("\033[0;34m" + componentType + "\033[0m")
				currentComponentType = componentType // Update current component type
			}
			continue
		}

		if len(line) > 0 {
			fields := strings.Fields(line)

			if strings.HasPrefix(fields[0], "NAME") {
				// Adjust and set new headers
				headers = adjustColumns(fields)
			} else if len(headers) > 0 {
				// Adjust rows to match the headers and collect them
				row := adjustColumns(fields)
				rows = append(rows, row)
			}
		}
	}

	// Print any remaining rows as a table for the last component type
	if len(rows) > 0 {
		printTable(headers, rows, currentComponentType)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
	}
}

func adjustColumns(columns []string) []string {
	if len(columns) > 8 {
		return columns[:8]
	}
	return columns
}

func printTable(headers []string, rows [][]string, componentType string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	for _, row := range rows {
		table.Append(row)
	}
	colors := make([]tablewriter.Colors, len(headers))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold}
	}
	table.SetHeaderColor(colors...)
	table.Render()
}
