package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ThanhNV121097/project-6c2383e6/backend/migrations"
	"github.com/jackc/pgx/v5/stdlib"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := migrate(ctx, db); err != nil {
		log.Fatal(err)
	}
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name text PRIMARY KEY, checksum text NOT NULL, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}
	entries, err := fs.ReadDir(migrations.Files, ".")
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		body, err := migrations.Files.ReadFile(name)
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(body))
		var stored string
		err = db.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE name = $1`, name).Scan(&stored)
		if err == nil {
			if stored != checksum {
				return fmt.Errorf("migration checksum mismatch: %s", name)
			}
			continue
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, string(body)); err == nil {
			_, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations (name, checksum) VALUES ($1, $2)`, name, checksum)
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func init() { stdlib.GetDefaultDriver() }
