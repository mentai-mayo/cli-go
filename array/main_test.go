package array_test

import (
	"testing"

	"github.com/mentai-mayo/cli-go/array"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	arr := array.New[uint](4)
	arr.Push(2)
	arr.Push(2)
	arr.Push(2)
	arr.Push(2)
	assert.Equal(t, arr.Len(), 4)
}
