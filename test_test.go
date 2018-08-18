package main

/*
func TestValidate(t *testing.T) {
	ec := NewExecutionContext()
	ec.groupsCount = 2

	t.Run("SperateGroup", func(t *testing.T) {
		addPupil(ec, "ariel", true, 0, 1, -1, -1)
		addPupil(ec, "maayan", false, 0, 1, -1, -1)

		sg := addGroup(ec, 1, "t1", false, 125, 0, 1)
		if len(ec.pupils) == 0 {
			t.Errorf("Wrong number of Pupils %d", len(ec.pupils))
		}
		sat := sg.Validate(ec)
		if sat {
			t.Errorf("t1 should NOT be  valid")
		}
		addPupil(ec, "Yuval", true, 1, 1, -1, -1)
		sg.AddMember(2)
		addPupil(ec, "Guy", true, 0, 1, -1, -1)
		sg.AddMember(3)
		inbal := addPupil(ec, "Inbal", false, 1, 1, -1, -1)
		sg.AddMember(4)
		sat = sg.Validate(ec)
		if !sat {
			t.Errorf("t1 should be valid")
		}
		if inbal.score != 0 {

			t.Errorf("inbal should have score 0")
		}

	})
}

func addPupil(ec *ExecutionContext, name string, isMale bool, group int, ref0 int, ref1 int, ref2 int) *Pupil {
	p := &Pupil{name: name, isMale: isMale, group: group}
	ec.pupils = append(ec.pupils, p)
	return p
}

func addGroup(ec *ExecutionContext, id int, name string, isUnite bool, weight int, members ...int) *SubGroupConstraint {
	sg := &SubGroupConstraint{id: id, desc: name, IsUnite: isUnite, weight: weight, Level: 1}
	for i := 0; i < len(members); i++ {
		sg.AddMember(members[i])
	}

	ec.Constraints = append(ec.Constraints, sg)
	return sg
}
*/
