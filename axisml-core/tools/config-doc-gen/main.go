// Command config-doc-gen prints the axisml-core (Lite) configuration reference
// section (Markdown) for assembly into docs/configuration.md. axisml-core reads
// no config file (env-only), so the section is rendered with envOnly=true.
package main

import (
	"fmt"
	"strings"

	"github.com/axisml/axisml-lite/axisml-core/internal/configutil"
	"github.com/axisml/axisml-lite/axisml-core/pkg/core"
)

func main() {
	fmt.Print(section("axisml-core (Lite)", &core.Config{}))
}

func section(name string, into any) string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s\n\n", name)
	b.WriteString("Configuration source: **environment only** — this binary reads no config file. Each key below is supplied as its `AXISML_` variable.\n\n")
	b.WriteString("| Key | Environment variable | Default | Secret | Description |\n")
	b.WriteString("|---|---|---|---|---|\n")
	for _, field := range configutil.Walk(into) {
		env, def, secret := "`"+field.EnvVar+"`", mdCode(field.Default), "—"
		if field.Secret {
			env, def, secret = "`"+field.EnvVar+"`<br>`"+field.EnvVar+"_FILE`", "—", "yes"
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s | %s |\n", field.Path, env, def, secret, strings.ReplaceAll(field.Doc, "|", "\\|"))
	}
	b.WriteString("\n")
	return b.String()
}

func mdCode(value string) string {
	if value == "" {
		return "—"
	}
	return "`" + value + "`"
}
