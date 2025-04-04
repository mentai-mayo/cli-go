package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mentai-mayo/cli-go"
)

type Arguments struct {
	name string `pos:"1"`
}

type ErrArguments struct {
	name Arguments
}

func TestParse(t *testing.T) {
	rawargs := []string{"bin", "mentai-mayo"}
	_, err := cli.Parse[Arguments](rawargs)
	require.NoError(t, err)
}

func TestParse_ErrNonStructTarget(t *testing.T) {
	rawargs := []string{"bin", "mentai-mayo"}
	_, err := cli.Parse[string](rawargs)
	_, ok := err.(cli.NonStructTargetErr)
	assert.True(t, ok, err)
}

func TestParse_ErrInvliadExpectedType(t *testing.T) {
	rawargs := []string{"bin", "mentai-mayo"}
	_, err := cli.Parse[ErrArguments](rawargs)
	_, ok := err.(*cli.NonStructTargetErr)
	assert.True(t, ok, err)
}
