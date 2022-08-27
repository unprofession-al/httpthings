package main

type TodoSet map[string]*Todo

func NewTodoSet() *TodoSet {
	return &TodoSet{}
}

func (ts TodoSet) AsSlice() []*Todo {
	out := []*Todo{}
	for _, v := range ts {
		out = append(out, v)
	}
	return out
}

func (ts TodoSet) Add(in *Todo) bool {
	if _, found := ts[in.Name]; found {
		return false
	}
	ts[in.Name] = in
	return true
}

type Todo struct {
	Name        string `json:"name" yaml:"name" jsonschema:"minLength=3"`
	Description string `json:"description" yaml:"description"`
	Done        bool   `json:"done" yaml:"done"`
}

func (t *Todo) Finish() {
	t.Done = true
}

type TodoRequest struct {
	Name        string `json:"name" yaml:"name" jsonschema:"minLength=3"`
	Description string `json:"description" yaml:"description"`
}

func (tr *TodoRequest) AsTodo() *Todo {
	return &Todo{
		Name:        tr.Name,
		Description: tr.Description,
		Done:        false,
	}
}
