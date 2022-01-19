package patch_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/patch"
)

var _ = Describe("Map", func() {
	When("primative map", func() {
		var (
			m map[string]int
		)
		BeforeEach(func() {
			m = map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			}
		})
		It("can add", func() {
			p := patch.NewMap(
				patch.MapSet("four", 4))
			result, ok := p.Apply(reflect.ValueOf(m))
			Expect(ok).To(BeTrue())
			Expect(result.IsNil()).To(BeFalse())
			Expect(len(result.MapKeys())).To(Equal(4))
			v, ok := m["four"]
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(4))
		})
		It("can remove", func() {
			p := patch.NewMap(
				patch.MapRemove("two"))
			result, ok := p.Apply(reflect.ValueOf(m))
			Expect(ok).To(BeTrue())
			Expect(result.IsNil()).To(BeFalse())
			Expect(len(result.MapKeys())).To(Equal(2))
			_, ok = m["two"]
			Expect(ok).To(BeFalse())
		})
		It("can set", func() {
			p := patch.NewMap(
				patch.MapSet("two", 4))
			result, ok := p.Apply(reflect.ValueOf(m))
			Expect(ok).To(BeTrue())
			Expect(result.IsNil()).To(BeFalse())
			Expect(len(result.MapKeys())).To(Equal(3))
			v, ok := m["two"]
			Expect(ok).To(BeTrue())
			Expect(v).To(Equal(4))
		})
	})
})
