package vault_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/vault"
	"github.com/stretchr/testify/require"
	vaultcontainer "github.com/testcontainers/testcontainers-go/modules/vault"
)

const (
	vaultVersion = "hashicorp/vault:1.13.3"
	rootToken    = "root-token"
)

// setupVaultContainer starts a Vault testcontainer and returns the address and root token
func setupVaultContainer(ctx context.Context, t *testing.T) (string, string, func()) {
	container, err := vaultcontainer.Run(ctx,
		vaultVersion,
		vaultcontainer.WithToken(rootToken),
	)
	require.NoError(t, err)

	addr, err := container.HttpHostAddress(ctx)
	require.NoError(t, err)

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	// HttpHostAddress already includes the protocol
	return addr, rootToken, cleanup
}

// setupAppRole configures AppRole authentication in Vault and returns role_id and secret_id
func setupAppRole(ctx context.Context, t *testing.T, vaultAddr, token string) (string, string) {
	// Create Vault client for setup
	client, err := vault.NewVault(&vault.VaultOptions{
		Address: vaultAddr,
		Token:   token,
		Path:    "secret",
	})
	require.NoError(t, err)

	// Enable AppRole auth method
	err = client.Client().Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	require.NoError(t, err)

	// Create a policy that allows access to secret/*
	policy := `path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}`
	err = client.Client().Sys().PutPolicy("test-policy", policy)
	require.NoError(t, err)

	// Create an AppRole with the policy
	_, err = client.Client().Logical().Write("auth/approle/role/test-role", map[string]interface{}{
		"token_policies": []string{"test-policy"},
		"token_ttl":      "1h",
		"token_max_ttl":  "4h",
	})
	require.NoError(t, err)

	// Get role_id
	roleIDResp, err := client.Client().Logical().Read("auth/approle/role/test-role/role-id")
	require.NoError(t, err)
	roleID := roleIDResp.Data["role_id"].(string)

	// Generate secret_id
	secretIDResp, err := client.Client().Logical().Write("auth/approle/role/test-role/secret-id", nil)
	require.NoError(t, err)
	secretID := secretIDResp.Data["secret_id"].(string)

	return roleID, secretID
}

// runStoreTests runs common store tests for the given store instance
func runStoreTests(t *testing.T, store stores.Store, keyPrefix string) {
	t.Run("set", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: keyPrefix + "-test",
		}}
		value := keyPrefix + "-value"
		err := store.Set(key, value)
		require.NoError(t, err)

		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, value, v)
	})

	t.Run("get", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: keyPrefix + "-test",
		}}
		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, keyPrefix+"-value", v)
	})

	t.Run("list", func(t *testing.T) {
		items, err := store.List()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(items), 1)
	})
}

func TestVaultWithToken(t *testing.T) {

	var INTEGRATION = "INTEGRATION"
	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
		return
	}

	// only run docker dependent tests if docker is installed
	if !dockerIsInstalled() {
		t.Skipf("skipping test: TestVaultWithToken tests currently only supported where docker is installed")
		return
	}

	ctx := context.Background()
	// Setup Vault container
	address, token, cleanup := setupVaultContainer(ctx, t)
	defer cleanup()

	// Create store with token authentication
	factory := vault.NewFactory()
	store, err := factory.Create(map[string]any{
		"address": address,
		"token":   token,
		"path":    "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	runStoreTests(t, store, "token")
}

func TestVaultWithAppRole(t *testing.T) {

	var INTEGRATION = "INTEGRATION"
	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
		return
	}

	// only run docker dependent tests if docker is installed
	if !dockerIsInstalled() {
		t.Skipf("skipping test: TestVaultWithAppRole tests currently only supported where docker is installed")
		return
	}

	ctx := context.Background()

	// Setup Vault container
	address, token, cleanup := setupVaultContainer(ctx, t)
	defer cleanup()

	// Setup AppRole authentication
	roleID, secretID := setupAppRole(ctx, t, address, token)

	// Create store with AppRole authentication
	factory := vault.NewFactory()
	store, err := factory.Create(map[string]any{
		"address":   address,
		"role_id":   roleID,
		"secret_id": secretID,
		"path":      "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	runStoreTests(t, store, "approle")
}

func dockerIsInstalled() bool {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Docker is not installed or not found in PATH.")
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Output: %s\n", output)
		return false
	}
	return true
}
