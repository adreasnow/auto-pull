package puller

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	t.Parallel()

	t.Run("https", func(t *testing.T) {
		t.Parallel()

		changes, err := Pull("/users/adrea.snow/scratch/github-actions")
		require.NoError(t, err)
		fmt.Println(changes)
	})

	t.Run("ssh", func(t *testing.T) {
		t.Parallel()

		changes, err := Pull("/users/adrea.snow/scratch/AdventOfCode/")
		require.NoError(t, err)
		fmt.Println(changes)
	})
}
