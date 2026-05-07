package tests

import (
	"backend/migrations"
	"backend/pkg/api"
	"backend/pkg/postgres"
	"fmt"
	"os"
	"testing"
)

var (
	db     *postgres.DB
	server *api.Server
)

func TestMain(m *testing.M) {
	var code int
	defer func() { os.Exit(code) }()

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("TEST_POSTGRES_HOST"),
		os.Getenv("TEST_POSTGRES_PORT"),
		os.Getenv("TEST_POSTGRES_DB"),
		os.Getenv("TEST_POSTGRES_USER"),
		os.Getenv("TEST_POSTGRES_PASSWORD"),
	)

	var err error
	db, err = postgres.NewDB(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect test db: %v\n", err)
		code = 1
		return
	}
	defer db.Close()

	if err := db.Migrate(migrations.FS); err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate test db: %v\n", err)
		code = 1
		return
	}

	server = api.NewServer([]string{"*"})

	// TODO: ドメイン追加時に Module の DI・ルーティング登録を行う
	// 例:
	// tx := postgres.NewTransactor(db)
	// customerMod := customer.NewModule(db, tx)
	// customerMod.RegisterRoutes(server.Group("/customers"))

	code = m.Run()
}
