package cmd

import (
	"gopkg.in/errgo.v2/errors"
)

type ChanSync struct {
	chSync chan int
	Size   int
	open   bool
}

const tokenValue int = 1

func (c *ChanSync) Closed() bool {
	return !c.open
}

func (c *ChanSync) Dec() error {
	// remove a token from the channel to mark a goroutine finished
	if len(c.chSync) < 1 {
		return errors.New("empty channel, cannot report Done")
	}
	_, ok := <-c.chSync
	if !ok {
		c.Close()
	}
	return nil
}

func (c *ChanSync) Inc() error {
	// insert a token into the channel to mark a goroutine started
	if c.chSync == nil {
		return errors.New("nil struct")
	}
	if c.Closed() {
		return errors.New("cannot write token to a closed channel")
	}
	c.chSync <- tokenValue
	return nil
}

func (c *ChanSync) Wait() error {
	return nil
}

func (c *ChanSync) Length() (int, error) {
	if c.chSync == nil {
		return 0, errors.New("nil struct")
	}
	return len(c.chSync), nil
}

func NewChanSync(size int) ChanSync {
	c := ChanSync{}
	c.Size = size
	c.chSync = make(chan int, c.Size)
	c.open = true
	return c
}

func (c *ChanSync) Close() {
	close(c.chSync)
	c.open = false
}
