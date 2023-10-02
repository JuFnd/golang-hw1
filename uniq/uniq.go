package uniq

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

func createScanner(inputFilename *string) (*os.File, error) {
	if *inputFilename == "" {
		return os.Stdin, nil
	} else {
		file, err := openFile(inputFilename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func createWriter(outputFilename *string) (*os.File, error) {
	if *outputFilename == "" {
		return os.Stdout, nil
	} else {
		file, err := openFile(outputFilename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
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

func isSymCounter(isCounter bool, prevLine string, counter int, writer io.Writer) {
	if isCounter {
		writer.Write([]byte(fmt.Sprintf("%d %s\n", counter+1, prevLine)))
	} else if !isCounter && prevLine != "" {
		writer.Write([]byte(fmt.Sprintf("%s\n", prevLine)))
	}
}

func parseFile(inputFilename *string, outputFilename *string, flags *UniqFlags) error {
	inputFile, err := createScanner(inputFilename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := createWriter(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var scanner io.Reader = inputFile
	var writer io.Writer = outputFile

	counter := 0
	prevLine := ""
	sc := bufio.NewScanner(scanner)
	for sc.Scan() {
		line := sc.Text()

		if isRegIgnore(flags.ignoreCase, line, prevLine, flags.numFields, flags.numChars) {
			counter++
		} else {
			if (flags.unique && counter == 0 || flags.duplicates && counter > 0 || !flags.duplicates && !flags.unique) && prevLine != "" {
				isSymCounter(flags.count, prevLine, counter, writer)
			}

			counter = 0
		}
		prevLine = line
	}

	if (flags.unique && counter == 0) || (flags.duplicates && counter > 0) || (!flags.duplicates && !flags.unique) {
		isSymCounter(flags.count, prevLine, counter, writer)
	}

	if w, ok := writer.(*bufio.Writer); ok {
		w.Flush()
	}

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
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}
}
