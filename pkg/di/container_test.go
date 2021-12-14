package di_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/di"
)

type SampleStruct struct {
	name string
}

type AggregateStruct struct {
	names []string
}

type SampleInterface interface {
	Name() string
}

type DependencyInterface interface {
	Name() string
}

type AggregateInterface interface {
	Names() []string
}

func NewSample(name string) SampleInterface {
	return &SampleStruct{
		name: name,
	}
}

func NewVariadic(dependencies ...DependencyInterface) AggregateInterface {
	names := []string{}
	for _, d := range dependencies {
		names = append(names, d.Name())
	}
	return &AggregateStruct{
		names: names,
	}
}

func NewAggregate(dependencies []DependencyInterface) AggregateInterface {
	names := []string{}
	for _, d := range dependencies {
		names = append(names, d.Name())
	}
	return &AggregateStruct{
		names: names,
	}
}

func NewWithNilError() (SampleInterface, error) {
	return &SampleStruct{
		name: "test",
	}, nil
}

func NewWithError() (SampleInterface, error) {
	return nil, fmt.Errorf("this is an error")
}

func TwoReturnTypes() (SampleInterface, AggregateInterface) {
	return nil, nil
}

func (s *SampleStruct) Name() string {
	return s.name
}

func (a *AggregateStruct) Names() []string {
	return a.names
}

var StringType = reflect.TypeOf((*string)(nil)).Elem()
var SampleInterfaceType = reflect.TypeOf((*SampleInterface)(nil)).Elem()
var DependencyInterfaceType = reflect.TypeOf((*DependencyInterface)(nil)).Elem()
var AggregateInterfaceType = reflect.TypeOf((*AggregateInterface)(nil)).Elem()

var _ = Describe("Container", func() {
	It("can resolve type", func() {
		container := di.NewContainer()
		sample := NewSample("test")
		container.RegisterInstance(SampleInterfaceType, sample)
		instance, err := container.Resolve(SampleInterfaceType)
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		_, ok := instance.(SampleInterface)
		Expect(ok).To(BeTrue())
	})
	It("can register constructor", func() {
		container := di.NewContainer()
		name := "myname"
		container.RegisterInstance(StringType, name)
		container.RegisterConstructor(NewSample)
		instance, err := container.Resolve(SampleInterfaceType)
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		sample, ok := instance.(SampleInterface)
		Expect(ok).To(BeTrue())
		Expect(sample.Name()).To(Equal(name))
	})
	It("can register array parameter", func() {
		container := di.NewContainer()
		dependencies := []*SampleStruct{
			{name: "sample 1"},
			{name: "sample 2"},
		}
		container.RegisterInstance(DependencyInterfaceType, dependencies[0])
		container.RegisterInstance(DependencyInterfaceType, dependencies[1])
		err := container.RegisterConstructor(NewAggregate)
		Expect(err).To(BeNil())
		instance, err := container.Resolve(AggregateInterfaceType)
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		a, ok := instance.(AggregateInterface)
		Expect(ok).To(BeTrue())
		Expect(a).ToNot(BeNil())
		Expect(len(a.Names())).To(Equal(2))
	})
	It("can register variadic parameter", func() {
		container := di.NewContainer()
		dependencies := []*SampleStruct{
			{name: "sample 1"},
			{name: "sample 2"},
		}
		container.RegisterInstance(DependencyInterfaceType, dependencies[0])
		container.RegisterInstance(DependencyInterfaceType, dependencies[1])
		err := container.RegisterConstructor(NewVariadic)
		Expect(err).To(BeNil())
		instance, err := container.Resolve(AggregateInterfaceType)
		Expect(err).To(BeNil())
		Expect(instance).ToNot(BeNil())
		_, ok := instance.(AggregateInterface)
		Expect(ok).To(BeTrue())
	})
	It("can invoke constructor that returns error", func() {
		container := di.NewContainer()
		err := container.RegisterConstructor(NewWithError)
		Expect(err).To(BeNil())
		i, err := container.Resolve(SampleInterfaceType)
		Expect(err).ToNot(BeNil())
		Expect(i).To(BeNil())
	})
	It("can invoke constructor that returns value and nil error", func() {
		container := di.NewContainer()
		err := container.RegisterConstructor(NewWithNilError)
		Expect(err).To(BeNil())
		i, err := container.Resolve(SampleInterfaceType)
		Expect(err).To(BeNil())
		Expect(i).ToNot(BeNil())
		_, ok := i.(SampleInterface)
		Expect(ok).To(BeTrue())
	})
	It("throws error when second return type is not error", func() {
		container := di.NewContainer()
		err := container.RegisterConstructor(TwoReturnTypes)
		Expect(err).ToNot(BeNil())
	})
	It("throws error when no return type", func() {
		container := di.NewContainer()
		err := container.RegisterConstructor(func() {})
		Expect(err).ToNot(BeNil())
	})
})
