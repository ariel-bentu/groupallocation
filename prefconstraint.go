package main

import (
	"fmt"
)

type PrefConstraint struct {
	id       int
	desc     string
	Level    int
	weight   int
	pupil    int
	refPupil int
	prefPrio int
}

func NewPrefConstraint(id int, desc string, weight int, pupil int, refPupil int, prefPrio int) Constraint {
	return &PrefConstraint{id: id, desc: desc, weight: weight,
		pupil: pupil, refPupil: refPupil, prefPrio: prefPrio, Level: 1}
}
func (c *PrefConstraint) ID() int {
	return c.id
}

func (c *PrefConstraint) Description() string {
	return c.desc
}
func (c *PrefConstraint) Weight() int {
	return c.weight
}
func (c *PrefConstraint) Validate(ec *ExecutionContext) bool {
	satisfied := ec.pupils[c.refPupil].group == ec.pupils[c.pupil].group
	if !satisfied {
		ec.pupils[c.pupil].score -= c.weight
	}
	return satisfied
}

func (c *PrefConstraint) Message(ec *ExecutionContext) string {
	satisfied := ec.pupils[c.refPupil].group == ec.pupils[c.pupil].group
	if !satisfied {
		return fmt.Sprintf("לא כובדה העדפה %d של %s (%s)", c.prefPrio, ec.pupils[c.pupil].name, ec.pupils[c.refPupil].name)
	}
	return ""
}
