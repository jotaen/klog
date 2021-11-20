package lib

import (
	"errors"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/parser"
)

// ReconcilerChain is an automatism that reads input from a file, runs one or
// more reconcilers on it, and then writes the result back to the file.
type ReconcilerChain struct {
	File app.FileOrBookmarkName
	Ctx  app.Context
}

type NotEligibleError struct{}

func (e NotEligibleError) Error() string { return "No record found at that date" }

func (c ReconcilerChain) Apply(
	applicators ...func(pr *parser.ParseResult) (*parser.ReconcileResult, error),
) error {
	pr, targetFilePath, err := c.Ctx.ReadFileInput(c.File)
	if err != nil {
		return err
	}
	result, err := func() (*parser.ReconcileResult, error) {
		for i, a := range applicators {
			result, err := a(pr)
			if result != nil {
				return result, nil
			}
			_, isNotEligibleError := err.(NotEligibleError)
			if isNotEligibleError && i < len(applicators)-1 {
				// Try next reconcile function
				continue
			}
			return nil, err
		}
		return nil, errors.New("No applicable record found")
	}()
	if err != nil {
		return err
	}
	err = c.Ctx.WriteFile(targetFilePath, result.NewText)
	if err != nil {
		return err
	}
	c.Ctx.Print("\n" + c.Ctx.Serialiser().SerialiseRecords(result.NewRecord) + "\n")
	return nil
}
