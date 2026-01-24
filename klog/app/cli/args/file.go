package args

import "github.com/jotaen/klog/klog/app"

type InputFilesArgs struct {
	File []app.FileOrBookmarkName `arg:"" optional:"" type:"string" completion-predictor:"file_or_bookmark" name:"file or bookmark" help:"One or more .klg source files or bookmarks. If absent, klog tries to use the default bookmark."`
}

type OutputFileArgs struct {
	File app.FileOrBookmarkName `arg:"" optional:"" type:"string" completion-predictor:"file_or_bookmark" name:"file or bookmark" help:"One .klg source file or bookmark. If absent, klog tries to use the default bookmark."`
}
