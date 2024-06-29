package database

import (
    "database/sql"
    "embed"
    "fmt"
    "github.com/rubenv/sql-migrate"
    "log"
)

//go:embed sql_migrations/*.sql
var migrationFiles embed.FS

var (
    DbConnection *sql.DB
)

func DbMigrate(dbParam *sql.DB) {
    migrations := &migrate.EmbedFileSystemMigrationSource{
        FileSystem: migrationFiles,
        Root:       "sql_migrations",
    }

    n, err := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
    if err != nil {
        log.Fatalf("Failed to apply migrations: %v", err)
    }

    DbConnection = dbParam

    fmt.Printf("Applied %d migrations!\n", n)
}
