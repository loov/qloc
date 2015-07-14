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
func (xs Counts) Swap(i, j int) { xs[i], xs[j] = xs[j], xs[i] }

type ByCode struct{ Counts }

func (xs ByCode) Less(i, j int) bool { return xs.Counts[i].Code > xs.Counts[j].Code }

type ByExt struct{ Counts }

func (xs ByExt) Less(i, j int) bool { return xs.Counts[i].Ext < xs.Counts[j].Ext }
