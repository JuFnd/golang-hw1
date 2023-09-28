package uniq

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

func parseFlags() (UniqFlags, error) {
	flags := UniqFlags{}

	flag.BoolVar(&flags.count, "c", false, "Count the number of occurrences of each line")
	flag.BoolVar(&flags.duplicates, "d", false, "Print only duplicate lines")
	flag.BoolVar(&flags.unique, "u", false, "Print only unique lines")
	flag.BoolVar(&flags.ignoreCase, "i", false, "Ignore case when comparing lines")
	flag.IntVar(&flags.numFields, "f", 0, "Ignore the first num_fields fields in each line")
	flag.IntVar(&flags.numChars, "s", 0, "Ignore the first num_chars characters in each line")

	flag.Parse()

	if flags.duplicates && flags.unique {
		return flags, fmt.Errorf("")
	}

	if flags.numFields < 0 || flags.numChars < 0 {
		return flags, fmt.Errorf("numeric values must be a non-negative integer")
	}

	return flags, nil
}

type Scan struct {
	reader *bufio.Scanner
	file   *os.File
}

type Write struct {
	writer *bufio.Writer
	file   *os.File
}

var writer Write
var scanner Scan

func createScanner(inputFilename *string) error {
	scanner = Scan{}
	if *inputFilename != "" {
		file, err := openFile(inputFilename)
		if err != nil {
			return err
		}
		scanner.file = file
		scanner.reader = bufio.NewScanner(scanner.file)
	} else {
		scanner.reader = bufio.NewScanner(os.Stdin)
	}
	return nil
}

func createWriter(outputFilename *string) error {
	writer = Write{}
	if *outputFilename != "" {
		file, err := openFile(outputFilename)
		if err != nil {
			return err
		}
		writer.file = file
		writer.writer = bufio.NewWriter(writer.file)
	} else {
		writer.writer = bufio.NewWriter(os.Stdout)
	}
	return nil
}

func openFile(filename *string) (*os.File, error) {
	file, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}
	return file, nil
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

func isSymCounter(isCounter bool, prevLine string, counter int) {
	if isCounter {
		writer.writer.WriteString(fmt.Sprintf("%d %s\n", counter+1, prevLine))
	} else if !isCounter && prevLine != "" {
		writer.writer.WriteString(fmt.Sprintf("%s\n", prevLine))
	}
}

func parseFile(inputFilename *string, outputFilename *string, flags *UniqFlags) error {
	err := createScanner(inputFilename)
	if err != nil {
		return err
	}
	err = createWriter(outputFilename)
	if err != nil {
		return err
	}

	counter := 0
	prevLine := ""
	for scanner.reader.Scan() {
		line := scanner.reader.Text()

		if isRegIgnore(flags.ignoreCase, line, prevLine, flags.numFields, flags.numChars) {
			counter++
		} else {
			if (flags.unique && counter == 0 || flags.duplicates && counter > 0 || !flags.duplicates && !flags.unique) && prevLine != "" {
				isSymCounter(flags.count, prevLine, counter)
			}

			counter = 0
		}
		prevLine = line
	}

	if (flags.unique && counter == 0) || (flags.duplicates && counter > 0) || (!flags.duplicates && !flags.unique) {
		isSymCounter(flags.count, prevLine, counter)
	}

	writer.writer.Flush()

	return nil
}

func Uniq() {
	flags, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)
	err = parseFile(&inputFilename, &outputFilename, &flags)
	writer.writer.Flush()
	defer scanner.file.Close()
	defer writer.file.Close()
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}
}
