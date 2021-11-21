package reconciling

import (
	"github.com/jotaen/klog/src/parser/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultStylePrefs(t *testing.T) {
	assert.Equal(t, stylePreferences{
		indentationStyle: "    ",
		lineEndingStyle:  "\n",
	}, stylePreferencesOrDefault(nil))
}

func TestInsertRespectsExplicitStylePrefs(t *testing.T) {
	result := insert(
		[]engine.Line{
			engine.NewLineFromString("Hello\r\n", 1),
			engine.NewLineFromString("World!\r\n", 2),
			engine.NewLineFromString("How are you?\r\n", 3),
			engine.NewLineFromString("Bye.\r\n", 4),
		},
		3,
		[]InsertableText{
			{"I’m great.", 0},
			{"(I hope you too.)", 1},
		},
		stylePreferences{"  ", "\r\n"},
	)
	assert.Equal(t, "Hello\r\nWorld!\r\nHow are you?\r\nI’m great.\r\n  (I hope you too.)\r\nBye.\r\n", join(result))
}

func TestDetectsStylePreferencesFromOtherBlock(t *testing.T) {
	block := []engine.Line{
		engine.NewLineFromString("Hello\r\n", 1),
		engine.NewLineFromString("   World!\r\n", 2),
	}
	result := insert(
		block,
		2,
		[]InsertableText{
			{"Hi.", 1},
		},
		stylePreferencesOrDefault(block),
	)
	assert.Equal(t, "Hello\r\n   World!\r\n   Hi.\r\n", join(result))
}
