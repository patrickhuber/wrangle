package patch_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/patch"
)

var _ = Describe("Primative", func() {
	Describe("StringUpdate", func() {
		When("are different", func() {
			It("returns new value and true", func() {
				update := patch.NewString("new")
				v, ok := update.Apply(reflect.ValueOf("old"))
				Expect(ok).To(BeTrue())
				Expect(v.String()).To(Equal("new"))
			})
		})
		When("are the same", func() {
			It("returns false", func() {
				update := patch.NewString("old")
				_, ok := update.Apply(reflect.ValueOf("old"))
				Expect(ok).To(BeFalse())
			})
		})
	})
	Describe("IntUpdate", func() {
		When("are different", func() {
			It("returns new value and true", func() {
				update := patch.NewInt(10)
				v, ok := update.Apply(reflect.ValueOf(20))
				Expect(ok).To(BeTrue())
				Expect(int(v.Int())).To(Equal(10))
			})
		})
		When("are the same", func() {
			It("returns false", func() {
				update := patch.NewInt(10)
				_, ok := update.Apply(reflect.ValueOf(10))
				Expect(ok).To(BeFalse())
			})
		})
	})
	Describe("BoolUpdate", func() {
		When("are different", func() {
			It("returns new value and true", func() {
				update := patch.NewBool(false)
				v, ok := update.Apply(reflect.ValueOf(true))
				Expect(ok).To(BeTrue())
				Expect(v.Bool()).To(Equal(false))
			})
		})
		When("are the same", func() {
			It("returns false", func() {
				update := patch.NewBool(true)
				_, ok := update.Apply(reflect.ValueOf(true))
				Expect(ok).To(BeFalse())
			})
		})
	})
})
