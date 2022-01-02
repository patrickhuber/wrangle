package patch_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/patch"
)

var _ = Describe("Slice", func() {
	Describe("Apply", func() {
		When("primative slice", func() {
			It("can add", func() {
				slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
				update := patch.NewSlice(
					patch.SliceAppend(10))
				val, ok := update.Apply(reflect.ValueOf(slice))
				Expect(ok).To(BeTrue())
				Expect(val.IsNil()).To(BeFalse())
				Expect(val.Len()).To(Equal(len(slice) + 1))
			})
			It("can remove at", func() {
				slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
				update := patch.NewSlice(
					patch.SliceRemoveAt(2))
				val, ok := update.Apply(reflect.ValueOf(slice))
				Expect(ok).To(BeTrue())
				Expect(val.IsNil()).To(BeFalse())
				Expect(val.Len()).To(Equal(len(slice) - 1))
				Expect(val.Index(2).Interface()).To(Equal(4))

			})
			It("can remove", func() {
				slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
				length := len(slice)
				update := patch.NewSlice(
					patch.SliceRemove(func(v reflect.Value) bool {
						i := v.Int()
						return i == int64(8)
					}))
				val, ok := update.Apply(reflect.ValueOf(slice))
				Expect(ok).To(BeTrue())
				Expect(val.IsNil()).To(BeFalse())
				Expect(val.Len()).To(Equal(length - 1))
			})
			It("can modify", func() {
				slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
				update := patch.NewSlice(
					patch.SliceModifyAt(5, 10))
				val, ok := update.Apply(reflect.ValueOf(slice))
				Expect(ok).To(BeTrue())
				Expect(val.IsNil()).To(BeFalse())
				Expect(val.Len()).To(Equal(len(slice)))
				Expect(val.Index(5).Interface()).To(Equal(10))
			})
		})
	})
})
