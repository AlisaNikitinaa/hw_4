package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	count := flag.Bool("c", false, "Показать количество повторений каждой строки")
	duplicate := flag.Bool("d", false, "Показать только строки которые повторяются")
	unique := flag.Bool("u", false, "Показать только уникальные строки")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр")
	numFields := flag.Int("f", 0, "Пропустить первые n слов в каждой строке")
	numChars := flag.Int("s", 0, "Пропустить первые n символов в каждой строке")

	flag.Parse()

	if (*count && *duplicate) || (*count && *unique) || (*duplicate && *unique) {
		fmt.Println("Ошибка, так как параметры -c, -d, -u нельзя использовать вместе")
		os.Exit(1)
	}

	args := flag.Args()
	var inputFile, outputFile string

	if len(args) > 0 {
		inputFile = args[0]
	}
	if len(args) > 1 {
		outputFile = args[1]
	}

	lines := readLines(inputFile)

	result := processLines(lines, *count, *duplicate, *unique, *ignoreCase, *numFields, *numChars)

	writeLines(outputFile, result)
}

func readLines(filename string) []string {
	var lines []string
	var scanner *bufio.Scanner

	if filename == "" {

		fmt.Println("Введите строки (Ctrl+D или Ctrl+Z для завершения):")
		scanner = bufio.NewScanner(os.Stdin)
	} else {

		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Ошибка открытия файла %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	return lines
}

func writeLines(filename string, lines []string) {
	var output *os.File

	if filename == "" {

		output = os.Stdout
	} else {

		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Ошибка создания файла %s: %v\n", filename, err)
			os.Exit(1)
		}
		defer file.Close()
		output = file
	}

	for _, line := range lines {
		fmt.Fprintln(output, line)
	}
}

func processLines(lines []string, count, duplicate, unique, ignoreCase bool, numFields, numChars int) []string {
	if len(lines) == 0 {
		return lines
	}

	lineCounts := make(map[string]int)
	originalLines := make(map[string]string)

	for _, line := range lines {

		key := prepareKey(line, ignoreCase, numFields, numChars)

		if _, exists := originalLines[key]; !exists {
			originalLines[key] = line
		}

		lineCounts[key]++
	}

	var result []string
	processed := make(map[string]bool)
	for _, line := range lines {
		key := prepareKey(line, ignoreCase, numFields, numChars)

		if processed[key] {
			continue
		}
		processed[key] = true

		countValue := lineCounts[key]
		originalLine := originalLines[key]

		switch {
		case count:

			result = append(result, fmt.Sprintf("%d %s", countValue, originalLine))
		case duplicate && countValue > 1:

			result = append(result, originalLine)
		case unique && countValue == 1:

			result = append(result, originalLine)
		case !count && !duplicate && !unique:

			result = append(result, originalLine)
		}
	}

	return result
}

func prepareKey(line string, ignoreCase bool, numFields, numChars int) string {
	result := line

	if numFields > 0 {
		result = skipFields(result, numFields)
	}

	if numChars > 0 {
		if numChars < len(result) {
			result = result[numChars:]
		} else {
			result = ""
		}
	}

	if ignoreCase {
		result = strings.ToLower(result)
	}

	return result
}

func skipFields(line string, n int) string {

	fields := strings.Fields(line)

	if n >= len(fields) {
		return ""
	}

	return strings.Join(fields[n:], " ")
}
