package main

type Runner struct {
	Id          int64  `json:"id"`
	Active      bool   `json:"active"`
	Name        string `json:"name"`
	Online      bool   `json:"online"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type Runners []Runner

func (t Runners) Filter(f func(r Runner, index int) bool) Runners {
	ret := make(Runners, 0)
	for i, item := range t {
		if f(item, i) {
			ret = append(ret, item)
		}
	}
	return ret
}
