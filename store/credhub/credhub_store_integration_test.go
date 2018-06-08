// +build integration

package credhub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredHubIntegration(t *testing.T) {

	r := require.New(t)
	r.Equal(1, 2)
}
