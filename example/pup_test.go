package pup

import (
	"io"
	"testing"

	. "github.com/stevegt/goadapt"
)

type MockReadWriteCloser struct {
	readbuf  []byte
	writebuf []byte
	readpos  int
}

func (c *MockReadWriteCloser) Read(out []byte) (n int, err error) {
	if c.readpos >= len(c.readbuf) {
		return 0, io.EOF
	}
	// set end to the end of the readbuf or out, depending on which comes
	// first
	end := c.readpos + len(out)
	if end > len(c.readbuf) {
		end = len(c.readbuf)
	}
	n = copy(out, c.readbuf[c.readpos:end])
	c.readpos += n
	return n, nil
}

func (c *MockReadWriteCloser) Write(data []byte) (n int, err error) {
	c.writebuf = append(c.writebuf, data...)
	return len(data), nil
}

func (c *MockReadWriteCloser) Close() error {
	return nil
}

var s1hash = "somehash"
var s1content = "first line\nsecond line\n"
var s1 = Spf("%s\n%s", s1hash, s1content)

func TestStream(t *testing.T) {
	rwc := &MockReadWriteCloser{readbuf: []byte(s1)}
	err := handleStream(rwc)
	Tassert(t, err == nil, "handleStream %v", err)
	Tassert(t, string(rwc.writebuf) == s1content, "writebuf '%v'", string(rwc.writebuf))
}
