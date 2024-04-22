package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	NumberOfBytes      = "c"
	NumberOfLines      = "l"
	NumberOfWords      = "w"
	NumberOfCharacters = "m"
	DefaultValue       = ""
	DefaultUsage       = "Input Valid File Path"
)

type FileInfo struct {
	IsExist         bool
	SizeInBytes     int64
	TotalLines      int
	TotalWords      int
	TotalCharacters int
}

func main() {
	var isInputFromPipe = isInputFromPipe()
	if isInputFromPipe {
		inputFromPipe, _ := readInputFromPipe()
		err := os.WriteFile("/tmp/temp.txt", []byte(inputFromPipe), 0644)
		if err != nil {
			exit()
		}

		os.Args = append(os.Args, "/tmp/temp.txt")
	}

	if len(os.Args) <= 1 {
		exit()
	}

	command := os.Args[1:2]
	if command == nil {
		exit()
	}

	chosenFlag := getFlag(command[0])
	filePath := flag.String(chosenFlag, DefaultValue, DefaultUsage)
	flag.Parse()
	if *filePath == "" {
		filePath = &chosenFlag
	}

	fileInfo := getFileInfo(*filePath)
	if !fileInfo.IsExist {
		fmt.Println("File is not exists")

		os.Exit(1)
	}

	var displayFilePath string
	if isInputFromPipe {
		displayFilePath = ""
	} else {
		displayFilePath = *filePath
	}

	switch chosenFlag {
	case NumberOfBytes:
		fmt.Println(fileInfo.SizeInBytes, displayFilePath)

	case NumberOfLines:
		fmt.Println(fileInfo.TotalLines, displayFilePath)

	case NumberOfWords:
		fmt.Println(fileInfo.TotalWords, displayFilePath)

	case NumberOfCharacters:
		fmt.Println(fileInfo.TotalCharacters, displayFilePath)

	default:
		fmt.Println(fileInfo.TotalLines, fileInfo.TotalWords, fileInfo.SizeInBytes, displayFilePath)
	}

}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func readInputFromPipe() (string, error) {
	input, err := io.ReadAll(os.Stdin)

	return string(input), err
}

func getFlag(command string) string {
	if len(command) <= 1 || !strings.Contains(command, "-") {
		return command
	}

	return command[len(command)-1:]
}

func getFileInfo(filePath string) FileInfo {
	data, err := os.Stat(filePath)
	if err != nil || errors.Is(err, os.ErrNotExist) {

		return FileInfo{}
	}

	file, _ := os.Open(filePath)
	totalLines, _ := CalculateTotalLines(file)
	fileContent, _ := os.ReadFile(filePath)
	totalWords := CalculateTotalWords(string(fileContent))
	totalCharacters := CalculateTotalCharacters(string(fileContent))

	return FileInfo{
		IsExist:         true,
		SizeInBytes:     data.Size(),
		TotalLines:      totalLines,
		TotalWords:      totalWords,
		TotalCharacters: totalCharacters,
	}
}

func exit() {
	fmt.Println(`Invalid Arguments.
	See manual: 
		-c: Calculate Number of Bytes
		-l: Calculate Number of Lines
		-w: Calculate Number of Words
		-m: Calculate Number of Characters
	`)
	os.Exit(1)
}

func CalculateTotalWords(fileContent string) int {

	return len(strings.Fields(fileContent))
}

func CalculateTotalCharacters(fileContent string) int {

	return utf8.RuneCountInString(fileContent)
}

func CalculateTotalLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
