package logging

import "sync"

type pcs struct {
	pcs []uintptr
}

var _pcsPool = &sync.Pool{
	New: func() interface{} {
		return &pcs{make([]uintptr, 64)}
	},
}

func newPcs() *pcs {
	return _pcsPool.Get().(*pcs)
}

func putPcs(p *pcs) {
	_pcsPool.Put(p)
}

type caller struct {
	ok   bool
	pc   uintptr
	file string
	line int
	fc   string
}

func (c *caller) GetOK() bool {
	return c.ok
}

func (c *caller) GetPC() uintptr {
	return c.pc
}

func (c *caller) GetFile() string {
	return c.file
}

func (c *caller) GetLine() int {
	return c.line
}

func (c *caller) GetFunc() string {
	return c.fc
}
