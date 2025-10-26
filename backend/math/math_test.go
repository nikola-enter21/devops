package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	require.Equal(t, 5, result)
}
