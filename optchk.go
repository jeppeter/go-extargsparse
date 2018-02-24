package extargsparse

type optCheck struct {
	longopt  []string
	shortopt []string
	varname  []string
}

func newOptCheck() *optCheck {
	self := &optCheck{}
	self.reset()
	return self
}

func (self *optCheck) reset() {
	self.longopt = make([]string, 0)
	self.shortopt = make([]string, 0)
	self.varname = make([]string, 0)
	return
}

func (self *optCheck) Copy(other *optCheck) {
	var k string
	self.reset()
	for _, k = range other.longopt {
		self.longopt = append(self.longopt, k)
	}
	for _, k = range other.shortopt {
		self.shortopt = append(self.shortopt, k)
	}

	for _, k = range other.varname {
		self.varname = append(self.varname, k)
	}
	return
}

func (self *optCheck) check_in_array(sarr []string, s string) bool {
	for _, k := range sarr {
		if k == s {
			return true
		}
	}
	return false
}

func (self *optCheck) AddAndCheck(typename string, value string) bool {
	if typename == "longopt" {
		if self.check_in_array(self.longopt, value) {
			return false
		}
		self.longopt = append(self.longopt, value)
		return true
	} else if typename == "shortopt" {
		if self.check_in_array(self.shortopt, value) {
			return false
		}
		self.shortopt = append(self.shortopt, value)
		return true
	} else if typename == "varname" {
		if self.check_in_array(self.varname, value) {
			return false
		}
		self.varname = append(self.varname, value)
		return true
	}
	return false
}
