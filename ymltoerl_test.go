package ymltoerl_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"path/filepath"

	"github.com/makasim/ymltoerl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixtures(main *testing.T) {
	files, err := ioutil.ReadDir("_fixtures")
	require.NoError(main, err)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		match, err := filepath.Match("*.yaml", f.Name())
		require.NoError(main, err)
		if !match {
			continue
		}

		f := f
		main.Run(f.Name(), func(t *testing.T) {
			expectedFile := filepath.Join("_fixtures", strings.Replace(f.Name(), ".yaml", ".config", 1))
			expected, err := ioutil.ReadFile(expectedFile)
			require.NoError(t, err)

			actual, err := ymltoerl.ConvertFile(filepath.Join("_fixtures", f.Name()))
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(actual), string(actual))
		})
	}
}
