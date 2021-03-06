package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnum(t *testing.T) {
	assert.Equal(t, 131, Enum("Go"))
	assert.Equal(t, 141, Enum("HTML"))
}

func TestOk(t *testing.T) {
	assert.Equal(t, 303, Enum(`Ren'Py`))
	assert.False(t, Ok(`Ren'Py`))
	assert.True(t, Ok(`Go`))
}
