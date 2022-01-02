package patch_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/patch"
)

type Car struct {
	Make  string
	Model string
	Year  int
}

type Parent struct {
	Name     string
	Child    Child
	ChildPtr *Child
}

type Child struct {
	Name string
}

var _ = Describe("Object", func() {
	It("can modify", func() {
		car := &Car{
			Make:  "Ford",
			Model: "F150",
			Year:  2000,
		}
		update := patch.NewObject(
			patch.ObjectSetField("Make", "Tesla"))

		_, ok := update.Apply(reflect.ValueOf(car))
		Expect(ok).To(BeTrue())
		Expect(car.Make).To(Equal("Tesla"))
	})
	It("can set child field", func() {
		parent := &Parent{
			Name: "parent",
			Child: Child{
				Name: "child",
			},
		}
		update := patch.NewObject(
			patch.ObjectSetField("Child",
				patch.NewObject(
					patch.ObjectSetField("Name", "test"))))

		_, ok := update.Apply(reflect.ValueOf(parent))
		Expect(ok).To(BeTrue())
		Expect(parent.Child.Name).To(Equal("test"))
	})
	It("can set child ptr", func() {
		parent := &Parent{
			Name: "parent",
		}
		update := patch.NewObject(
			patch.ObjectSetField("ChildPtr",
				patch.NewObject(
					patch.ObjectSetField("Name", "test"))))

		_, ok := update.Apply(reflect.ValueOf(parent))
		Expect(ok).To(BeTrue())
		Expect(parent.ChildPtr).ToNot(BeNil())
		Expect(parent.ChildPtr.Name).To(Equal("test"))
	})
})
