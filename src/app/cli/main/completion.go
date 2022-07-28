package klog

import (
	"github.com/jotaen/klog/src/app"
	"github.com/posener/complete"
)

func PredictBookmarks(ctx app.Context) complete.Predictor {
	thunk := func() []string {
		names := make([]string, 0)
		bookmarksCollection, err := ctx.ReadBookmarks()
		if err != nil {
			return names
		}
		for _, bookmark := range bookmarksCollection.All() {
			names = append(names, bookmark.Name().ValuePretty())
		}
		return names
	}
	return complete.PredictFunc(func(a complete.Args) []string { return thunk() })
}

func Predictors(ctx app.Context) map[string]complete.Predictor {
	return map[string]complete.Predictor{
		"file":             complete.PredictFiles("*.klg"),
		"bookmark":         PredictBookmarks(ctx),
		"file or bookmark": complete.PredictOr(complete.PredictFiles("*.klg"), PredictBookmarks(ctx)),
	}
}
