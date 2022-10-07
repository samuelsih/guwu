package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSpamEmail(t *testing.T) {
	email := "test@mail7.io"

	result, err := IsSpamEmail(context.Background(), email)

	assert.NoError(t, err)
	assert.Equal(t, result, false)
}
