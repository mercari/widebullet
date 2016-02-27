package wlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccessLogger(t *testing.T) {
	assert := assert.New(t)

	al := AccessLogger("debug")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Debug, al.Level)

	al = AccessLogger("info")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Info, al.Level)

	al = AccessLogger("notice")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Notice, al.Level)

	al = AccessLogger("warn")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Warn, al.Level)

	al = AccessLogger("error")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Error, al.Level)

	al = AccessLogger("crit")
	assert.Equal(Stdout, al.Rdr)
	assert.Equal(Crit, al.Level)
}

func TestNewErrorLogger(t *testing.T) {
	assert := assert.New(t)

	el := ErrorLogger("debug")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Debug, el.Level)

	el = ErrorLogger("info")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Info, el.Level)

	el = ErrorLogger("notice")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Notice, el.Level)

	el = ErrorLogger("warn")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Warn, el.Level)

	el = ErrorLogger("error")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Error, el.Level)

	el = ErrorLogger("crit")
	assert.Equal(Stderr, el.Rdr)
	assert.Equal(Crit, el.Level)
}

func TestLevel2String(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("debug", level2String(Debug))
	assert.Equal("info", level2String(Info))
	assert.Equal("notice", level2String(Notice))
	assert.Equal("warn", level2String(Warn))
	assert.Equal("error", level2String(Error))
	assert.Equal("crit", level2String(Crit))
}
