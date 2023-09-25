package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type UniqFlags struct {
	count      bool
	duplicates bool
	unique     bool
	ignoreCase bool
	numFields  int
	numChars   int
}

func parseFlags() UniqFlags {
	flags := UniqFlags{}

	flag.BoolVar(&flags.count, "c", false, "Count the number of occurrences of each line")
	flag.BoolVar(&flags.duplicates, "d", false, "Print only duplicate lines")
	flag.BoolVar(&flags.unique, "u", false, "Print only unique lines")
	flag.BoolVar(&flags.ignoreCase, "i", false, "Ignore case when comparing lines")
	flag.IntVar(&flags.numFields, "f", 0, "Ignore the first num_fields fields in each line")
	flag.IntVar(&flags.numChars, "s", 0, "Ignore the first num_chars characters in each line")

	flag.Parse()

	return flags
}

type Scan struct {
	reader *bufio.Scanner
	file   *os.File
	err    error
}

func createScanner(inputFilename *string) Scan {
	scanner := Scan{}
	if *inputFilename != "" {
		scanner.file, scanner.err = openFile(inputFilename)
		scanner.reader = bufio.NewScanner(scanner.file)
	} else {
		scanner.reader = bufio.NewScanner(os.Stdin)
	}
	return scanner
}

type Write struct {
	writer *bufio.Writer
	file   *os.File
	err    error
}

func createWriter(outputFilename *string) Write {
	writer := Write{}
	if *outputFilename != "" {
		writer.file, writer.err = openFile(outputFilename)
		writer.writer = bufio.NewWriter(writer.file)
	} else {
		writer.writer = bufio.NewWriter(os.Stdout)
	}
	return writer
}

func truncateString(line string, numChars int) string {
	if len(line) < numChars {
		return ""
	}

	if numChars == 0 {
		return line
	}
	truncated := line[numChars:]
	return truncated
}

func truncateWords(words []string, numFields int, numChars int) string {
	if len(words) < numFields {
		return ""
	}
	truncatedWords := words[numFields:]
	cuttedString := strings.Join(truncatedWords, " ")
	return truncateString(cuttedString, numChars)
}

func skipFields(line string, prevLine string, numFields int, numChars int) bool {
	splitLine := strings.Fields(line)
	splitPrevLine := strings.Fields(prevLine)
	compareLine := truncateWords(splitLine, numFields, numChars)
	comparePrevLine := truncateWords(splitPrevLine, numFields, numChars)
	return compareLine == comparePrevLine
}

func isRegIgnore(isIgnore bool, line string, prevLine string, numFields int, numChars int) bool {
	if isIgnore {
		regIgnoreLine := strings.ToLower(line)
		regIgnorePrevLine := strings.ToLower(prevLine)
		return skipFields(regIgnoreLine, regIgnorePrevLine, numFields, numChars)
	} else {
		return skipFields(line, prevLine, numFields, numChars)
	}
}

func isSymCounter(isCounter bool, prevLine string, counter int, writer *Write) {
	if isCounter {
		writer.writer.WriteString(fmt.Sprintf("%d %s\n", counter+1, prevLine))
	} else if !isCounter && prevLine != "" {
		writer.writer.WriteString(fmt.Sprintf("%s\n", prevLine))
	}
}

func openFile(filename *string) (*os.File, error) {
	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func parseFile(inputFilename *string, outputFilename *string, flags *UniqFlags) {
	scanner := createScanner(inputFilename)
	writer := createWriter(outputFilename)
	counter := 0
	prevLine := ""
	for scanner.reader.Scan() {
		line := scanner.reader.Text()
		if *&flags.duplicates && *&flags.unique {
			break
		}
		if isRegIgnore(*&flags.ignoreCase, line, prevLine, *&flags.numFields, *&flags.numChars) {
			counter++
		} else {
			if *&flags.unique && counter == 0 && prevLine != "" {
				isSymCounter(*&flags.count, prevLine, counter, &writer)
			}

			if *&flags.duplicates && counter > 0 && prevLine != "" {
				isSymCounter(*&flags.count, prevLine, counter, &writer)
			}

			if !*&flags.duplicates && !*&flags.unique && prevLine != "" {
				isSymCounter(*&flags.count, prevLine, counter, &writer)
			}
			counter = 0
		}
		prevLine = line
	}
	defer scanner.file.Close()

	if *&flags.unique && counter == 0 {
		isSymCounter(*&flags.count, prevLine, counter, &writer)
	}

	if *&flags.duplicates && counter > 0 {
		isSymCounter(*&flags.count, prevLine, counter, &writer)
	}

	if !*&flags.duplicates && !*&flags.unique {
		isSymCounter(*&flags.count, prevLine, counter, &writer)
	}
	writer.writer.Flush()
	defer writer.file.Close()
}

func main() {
	flags := parseFlags()
	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)

	parseFile(&inputFilename, &outputFilename, &flags)
}
