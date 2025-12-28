package kfl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokeniseEmptyToken(t *testing.T) {
	{ // Empty
		p, _, err := tokenise("")
		require.Nil(t, err)
		assert.Equal(t, p, []token{})
	}
	{ // Blank
		p, _, err := tokenise("    ")
		require.Nil(t, err)
		assert.Equal(t, p, []token{})
	}
}

func TestTokeniseAllTokens(t *testing.T) {
	p, _, err := tokenise("2020-01-01 && #hello || (2020-02-02 && !2021-Q4) && type:duration")
	require.Nil(t, err)
	assert.Equal(t, []token{
		tokenDate{"2020-01-01"},
		tokenAnd{},
		tokenTag{"hello"},
		tokenOr{},
		tokenOpenBracket{},
		tokenDate{"2020-02-02"},
		tokenAnd{},
		tokenNot{},
		tokenPeriod{"2021-Q4"},
		tokenCloseBracket{},
		tokenAnd{},
		tokenEntryType{"duration"},
	}, p)
}

func TestDisregardWhitespaceBetweenTokens(t *testing.T) {
	p, _, err := tokenise("   2020-01-01    &&     #hello    ||    (   2020-02-02   &&   !   2021-Q4  )  &&    type:duration")
	require.Nil(t, err)
	assert.Equal(t, []token{
		tokenDate{"2020-01-01"},
		tokenAnd{},
		tokenTag{"hello"},
		tokenOr{},
		tokenOpenBracket{},
		tokenDate{"2020-02-02"},
		tokenAnd{},
		tokenNot{},
		tokenPeriod{"2021-Q4"},
		tokenCloseBracket{},
		tokenAnd{},
		tokenEntryType{"duration"},
	}, p)
}

func TestFailsOnUnrecognisedToken(t *testing.T) {
	for _, txt := range []string{
		"abcde",
		"2020-01-01 & 2020-01-02",
		"2020-01-01 * 2020-01-02",
		"2020-01-01 {2020-01-02}",
	} {
		t.Run(txt, func(t *testing.T) {
			p, _, err := tokenise(txt)
			require.ErrorIs(t, err.Original(), ErrUnrecognisedToken)
			assert.Nil(t, p)
		})
	}
}

func TestFailsOnMissingWhitespace(t *testing.T) {
	for _, txt := range []string{
		"2021-12-12 &&&",
		"2021-12-12 &&&&",
		"2021-12-12 &&||",
		"2021-12-12 &&2021-12-12",
		"2021-12-12 &&#tag",
		"2021-12-12 &&(2021-12-12 || #foo)",

		"2021-12-12 |||",
		"2021-12-12 ||||",
		"2021-12-12 ||&&",
		"2021-12-12 ||2021-12-12",
		"2021-12-12 ||#tag",
		"2021-12-12 ||(2021-12-12 || #foo)",

		"(#foo)(#bar)",
		"( #foo )( #bar )",

		"2020-01-01&&",
		"2020-01-01||",
		"2020-01-01( #foo )",
		"2020-01-01#foo",

		"2020-01-01...2020-01-31&&",
		"2020-01-01...2020-01-31( #foo )",
		"2020-01-01...&&",
		"2020-01-01...( #foo )",

		"(2021-12-12 || #foo)2020-01-01",
		"(2021-12-12 || #foo)&& #foo",
		"(2021-12-12 || #foo)#foo",

		"#tag&& #tag",
		"#tag|| #tag",
		"#tag( 2020-01-01)",

		"2020-Q4&&",
		"2020-Q4||",
		"2020-Q4( 2020-01-01 )",
		"2020-Q4!( 2020-01-01 )",
	} {
		t.Run(txt, func(t *testing.T) {
			p, _, err := tokenise(txt)
			require.ErrorIs(t, err.Original(), ErrMissingWhiteSpace)
			assert.Nil(t, p)
		})
	}
}
