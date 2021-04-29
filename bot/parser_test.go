package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseText(t *testing.T) {
	t.Parallel()

	result, err := ParseText("1 Mark 2:12")

	expected := &ParsedText{
		Book:    "1 Mark",
		Chapter: "2",
		Start:   "12",
		End:     "",
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	result, err = ParseText("Mark 2:12 - 13")

	expected = &ParsedText{
		Book:    "Mark",
		Chapter: "2",
		Start:   "12",
		End:     "13",
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, result)

}

func TestParseTextInvalid(t *testing.T) {
	t.Parallel()

	expected := "incorrect text provided"

	result, err := ParseText("1 Mark 212")
	assert.Nil(t, result)
	assert.Equal(t, expected, err.Error())
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	parsed := &ParsedText{
		Book:    "1 Mark",
		Chapter: "10",
		Start:   "11",
		End:     "13",
	}

	assert.Equal(t, true, parsed.IsValid())

	parsed = &ParsedText{
		Book:    "1 Mark",
		Chapter: "10",
		Start:   "11",
		End:     "",
	}

	assert.Equal(t, true, parsed.IsValid())
}

func TestIsValidNot(t *testing.T) {
	t.Parallel()

	t.Run("No Start", func(t *testing.T) {
		parsed := &ParsedText{
			Book:    "1 Mark",
			Chapter: "10",
			Start:   "",
			End:     "",
		}
		assert.Equal(t, false, parsed.IsValid())
	})

	t.Run("No Chapter", func(t *testing.T) {
		parsed := &ParsedText{
			Book:    "1 Mark",
			Chapter: "",
			Start:   "10",
			End:     "",
		}

		assert.Equal(t, false, parsed.IsValid())
	})

	t.Run("no_book", func(t *testing.T) {
		parsed := &ParsedText{
			Book:    "",
			Chapter: "10",
			Start:   "10",
			End:     "",
		}

		assert.Equal(t, false, parsed.IsValid())
	})

	t.Run("end<start", func(t *testing.T) {
		parsed := &ParsedText{
			Book:    "1 Mark",
			Chapter: "10",
			Start:   "10",
			End:     "9",
		}

		assert.Equal(t, false, parsed.IsValid())
	})

}
