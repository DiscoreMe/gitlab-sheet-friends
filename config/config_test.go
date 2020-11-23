package config_test

import (
	"os"
	"path"
	"testing"

	"github.com/DiscoreMe/gitlab-sheets-friends/config"
	"github.com/stretchr/testify/assert"
)

const testConfig = `
db: 'test.db'
sheet_id: '1'
`

func TestLoadConfig(t *testing.T) {
	filepath := path.Join(os.TempDir(), "config.test.yaml")
	file, err := os.Create(filepath)
	assert.NoError(t, err)
	_, err = file.WriteString(testConfig)
	assert.NoError(t, err)
	assert.NoError(t, file.Close())

	conf, err := config.LoadConfig(filepath)
	assert.NoError(t, err)

	assert.NoError(t, os.Remove(filepath))

	assert.Equal(t, "test.db", conf.DB)
	assert.Equal(t, "1", conf.SpreadSheetID)
}
