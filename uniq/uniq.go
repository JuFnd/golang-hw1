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
	Count      bool
	Duplicates bool
	Unique     bool
	IgnoreCase bool
	NumFields  int
	NumChars   int
}

func parseFlags() (UniqFlags, error) {
	flags := UniqFlags{}

	flag.BoolVar(&flags.Count, "c", false, "Count the number of occurrences of each line")
	flag.BoolVar(&flags.Duplicates, "d", false, "Print only duplicate lines")
	flag.BoolVar(&flags.Unique, "u", false, "Print only unique lines")
	flag.BoolVar(&flags.IgnoreCase, "i", false, "Ignore case when comparing lines")
	flag.IntVar(&flags.NumFields, "f", 0, "Ignore the first num_fields fields in each line")
	flag.IntVar(&flags.NumChars, "s", 0, "Ignore the first num_chars characters in each line")

	flag.Parse()

	if flags.Duplicates && flags.Unique {
		return flags, fmt.Errorf("")
	}

	if flags.NumFields < 0 || flags.NumChars < 0 {
		return flags, fmt.Errorf("numeric values must be a non-negative integer")
	}

	return flags, nil
}

func createScanner(inputFilename string) (*os.File, error) {
	if inputFilename == "" {
		return os.Stdin, nil
	} else {
		file, err := openFile(inputFilename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func createWriter(outputFilename string) (*os.File, error) {
	if outputFilename == "" {
		return os.Stdout, nil
	} else {
		file, err := openFile(outputFilename)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func openFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
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
		if prevLine != " " {
			writer.Write([]byte(fmt.Sprintf("%d %s\n", counter+1, prevLine)))
		} else {
			writer.Write([]byte(fmt.Sprintf("%d\n", counter+1)))
		}
	} else if !isCounter && prevLine != "" {
		writer.Write([]byte(fmt.Sprintf("%s\n", prevLine)))
	}
}

func Uniq(scanner io.Reader, writer io.Writer, flags *UniqFlags) error {
	counter := 0
	prevLine := ""
	sc := bufio.NewScanner(scanner)
	for sc.Scan() {
		line := sc.Text()

		if isRegIgnore(flags.IgnoreCase, line, prevLine, flags.NumFields, flags.NumChars) {
			counter++
		} else {
			if (flags.Unique && counter == 0 || flags.Duplicates && counter > 0 || !flags.Duplicates && !flags.Unique) && prevLine != "" {
				isSymCounter(flags.Count, prevLine, counter, writer)
			}

			counter = 0
		}
		prevLine = line
	}

	if (flags.Unique && counter == 0) || (flags.Duplicates && counter > 0) || (!flags.Duplicates && !flags.Unique) {
		isSymCounter(flags.Count, prevLine, counter, writer)
	}

	return nil
}

func ParseFile() error {
	flags, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		return err
	}

	inputFilename := flag.Arg(0)
	outputFilename := flag.Arg(1)

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

	err = Uniq(scanner, writer, &flags)
	if err != nil {
		return err
	}

	return nil
}
