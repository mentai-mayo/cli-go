package cli_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mentai-mayo/cli-go"
)

type Arguments struct {
	Name string `pos:"1"`
}

type ErrArguments struct {
	Name Arguments
}

func TestParse(t *testing.T) {
	rawargs := []string{"bin", "mentai-mayo"}
	args, err := cli.Parse[Arguments](rawargs)
	require.NoError(t, err)
	fmt.Printf("%#v", args)
	assert.Equal(t, args.Name, "mentai-mayo")
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
