package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestParseHeaders_Extended(t *testing.T) {
	t.Run("ValidSpecialCharactersAndLowercasing", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("X-Test_Token!$: value\r\n\r\n")
		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.Equal(t, "value", headers["x-test_token!$"])
		assert.Equal(t, 25, n)
		assert.True(t, done)
	})

	t.Run("InvalidCharacterInName", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost:42069\r\n\r\n")
		_, _, err := headers.Parse(data)
		require.Error(t, err)
	})

	t.Run("SpaceBeforeColonFails", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Content-Type : application/json\r\n\r\n")
		_, _, err := headers.Parse(data)
		require.Error(t, err)
	})

	t.Run("EmptyValue", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Empty-Header:\r\n\r\n")
		_, done, err := headers.Parse(data)
		require.NoError(t, err)
		val, ok := headers["empty-header"]
		assert.True(t, ok)
		assert.Equal(t, "", val)
		assert.True(t, done)
	})
	t.Run("DuplicateHeaderKeys", func(t *testing.T) {
		headers := NewHeaders()
		// Multiple instances of the same key
		data := []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")

		n, done, err := headers.Parse(data)
		require.NoError(t, err)
		assert.True(t, done)

		// Expected result: all values joined by a comma
		expected := "lane-loves-go, prime-loves-zig, tj-loves-ocaml"
		assert.Equal(t, expected, headers["set-person"])
		assert.Equal(t, len(data), n)
	})
}
