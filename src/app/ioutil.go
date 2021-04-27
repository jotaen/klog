package app

import (
	"io"
	"os"
)

func ReadFile(path string) (string, Error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", NewError(
			"Cannot read file",
			"Location: "+path,
			err,
		)
	}
	return string(contents), nil
}

func RemoveFile(path string) Error {
	err := os.Remove(path)
	if err != nil {
		return NewError(
			"Cannot remove file",
			"Location: "+path,
			err,
		)
	}
	return nil
}

func WriteToFile(path string, contents string) Error {
	err := os.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		return NewError(
			"Cannot write to file",
			"Location: "+path,
			err,
		)
	}
	return nil
}

func ReadStdin() (string, Error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", NewError(
			"Cannot read from Stdin",
			"Cannot open Stdin stream to check for input",
			err,
		)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", NewError(
			"Error while reading from Stdin",
			"An error occurred while processing the input stream",
			err,
		)
	}
	return string(bytes), nil
}
