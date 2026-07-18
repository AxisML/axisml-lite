package core_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/axisml/axisml-lite/axisml-core/pkg/core"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := core.Load()
	require.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 5432, cfg.Database.Port)
	assert.Equal(t, "axisml", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "http://zot:5000", cfg.OCI.Endpoint)
	assert.Equal(t, "admin", cfg.OCI.AdminUser)
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("AXISML_DATABASE_HOST", "axisml-database")
	t.Setenv("AXISML_OCI_ENDPOINT", "http://zot.internal:5000")
	t.Setenv("AXISML_LOG_FORMAT", "console")
	t.Setenv("AXISML_WORKLOAD_TENANT_PREFIX", "true")
	cfg, err := core.Load()
	require.NoError(t, err)
	assert.Equal(t, "axisml-database", cfg.Database.Host)
	assert.Equal(t, "http://zot.internal:5000", cfg.OCI.Endpoint)
	assert.Equal(t, "console", cfg.Log.Format)
	assert.True(t, cfg.Workload.TenantPrefix)
}

func TestLoad_SecretFromFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/oci-pw"
	require.NoError(t, os.WriteFile(path, []byte("zot-admin\n"), 0o600))
	t.Setenv("AXISML_OCI_ADMIN_PASSWORD_FILE", path)
	cfg, err := core.Load()
	require.NoError(t, err)
	assert.Equal(t, "zot-admin", cfg.OCI.AdminPassword)
}

func TestLoadWithOptions_CustomPrefixOnly(t *testing.T) {
	t.Setenv("AXISML_DATABASE_HOST", "axisml-database")
	t.Setenv("AXISML_DATABASE_PORT", "15432")
	t.Setenv("AXISML_WORKLOAD_TENANT_PREFIX", "false")
	t.Setenv("AXISML_GPU_DEVICES", "0")
	t.Setenv("AIOSML_DATABASE_HOST", "aiosml-database")
	t.Setenv("AIOSML_DATABASE_PORT", "25432")
	t.Setenv("AIOSML_WORKLOAD_TENANT_PREFIX", "true")
	t.Setenv("AIOSML_GPU_DEVICES", "2,3")

	cfg, err := core.LoadWithOptions(core.LoadOptions{EnvPrefix: "AIOSML"})
	require.NoError(t, err)
	assert.Equal(t, "aiosml-database", cfg.Database.Host)
	assert.Equal(t, 25432, cfg.Database.Port)
	assert.True(t, cfg.Workload.TenantPrefix)
	assert.Equal(t, "2,3", cfg.GPU.Devices)
	assert.Equal(t, "json", cfg.Log.Format)
}

func TestLoadWithOptions_EmptyPrefixUsesDefault(t *testing.T) {
	t.Setenv("AXISML_DATABASE_HOST", "axisml-database")
	t.Setenv("AIOSML_DATABASE_HOST", "aiosml-database")

	cfg, err := core.LoadWithOptions(core.LoadOptions{})
	require.NoError(t, err)
	assert.Equal(t, "axisml-database", cfg.Database.Host)
}

func TestLoadWithOptions_CustomPrefixSecretFilesOverrideEnvironment(t *testing.T) {
	dir := t.TempDir()
	databasePath := dir + "/database-password"
	ociPath := dir + "/oci-password"
	require.NoError(t, os.WriteFile(databasePath, []byte("database-file-secret\n"), 0o600))
	require.NoError(t, os.WriteFile(ociPath, []byte("oci-file-secret\n"), 0o600))

	t.Setenv("AIOSML_DATABASE_PASSWORD", "database-env-secret")
	t.Setenv("AIOSML_DATABASE_PASSWORD_FILE", databasePath)
	t.Setenv("AIOSML_OCI_ADMIN_PASSWORD", "oci-env-secret")
	t.Setenv("AIOSML_OCI_ADMIN_PASSWORD_FILE", ociPath)
	t.Setenv("AXISML_DATABASE_PASSWORD", "axisml-secret")
	t.Setenv("AXISML_OCI_ADMIN_PASSWORD", "axisml-secret")

	cfg, err := core.LoadWithOptions(core.LoadOptions{EnvPrefix: "AIOSML"})
	require.NoError(t, err)
	assert.Equal(t, "database-file-secret", cfg.Database.Password)
	assert.Equal(t, "oci-file-secret", cfg.OCI.AdminPassword)
}

func TestLoadWithOptions_InvalidPrefix(t *testing.T) {
	for _, prefix := range []string{"aiosml", "1AIOSML", "AIOS-ML", "AIOSML_"} {
		t.Run(prefix, func(t *testing.T) {
			_, err := core.LoadWithOptions(core.LoadOptions{EnvPrefix: prefix})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid environment variable prefix")
			assert.Contains(t, err.Error(), prefix)
		})
	}
}

func TestLoadWithOptions_CustomPrefixTypeParsingError(t *testing.T) {
	t.Setenv("AIOSML_DATABASE_PORT", "not-an-integer")

	_, err := core.LoadWithOptions(core.LoadOptions{EnvPrefix: "AIOSML"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode config")
}

func TestLoadWithOptions_CustomPrefixSecretFileErrorNamesSelectedVariable(t *testing.T) {
	t.Setenv("AIOSML_DATABASE_PASSWORD_FILE", t.TempDir()+"/missing")

	_, err := core.LoadWithOptions(core.LoadOptions{EnvPrefix: "AIOSML"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "AIOSML_DATABASE_PASSWORD_FILE")
	assert.NotContains(t, err.Error(), "AXISML_DATABASE_PASSWORD_FILE")
}

func TestPostgresDSN(t *testing.T) {
	cfg, err := core.Load()
	require.NoError(t, err)
	dsn := cfg.PostgresDSN()
	for _, want := range []string{"host=localhost", "port=5432", "dbname=axisml", "sslmode=disable"} {
		assert.Contains(t, dsn, want)
	}
}
