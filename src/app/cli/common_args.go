package cli

import (
	"klog"
	"klog/service"
)

type FilesArgs struct {
	File []string `arg type:"existingfile" name:"file" help:".klg source file(s)"`
}

type FilterArgs struct {
	Tags     []string `short:"t" name:"tag" help:"Only records that contain this tag"`
	Date     src.Date `short:"d" name:"date" help:"Only records at this date"`
	AfterEq  src.Date `short:"a" name:"after" help:"Only records at or after this date"`
	BeforeEq src.Date `short:"b" name:"before" help:"Only records at or before this date"`
}

func (args *FilterArgs) toFilter() service.Filter {
	filter := service.Filter{
		BeforeEq: args.BeforeEq,
		AfterEq:  args.AfterEq,
		Tags:     args.Tags,
	}
	if args.Date != nil {
		filter.BeforeEq = args.Date
		filter.AfterEq = args.Date
	}
	return filter
}
