package di_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/di"
)

type SampleStruct struct {
	name string
}

type SampleInterface interface {
	Name() string
}

type DependencyInterface interface {
	Name() string
}

func NewSample(name string) SampleInterface {
	return &SampleStruct{
		name: name,
	}
}

func (s *SampleStruct) Name() string {
	return s.name
}

var _ = Describe("Container", func() {
	It("can resolve type", func() {
		container := di.NewContainer()
		sample := NewSample("test")
		container.RegisterInstance(reflect.TypeOf((*SampleInterface)(nil)).Elem(), sample)
		instance, err := container.Resolve(reflect.TypeOf((*SampleInterface)(nil)).Elem())
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		_, ok := instance.(SampleInterface)
		Expect(ok).To(BeTrue())
	})
	It("can register constructor", func() {
		container := di.NewContainer()
		name := "myname"
		container.RegisterInstance(reflect.TypeOf((*string)(nil)).Elem(), name)
		container.RegisterConstructor(NewSample)
		instance, err := container.Resolve(reflect.TypeOf((*SampleInterface)(nil)).Elem())
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		sample, ok := instance.(SampleInterface)
		Expect(ok).To(BeTrue())
		Expect(sample.Name()).To(Equal(name))
	})
})
