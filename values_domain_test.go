package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_testMain(t *testing.T) {
	RegisterFailHandler(Fail)

	tests := []struct {
		name string
	}{
		{name: "main"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valuesDomain := NewValuesDomainMain(3, 100)
			val := valuesDomain.FirstValue(10)
			Expect(val).To(Equal(0))

			val = valuesDomain.NextValue(10, val)
			Expect(val).To(Equal(1))

			valuesDomain.PushDomainRestriction(2, 10, []int{1, 2})

			val = valuesDomain.FirstValue(10)
			Expect(val).To(Equal(1))

			val = valuesDomain.NextValue(10, val)
			Expect(val).To(Equal(2))

			valuesDomain.PushDomainRestriction(4, 10, []int{0, 2})
			val = valuesDomain.FirstValue(10)
			Expect(val).To(Equal(2))
			val = valuesDomain.NextValue(val, 10)
			Expect(val).To(Equal(-1))

			valuesDomain.PopAllDomainRestriction(4)
			val = valuesDomain.FirstValue(10)
			Expect(val).To(Equal(1))

			val = valuesDomain.NextValue(10, val)
			Expect(val).To(Equal(2))

		})
	}
}
