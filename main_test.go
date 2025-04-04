package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mentai-mayo/cli-go"
)

type Arguments struct {
	name string `arg:"1"`
}

func TestParse(t *testing.T) {
	rawargs := []string{"bin", "mentai-mayo"}
	_, err := cli.Parse[Arguments](rawargs)
	assert.NoError(t, err)
}
