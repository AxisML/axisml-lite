// Package core holds the axisml-core composition root: configuration, the
// config-backed Cluster Manager providers (single default ResourcePool +
// Tenant), the GORM/DB coordination and the assembly that mounts the three
// System modules (Cluster Manager, Compute Service, Artifact Hub) plus the
// in-process Standalone Runtime on one router.
package core

import (
	"fmt"

	"github.com/axisml/axisml-lite/axisml-core/internal/configutil"
)

// Config is the axisml-core process configuration. Lite runs under Docker
// Compose, where the environment is the idiomatic config mechanism, so this
// binary reads NO config file. Load uses the AXISML prefix by default; embedded
// callers can select another prefix with LoadWithOptions. Secret keys also
// accept a corresponding <PREFIX>_<KEY>_FILE variable. A single binary fronts
// all three System modules, so the database and OCI settings are shared.
// Everything operational (ports, reconcile cadence, GC, filesystem paths,
// Docker network) is a fixed constant — see consts.go.
type Config struct {
	Common `mapstructure:",squash"`

	OCI      OCI      `mapstructure:"oci"`
	GPU      GPU      `mapstructure:"gpu"`
	Workload Workload `mapstructure:"workload"`
}

// Workload controls physical workload resource naming.
type Workload struct {
	TenantPrefix bool `mapstructure:"tenant_prefix" default:"false" doc:"Prefix physical workload names with a readable, collision-resistant tenant token"`
}

// GPU configures single-host GPU scheduling. Devices names the physical GPU
// indices AxisML may schedule onto as a comma list ("0,1,2"), which turns on
// managed scheduling (pin to a free card, wait when none is free). Empty leaves
// managed scheduling off: GPU workloads use Docker's default count-based request.
type GPU struct {
	Devices string `mapstructure:"devices" doc:"Physical GPU indices to schedule onto (comma list, e.g. 0,1,2); empty falls back to Docker's default count-based GPU request"`
}

// OCI is the artifact registry (zot) connection. The scheme is derived from the
// endpoint URL, so there is no separate scheme key.
type OCI struct {
	Endpoint      string `mapstructure:"endpoint" default:"http://zot:5000" doc:"OCI registry endpoint (full URL; scheme derived from it)"`
	AdminUser     string `mapstructure:"admin_user" default:"admin" doc:"OCI registry admin username"`
	AdminPassword string `mapstructure:"admin_password" secret:"true" doc:"OCI registry admin password"`
}

// DefaultEnvPrefix is the environment variable prefix used by Load and by
// LoadWithOptions when EnvPrefix is empty.
const DefaultEnvPrefix = "AXISML"

// LoadOptions controls how axisml-core loads its process configuration.
type LoadOptions struct {
	// EnvPrefix does not include a trailing underscore, for example AXISML or
	// AIOSML. An empty value uses DefaultEnvPrefix.
	EnvPrefix string
}

// Load resolves the configuration from defaults < AXISML_ env < AXISML_<KEY>_FILE
// secret files. axisml-core reads no config file (env-only).
func Load() (Config, error) {
	return LoadWithOptions(LoadOptions{EnvPrefix: DefaultEnvPrefix})
}

// LoadWithOptions resolves the configuration from defaults, then environment
// variables under the selected prefix, then the selected prefix's secret
// *_FILE variables. It does not read or merge variables under other prefixes.
func LoadWithOptions(opts LoadOptions) (Config, error) {
	prefix := opts.EnvPrefix
	if prefix == "" {
		prefix = DefaultEnvPrefix
	}
	var c Config
	if err := configutil.Load(&c, prefix); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Validate is invoked by the loader (fail-fast).
func (c Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database.host is required")
	}
	return nil
}
