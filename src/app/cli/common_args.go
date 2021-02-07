package cli

import (
	"klog"
	"klog/service"
	"time"
)

type InputFilesArgs struct {
	File []string `arg optional type:"existingfile" name:"file" help:".klg source file(s) (if empty the bookmark is used)"`
}

type DiffArg struct {
	Diff bool `name:"diff" help:"Show difference between actual and should total time"`
}

type FilterArgs struct {
	Tags      []string    `name:"tag" help:"Only records (or particular entries) that match this tag"`
	Date      []klog.Date `name:"date" help:"Only records at this date"`
	Today     bool        `name:"today" help:"Shorthand for today’s date"`
	Yesterday bool        `name:"yesterday" help:"Shorthand for yesterday’s date"`
	AfterEq   klog.Date   `name:"after" help:"Only records after this date (inclusive)"`
	BeforeEq  klog.Date   `name:"before" help:"Only records before this date (inclusive)"`
}

func (args *FilterArgs) toFilter() service.Filter {
	filter := service.Filter{
		BeforeEq: args.BeforeEq,
		AfterEq:  args.AfterEq,
		Tags:     args.Tags,
		Dates:    args.Date,
	}
	if args.Today {
		filter.Dates = append(filter.Dates, klog.NewDateFromTime(time.Now()))
	}
	if args.Yesterday {
		filter.Dates = append(filter.Dates, klog.NewDateFromTime(time.Now().AddDate(0, 0, -1)))
	}
	return filter
}
