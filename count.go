package main

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
	Ext    string
	Files  int
	Binary int
	Blank  int
	Code   int
}

func (c *Count) Add(s *Count) {
	c.Ext = s.Ext
	c.Files += s.Files
	c.Binary += s.Binary
	c.Blank += s.Blank
	c.Code += s.Code
}

type Counts []*Count

func (xs Counts) Len() int      { return len(xs) }
func (xs Counts) Swap(i, k int) { xs[i], xs[k] = xs[k], xs[i] }

type ByCode struct{ Counts }

func (xs ByCode) Less(i, k int) bool {
	a, b := xs.Counts[i], xs.Counts[k]
	if a.Code == b.Code {
		if a.Binary == b.Binary {
			return a.Ext < b.Ext
		}
		return a.Binary > b.Binary
	}
	return a.Code > b.Code
}

type ByExt struct{ Counts }

func (xs ByExt) Less(i, k int) bool {
	a, b := xs.Counts[i], xs.Counts[k]
	if a.Ext == b.Ext {
		if a.Code == b.Code {
			if a.Binary == b.Binary {
				return a.Blank < b.Blank
			}
			return a.Binary > b.Binary
		}
	}
	return a.Ext < b.Ext
}
