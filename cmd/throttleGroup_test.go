package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChanSync_New(t *testing.T) {
	c := NewChanSync(100)
	assert.NotNil(t, c)
	assert.Equal(t, 100, c.Size, "Size of channel should match")
}

func TestChanSync_Done(t *testing.T) {
}

func TestChanSync_Close(t *testing.T) {
	c := NewChanSync(100)
	err := c.Inc()
	assert.NoError(t, err)
	l, _ := c.Length()
	assert.Equal(t, l, 1)
	err = c.Inc()
	l, _ = c.Length()
	assert.Equal(t, l, 2)
	c.Close()
	assert.Equal(t, c.Closed(), true)
	err = c.Dec()
	assert.NoError(t, err)
	err = c.Dec()
	assert.NoError(t, err)
	err = c.Dec()
	assert.Error(t, err)
}

func Test_Loop(t *testing.T) {
	c := NewChanSync(5)
	for i := 0; i < 10; i++ {
		c.Inc()
	}
}
