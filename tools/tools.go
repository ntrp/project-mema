//go:build tools

package tools

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/riverqueue/river/cmd/river"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
