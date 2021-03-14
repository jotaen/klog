package app

import (
	"bufio"
	"io"
	"os"
)

func ReadFile(path string) (string, Error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", appError{
			"Cannot read file",
			"Location: " + path,
		}
	}
	return string(contents), nil
}

func RemoveFile(path string) Error {
	err := os.Remove(path)
	if err != nil {
		return appError{
			"Cannot remove file",
			"Location: " + path,
		}
	}
	return nil
}

// Deprecated
func appendToFile(path string, textToAppend string) Error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return appError{
			"Cannot write to file",
			"Location: " + path,
		}
	}
	defer file.Close()
	if _, err := file.WriteString(textToAppend); err != nil {
		return appError{
			"Cannot write to file",
			"Location: " + path,
		}
	}
	return nil
}

func WriteToFile(path string, contents string) Error {
	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		return appError{
			message: "Cannot write to file",
			help:    "Location: " + path,
		}
	}
	return nil
}

func ReadStdin() (string, Error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", appError{
			"Cannot read from Stdin",
			"Cannot open Stdin stream to check for input",
		}
	}
	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return "", nil
	}
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", appError{
				"Error while reading from Stdin",
				"An error occurred while processing the input stream",
			}
		}
		output = append(output, input)
	}
	return string(output), nil
}
