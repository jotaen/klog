package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokeniseEmptyToken(t *testing.T) {
	{ // Empty
		p, err := tokenise("")
		require.Nil(t, err)
		assert.Equal(t, p, []token{})
	}
	{ // Blank
		p, err := tokenise("    ")
		require.Nil(t, err)
		assert.Equal(t, p, []token{})
	}
}

func TestTokeniseAllTokens(t *testing.T) {
	p, err := tokenise("2020-01-01 && #hello || (2020-02-02 && !2021-Q4) && type:duration")
	require.Nil(t, err)
	assert.Equal(t, []token{
		{tokenDate, "2020-01-01", 0},
		{tokenAnd, "&&", 11},
		{tokenTag, "#hello", 14},
		{tokenOr, "||", 21},
		{tokenOpenBracket, "(", 24},
		{tokenDate, "2020-02-02", 25},
		{tokenAnd, "&&", 36},
		{tokenNot, "!", 39},
		{tokenPeriod, "2021-Q4", 40},
		{tokenCloseBracket, ")", 47},
		{tokenAnd, "&&", 49},
		{tokenEntryType, "type:duration", 52},
	}, p)
}

func TestDisregardWhitespaceBetweenTokens(t *testing.T) {
	p, err := tokenise("   2020-01-01    &&     #hello    ||    (   2020-02-02   &&   !   2021-Q4  )  &&    type:duration")
	require.Nil(t, err)
	assert.Equal(t, []token{
		{tokenDate, "2020-01-01", 3},
		{tokenAnd, "&&", 17},
		{tokenTag, "#hello", 24},
		{tokenOr, "||", 34},
		{tokenOpenBracket, "(", 40},
		{tokenDate, "2020-02-02", 44},
		{tokenAnd, "&&", 57},
		{tokenNot, "!", 62},
		{tokenPeriod, "2021-Q4", 66},
		{tokenCloseBracket, ")", 75},
		{tokenAnd, "&&", 78},
		{tokenEntryType, "type:duration", 84},
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
			p, err := tokenise(txt)
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

		"type:duration&&",
		"type:duration||",
		"type:duration( 2020-01-01 )",
		"type:duration!( 2020-01-01 )",
	} {
		t.Run(txt, func(t *testing.T) {
			p, err := tokenise(txt)
			require.ErrorIs(t, err.Original(), ErrMissingWhiteSpace)
			assert.Nil(t, p)
		})
	}
}
