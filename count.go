package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

type CountByExt map[string]*Count

func (xs CountByExt) Add(s *Count) {
	r, ok := xs[s.Ext]
	if !ok {
		xs[s.Ext] = s
		return
	}
	r.Add(s)
}

type Count struct {
	Ext   string
	Files int
	Blank int
	Code  int
}

func (c *Count) Add(s *Count) {
	c.Ext = s.Ext
	c.Files += s.Files
	c.Blank += s.Blank
	c.Code += s.Code
}

type Counts []*Count

func (xs Counts) Len() int      { return len(xs) }
func (xs Counts) Swap(i, j int) { xs[i], xs[j] = xs[j], xs[i] }

type ByCode struct{ Counts }

func (xs ByCode) Less(i, j int) bool { return xs.Counts[i].Code > xs.Counts[i].Code }

type ByExt struct{ Counts }

func (xs ByExt) Less(i, j int) bool { return xs.Counts[i].Ext < xs.Counts[i].Ext }

var (
	BinaryExtCacheWrite = sync.Mutex{}
	BinaryExtCache      = atomic.Value{}

	ErrEmptyFile  = errors.New("empty file")
	ErrBinaryFile = errors.New("binary file")
)

func init() { BinaryExtCache.Store(make(map[string]bool)) }
func IsBinaryExt(ext string) bool {
	if ShouldExamine(ext) {
		return false
	}
	cache := BinaryExtCache.Load().(map[string]bool)
	return cache[ext]
}

func AddBinaryExt(ext string) {
	BinaryExtCacheWrite.Lock()
	defer BinaryExtCacheWrite.Unlock()

	cache := BinaryExtCache.Load().(map[string]bool)
	next := make(map[string]bool)
	for e := range cache {
		next[e] = true
	}
	next[ext] = true
	BinaryExtCache.Store(next)
}
