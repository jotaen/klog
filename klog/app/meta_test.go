package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMeta(t *testing.T) {
	m := NewMeta("Specification", "License", "v1.0", "1234abc", "2000-12-31")
	assert.Equal(t, "Specification", m.Specification)
	assert.Equal(t, "License", m.License)
	assert.Equal(t, "v1.0", m.Version)
	assert.Equal(t, "1234abc", m.SrcHash)
	assert.Equal(t, "2000-12-31", m.BuildDate)
}

func TestCreateMetaTrimsWhiteSpace(t *testing.T) {
	m := NewMeta("Specification", "License", "   v1.0   ", "  1234abc  ", " 2000-12-31 ")
	assert.Equal(t, "Specification", m.Specification)
	assert.Equal(t, "License", m.License)
	assert.Equal(t, "v1.0", m.Version)
	assert.Equal(t, "1234abc", m.SrcHash)
	assert.Equal(t, "2000-12-31", m.BuildDate)
}

func TestAcceptsValidVersion(t *testing.T) {
	for _, in := range []string{
		"v0.0",
		"v000.000",
		"v001.002",
		"v1.1",
		"v183.5",
		"v2.5910",
		"v28372.5910",
		"v5.23-alpha",
		"v5.23.42-alpha",
		"v5.23.42-alpha_SOMETHING&$^%!*@&#",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", in, "1234abc", "2000-12-31")
			assert.Equal(t, in, m.Version)
		})
	}
}

func TestRejectsInvalidVersion(t *testing.T) {
	for _, in := range []string{
		"",
		"foo",
		"0.0",
		"v1",
		"v1.",
		"v1.2.",
		"v1.2-",
		"v.12",
		"v1.22alpha",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", in, "1234abc", "2000-12-31")
			assert.Equal(t, "v?.?", m.Version)
		})
	}
}

func TestAcceptsValidSrcHash(t *testing.T) {
	for _, in := range []string{
		"1a2b3c4",
		"1A2B3C4",
		"1590acf",
		"1234567890abcdef",
		"1234567890ABCDEF",
		"7a4f7b6b0749adb164aa6ac67dee9ba6556f261e",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", "v1.2", in, "2000-12-31")
			assert.Equal(t, in[:7], m.SrcHash)
		})
	}
}

func TestRejectsInvalidSrcHash(t *testing.T) {
	for _, in := range []string{
		"",
		"0",
		"123456",
		"123456g",
		"abcdefghijklmnopqrst",
		"&^@!#@(*#&!*@&",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", "v1.2", in, "2000-12-31")
			assert.Equal(t, "???????", m.SrcHash)
		})
	}
}

func TestAcceptsValidBuildDate(t *testing.T) {
	for _, in := range []string{
		"2000-12-31",
		"3617-05-15",
		"9999-99-99",
		"0001-01-01",
		"2000-12-31T12:00:23Z",
		"2000-12-31Tanything",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", "v1.2", "1234567", in)
			assert.Equal(t, in[:10], m.BuildDate)
		})
	}
}

func TestRejectsInvalidBuildDate(t *testing.T) {
	for _, in := range []string{
		"",
		"153728182883",
		"2000_12_31",
		"2000 12 31",
		"2000/12/31",
		"abcd-ef-gh",
	} {
		t.Run(in, func(t *testing.T) {
			m := NewMeta("Specification", "License", "v1.2", "1234567", in)
			assert.Equal(t, "????-??-??", m.BuildDate)
		})
	}
}
