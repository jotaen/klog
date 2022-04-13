package klog

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatesNewTag(t *testing.T) {
	for _, x := range []struct {
		tag        string
		expectName string
	}{
		{"#tag", "tag"},
		{"#TAG", "tag"},
		{"#t-a-g", "t-a-g"},
		{"#t_a_g", "t_a_g"},
		{"#t1a2g3", "t1a2g3"},
		{"#---", "---"},
		{"#___", "___"},
	} {
		tag, err := NewTagFromString(x.tag)
		require.Nil(t, err)
		assert.Equal(t, x.expectName, tag.Name())
	}
}

func TestTagMatching(t *testing.T) {
	for _, x := range []struct {
		tag1 string
		tag2 string
	}{
		// Identity
		{`#tag`, `#tag`},
		{`#tag=value`, `#tag=value`},

		// Value empty is the same as value absent
		{`#tag`, `#tag=`},
		{`#tag`, `#tag=""`},
		{`#tag`, `#tag=''`},
		{`#tag=`, `#tag`},
		{`#tag=""`, `#tag`},
		{`#tag=''`, `#tag`},
		{`#tag=''`, `#tag=""`},
		{`#tag=''`, `#tag=`},

		// Case-insensitivity of name
		{`#TAG`, `#tag`},
		{`#TAG=value`, `#tag=value`},

		// Quotation style is irrelevant
		{`#tag=value`, `#tag="value"`},
		{`#tag=value`, `#tag='value'`},
		{`#tag="value"`, `#tag='value'`},
	} {
		first, err1 := NewTagFromString(x.tag1)
		require.Nil(t, err1)
		second, err2 := NewTagFromString(x.tag2)
		require.Nil(t, err2)
		assert.Equal(t, first, second)
	}
}

func TestTagIsNotMatching(t *testing.T) {
	for _, x := range []struct {
		tag1 string
		tag2 string
	}{
		// Name is different
		{`#tag`, `#t-a-g`},
		{`#tag`, `#tags`},

		// Query has value, but base hasn’t
		{`#tag`, `#tag=value`},

		// Query value is different than base’s
		{`#tag=value`, `#tag=VALUE`},
		{`#tag=value`, `#tag=foo`},
		{`#tag='V A L U E'`, `#tag='v a l u e'`},
		{`#tag=''`, `#tag=' '`},
	} {
		first, err1 := NewTagFromString(x.tag1)
		require.Nil(t, err1)
		second, err2 := NewTagFromString(x.tag2)
		require.Nil(t, err2)
		assert.NotEqual(t, first, second)
	}
}

func TestPrecedingHashCharIsOptional(t *testing.T) {
	tag, err := NewTagFromString("tag")
	require.Nil(t, err)
	assert.Equal(t, "tag", tag.Name())
}

func TestRejectsInvalidTags(t *testing.T) {
	for _, name := range []string{
		"",
		"##tag",
		"##tag",
		"a#tag",
		"a #tag",
		"#tag#tag",
		"#tag #tag",
		"#t^a*g",
		"#tag?",
		"#tag:tag",
		"#tag!!!",
		"#t-a?g",
		`#tag=foo=bar`,
		`#tag='foo`,
		`#tag='It's great'`,
		`#tag="foo`,
		`#tag="foo`,
		`#tag="`,
	} {
		_, err := NewTagFromString(name)
		require.Error(t, err)
	}
}

func TestCreatesNewTagWithValue(t *testing.T) {
	for _, x := range []struct {
		tag         string
		expectValue string
	}{
		{`#tag=value`, `value`},
		{`#tag=VALUE`, `VALUE`},
		{`#tag=V_A_L_U_E`, `V_A_L_U_E`},
		{`#tag=v-a-l-u-e`, `v-a-l-u-e`},
		{`#tag=v-a-l-u-e`, `v-a-l-u-e`},
		{`#tag="v a l u e"`, `v a l u e`},
		{`#tag='v!a?l,u=e'`, `v!a?l,u=e`},
		{`#tag='foo=bar'`, `foo=bar`},
	} {
		tag, err := NewTagFromString(x.tag)
		require.Nil(t, err)
		assert.Equal(t, x.expectValue, tag.Value())
	}
}

func TestSerialiseTag(t *testing.T) {
	tagWithoutValue := NewTagOrPanic("test", "")
	assert.Equal(t, "#test", tagWithoutValue.ToString())

	tagWithValue := NewTagOrPanic("test", "value")
	assert.Equal(t, "#test=value", tagWithValue.ToString())

	tagWithValueNeedingQuoting := NewTagOrPanic("test", "v a l u e")
	assert.Equal(t, `#test="v a l u e"`, tagWithValueNeedingQuoting.ToString())

	tagWithValueContainingDoubleQuote := NewTagOrPanic("test", `It's great`)
	assert.Equal(t, `#test="It's great"`, tagWithValueContainingDoubleQuote.ToString())

	tagWithValueContainingSingleQuote := NewTagOrPanic("test", `5"`)
	assert.Equal(t, `#test='5"'`, tagWithValueContainingSingleQuote.ToString())
}
