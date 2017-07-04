package main

//Constraint is a...
type Constraint interface {
	ID() int

	Weight() int

	Members() []int

	Description() string
	Validate(ec *ExecutionContext) bool

	Message(ec *ExecutionContext) string
}

type ByWeight []Constraint

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].Weight() < a[j].Weight() }

type ByMembers []Constraint

func (a ByMembers) Len() int           { return len(a) }
func (a ByMembers) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMembers) Less(i, j int) bool { return len(a[i].Members()) > len(a[j].Members()) }
