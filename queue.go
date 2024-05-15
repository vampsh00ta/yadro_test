package main

type Queue []string

func (self *Queue) Push(x string) {
	*self = append(*self, x)
}
func (self *Queue) Size() int {
	h := *self
	l := len(h)
	return l
}
func (self *Queue) Pop() string {
	h := *self
	var el string
	l := len(h)
	if l == 0 {
		return ""
	}
	el, *self = h[0], h[1:l]

	return el
}

func NewQueue() *Queue {
	return &Queue{}
}
