package main

type TodoSet map[string]*Todo

func NewTodoSet() *TodoSet {
	return &TodoSet{}
}

func (ts TodoSet) Add(n, d string) bool {
	if _, found := ts[n]; found {
		return false
	}
	ts[n] = &Todo{
		Name:        n,
		Description: d,
		Done:        false,
	}
	return true
}

type Todo struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Done        bool   `json:"done" yaml:"done"`
}

func (t *Todo) Finish() {
	t.Done = true
}
