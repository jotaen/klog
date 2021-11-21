package reconciler

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReconcileAddBlockIfOriginalIsEmpty(t *testing.T) {
	rs, bs, _ := parser.Parse("")
	reconciler := NewRecordReconciler(rs, bs, Ɀ_Date_(3333, 1, 1))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2000-05-05", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, "2000-05-05\n", result.NewText)
}

func TestReconcileAddBlockToEnd(t *testing.T) {
	original := `
2018-01-01
    1h`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, bs, Ɀ_Date_(2018, 1, 2))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2018-01-02", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
`, result.NewText)
}

func TestReconcileAddBlockToEndWithTrailingNewlines(t *testing.T) {
	original := `
2018-01-01
    1h

`
	rs, bs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, bs, Ɀ_Date_(2018, 1, 2))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2018-01-02", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02

`, result.NewText)
}

func TestReconcileAddBlockToBeginning(t *testing.T) {
	original := "2018-01-02"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, bs, Ɀ_Date_(2018, 1, 1))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2018-01-01", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `2018-01-01

2018-01-02`, result.NewText)
}

func TestReconcileAddBlockToBeginningWithLeadingNewlines(t *testing.T) {
	original := "\n\n2018-01-02"
	rs, bs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(rs, bs, Ɀ_Date_(2018, 1, 1))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2018-01-01", 0},
	})
	require.Nil(t, err)
	assert.Equal(t, `2018-01-01



2018-01-02`, result.NewText)
}

func TestReconcileAddBlockInBetween(t *testing.T) {
	original := `
2018-01-01
    1h

2018-01-03
    3h`
	pr, bs, _ := parser.Parse(original)
	reconciler := NewRecordReconciler(pr, bs, Ɀ_Date_(2018, 1, 2))
	result, err := reconciler.InsertBlock([]InsertableText{
		{"2018-01-02", 0},
		{"This and that", 0},
		{"30m worked", 1},
	})
	require.Nil(t, err)
	assert.Equal(t, `
2018-01-01
    1h

2018-01-02
This and that
    30m worked

2018-01-03
    3h`, result.NewText)
}
