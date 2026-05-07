package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/pkg/api"
	"backend/pkg/config"
	"backend/pkg/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	srv := api.NewServer(cfg.Origins)

	// TODO: ドメイン追加時に Module の DI・ルーティング登録を行う
	// tx := postgres.NewTransactor(db)
	// mod := xxx.NewModule(db, tx)
	// mod.RegisterRoutes(srv.Group("/api/xxx"))

	go func() {
		if err := srv.Run(cfg.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}
}
