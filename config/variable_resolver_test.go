package config

import "testing"

type FakeResolver struct {
}

type interpolator struct {
}

func (resolver *FakeResolver) Resolve(document interface{}) interface{} {

	return nil
}

func TestVariableResolver(t *testing.T) {
	t.Run("CanResolveVariable", func(t *testing.T) {

	})
}
