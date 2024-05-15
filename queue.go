package main

type Queue []string

func (q *Queue) Push(x string) {
	*q = append(*q, x)
}

func (q *Queue) Size() int {
	h := *q
	l := len(h)
	return l
}

func (q *Queue) Pop() string {
	h := *q
	var el string
	l := len(h)
	if l == 0 {
		return ""
	}
	el, *q = h[0], h[1:l]

	return el
}

func NewQueue() *Queue {
	return &Queue{}
}
