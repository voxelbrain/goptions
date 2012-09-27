package goptions

type MutexGroup []*Flag

func (mg MutexGroup) IsObligatory() bool {
	for _, flag := range mg {
		if flag.Obligatory {
			return true
		}
	}
	return false
}

func (mg MutexGroup) WasSpecified() bool {
	for _, flag := range mg {
		if flag.WasSpecified {
			return true
		}
	}
	return false
}

func (mg MutexGroup) IsValid() bool {
	c := 0
	for _, flag := range mg {
		if flag.WasSpecified {
			c++
		}
	}
	return c <= 1 && (!mg.IsObligatory() || c == 1)
}

func (mg MutexGroup) Names() []string {
	r := make([]string, len(mg))
	for i, flag := range mg {
		r[i] = flag.Name()
	}
	return r
}
