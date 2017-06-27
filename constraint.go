package main

//Constraint is a...
type Constraint interface {
	ID() int

	Weight() int

	Description() string
	Validate(ec *ExecutionContext) bool

	Message(ec *ExecutionContext) string
}

type ByWeight []Constraint

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight() < a[j].Weight() }
