package cli

import (
	"klog"
	"klog/service"
	"time"
)

type MultipleFilesArgs struct {
	File []string `arg required type:"existingfile" name:"file" help:".klg source file(s)"`
}

type SingleFileArgs struct {
	File string `arg required type:"existingfile" name:"file" help:".klg source file"`
}

type FilterArgs struct {
	Tags      []string    `name:"tag" help:"Only records that contain this tag"`
	Date      []klog.Date `name:"date" help:"Only records at this date"`
	Today     bool        `name:"today" help:"Shorthand for today’s date"`
	Yesterday bool        `name:"yesterday" help:"Shorthand for yesterday’s date"`
	AfterEq   klog.Date   `name:"after" help:"Only records at or after this date"`
	BeforeEq  klog.Date   `name:"before" help:"Only records at or before this date"`
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
