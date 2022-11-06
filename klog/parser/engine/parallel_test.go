package engine

import (
	"github.com/jotaen/klog/klog/parser/txt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var identityParser = ParallelBatchParser[string]{
	SerialParser: SerialParser[string]{
		ParseOne: func(b txt.Block) (string, []txt.Error) {
			original := ""
			for _, l := range b.Lines() {
				original += l.Original()
			}
			return original, nil
		},
	},
	NumberOfWorkers: 100,
}

func TestParallelParserDoesNotMessUpBatchOrder(t *testing.T) {
	// The mock parser has 100 workers, so the batch size will be 1 char per worker.
	// The serial parser is basically an identity function, so it returns the input
	// text of the block, i.e. that one char per worker. The parallel parser is now
	// expected to re-construct the original order of the input after batching.
	// If it wouldn’t do that, the return text would be messed up, e.g. `7369285014`
	// instead of `1234567890`.
	val, _, _ := identityParser.Parse("1234567890")
	assert.Equal(t, []string{"1234567890"}, val)
}

func TestParallelParser(t *testing.T) {
	for _, x := range []struct {
		txt    string
		chunks int
		exp    []string
	}{
		// Small ASCII strings:
		{"Hello", 1, []string{"Hello"}},
		{"Hello", 2, []string{"Hel", "lo"}},
		{"Hello", 3, []string{"He", "ll", "o"}},
		{"Hello", 4, []string{"He", "ll", "o", ""}},
		{"Hello", 5, []string{"H", "e", "l", "l", "o"}},
		{"Hello", 6, []string{"H", "e", "l", "l", "o", ""}},
		{"Hello", 8, []string{"H", "e", "l", "l", "o", "", "", ""}},

		// Larger ASCII strings:
		{"abcdefghijklmnopqrstuvwxyz", 3, []string{"abcdefghi", "jklmnopqr", "stuvwxyz"}},
		{"abcdefghijklmnopqrstuvwxyz", 13, []string{"ab", "cd", "ef", "gh", "ij", "kl", "mn", "op", "qr", "st", "uv", "wx", "yz"}},

		// UTF-8 strings: (reminder: the chunks are supposed to have similar byte-size, not character-count!)
		{"藤本太郎喜左衛門将時能", 4, []string{"藤本太", "郎喜左", "衛門将", "時能"}},
		{"藤本太郎喜左衛門将時能", 11, []string{"藤", "本", "太", "郎", "喜", "左", "衛", "門", "将", "時", "能"}},
		{"藤😀abcdef©½, ★Test🤡äß©•¥üöπგამარჯობა", 3, []string{"藤😀abcdef©½, ★Tes", "t🤡äß©•¥üöπგ", "ამარჯობა"}},
	} {
		chunks := splitIntoChunks(x.txt, x.chunks)
		assert.Equal(t, x.exp, chunks)

		val, _, errs := identityParser.Parse(x.txt)
		assert.Nil(t, errs)
		assert.Equal(t, []string{x.txt}, val)
	}
}
