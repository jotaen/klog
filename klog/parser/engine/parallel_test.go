package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIntoChunks(t *testing.T) {
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
	}
}
