package service

import (
	"log"
	"strings"

	"github.com/couatl/forum-db-api/modules/assets/assets_db"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

type ForumGeneric struct {
	db *sqlx.DB
}
type DatabaseType struct {
	Prefix      string
	Description string
	Example     string
	Factory     func(string) ForumHandler
}

var supportedDatabases = []DatabaseType{
	{"postgres", "PostgreSQL database", "postgres://docker:docker@localhost/docker", NewForumPgSQL},
}

func NewForum(database string) ForumHandler {
	help := "Supported database types:"
	for _, item := range supportedDatabases {
		var source string
		help += "\n- " + item.Prefix + "\t" + item.Description + " (example: " + item.Example + ")"
		if database == item.Prefix {
			source = ""
		} else if strings.HasPrefix(database, item.Prefix+"://") {
			source = database
		} else if strings.HasPrefix(database, item.Prefix+":") {
			source = database[len(item.Prefix)+1:]
		} else {
			continue
		}
		return item.Factory(source)
	}
	panic("Unsupported database type: " + database + "\n" + help)
}

func NewForumGeneric(dialect string, dataSourceName string) ForumGeneric {
	migrations := &migrate.AssetMigrationSource{
		Asset:    assets_db.Asset,
		AssetDir: assets_db.AssetDir,
		Dir:      "db/" + dialect,
	}
	db, err := sqlx.Open(dialect, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = migrate.Exec(db.DB, dialect, migrations, migrate.Up); err != nil {
		log.Fatal(err)
	}
	return ForumGeneric{db: db}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}
