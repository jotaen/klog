package app

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

func readFile(path string) (string, Error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", appError{
			"Cannot read file",
			"Location: " + path,
		}
	}
	return string(contents), nil
}

func removeFile(path string) Error {
	err := os.Remove(path)
	if err != nil {
		return appError{
			"Cannot remove file",
			"Location: " + path,
		}
	}
	return nil
}

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

func readStdin() (string, Error) {
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
