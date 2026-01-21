package db

import "embed"

// SQLRoot is the directory (within SQLFiles) that contains migration files.
const SQLRoot = "migrations"

// SQLFiles embeds all SQL migration files.
//
//go:embed migrations/*.sql
var SQLFiles embed.FS
