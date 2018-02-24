package extargsparse

type parseState struct {
	logObject
	cmdpaths      []*parserCompat
	curidx        int
	curcharidx    int
	shortcharargs int
	longargs      int
	keyidx        int
	validx        int
	args          []string
	ended         int
	longprefix    string
	shortprefix   string
	parseall      bool
	leftargs      []string
	bundlemode    bool
}

func newParseState(args []string, maincmd *parserCompat, optattr *ExtArgsOptions) *parseState {
	assert_test(maincmd != nil, "maincmd can not set nil")
	self := &parseState{logObject: *newLogObject("extargsparse")}
	if optattr == nil {
		optattr = NewExtArgsOptions("{}")
	}

	self.cmdpaths = make([]*parserCompat, 0)
	self.cmdpaths = append(self.cmdpaths, maincmd)
	self.curidx = 0
	self.curcharidx = -1
	self.shortcharargs = -1
	self.longargs = -1
	self.keyidx = -1
	self.validx = -1
	self.args = args
	self.ended = 0
	self.longprefix = optattr.GetString("longprefix")
	self.shortprefix = optattr.GetString("shortprefix")
	if len(self.shortprefix) == 0 || len(self.longprefix) == 0 ||
		self.shortprefix == self.longprefix {
		self.bundlemode = true
	} else {
		self.bundlemode = false
	}
	self.parseall = optattr.GetBool("parseall")
	self.leftargs = make([]string, 0)

	return self
}
