package cli

import (
	"klog/record"
	"klog/service"
)

type FileArgs struct {
	File string `arg optional name:"file" help:"File to read from"`
}

type FilterArgs struct {
	Tags     []string    `short:"t" long:"tag" help:"Only records that contain this tag"`
	Date     record.Date `short:"d" long:"date" help:"Only records at this date"`
	AfterEq  record.Date `short:"a" long:"after" help:"Only records at or after this date"`
	BeforeEq record.Date `short:"b" long:"before" help:"Only records at or before this date"`
}

func (args *FilterArgs) ToFilter() service.Filter {
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
