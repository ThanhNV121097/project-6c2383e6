package migrations

import "embed"

// Files contains ordered SQL migration files.
//
//go:embed *.sql
var Files embed.FS
