package ops

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)
	WriteTestDB(t, dataDir)
	err = DoStats(dataDir)
	assert.NoError(t, err)
}

func TestRecursiveStatistics(t *testing.T) {
	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	err = os.MkdirAll(baseDataDir, os.ModePerm)
	defer os.RemoveAll(baseDataDir)
	assert.NoError(t, err)

	paths := []string{
		"1/store_1/",
		"1/store_2/",
		"2/store_1/",
		"2/store_2/",
	}

	for _, relLocation := range paths {
		WriteTestDB(t, filepath.Join(baseDataDir, relLocation))
	}

	err = DoRecursiveStats(baseDataDir, 1)
	assert.NoError(t, err)
}
