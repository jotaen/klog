package lib

import (
	"errors"
	klog "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/app"
	"github.com/jotaen/klog/src/parser/lineparsing"
	"github.com/jotaen/klog/src/parser/reconciler"
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
	applicators ...func(records []klog.Record, blocks []lineparsing.Block) (*reconciler.ReconcileResult, error),
) error {
	records, blocks, targetFilePath, err := c.Ctx.ReadFileInput(c.File)
	if err != nil {
		return err
	}
	result, err := func() (*reconciler.ReconcileResult, error) {
		for i, a := range applicators {
			result, err := a(records, blocks)
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
