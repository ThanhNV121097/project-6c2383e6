package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ThanhNV121097/project-6c2383e6/backend/migrations"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
)

const defaultGreeting = "Hello, World!"

type greetingResponse struct {
	Text string `json:"text"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := db.PingContext(r.Context()); err != nil {
			writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Service unavailable.")
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/v1/greeting", greetingHandler(db))
	log.Fatal(http.ListenAndServe(":"+port, cors(mux)))
}

func greetingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			text, err := readGreeting(r.Context(), db)
			if err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, greetingResponse{Text: text})
		case http.MethodPut:
			text, ok := decodeGreetingRequest(w, r)
			if !ok {
				return
			}
			saved, err := saveGreeting(r.Context(), db, text)
			if err != nil {
				writeDBError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, greetingResponse{Text: saved})
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func readGreeting(ctx context.Context, db *sql.DB) (string, error) {
	var text string
	err := db.QueryRowContext(ctx, `SELECT text FROM greetings WHERE id = true`).Scan(&text)
	if errors.Is(err, sql.ErrNoRows) {
		return defaultGreeting, nil
	}
	return text, err
}

func saveGreeting(ctx context.Context, db *sql.DB, text string) (string, error) {
	_, err := db.ExecContext(ctx, `INSERT INTO greetings (id, text) VALUES (true, $1) ON CONFLICT (id) DO UPDATE SET text = EXCLUDED.text`, text)
	if err != nil {
		return "", err
	}
	return text, nil
}

func decodeGreetingRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	defer r.Body.Close()
	var body map[string]json.RawMessage
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := decoder.Decode(&body); err != nil {
		writeAPIError(w, http.StatusBadRequest, "MALFORMED_REQUEST", "Request body is malformed.")
		return "", false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeAPIError(w, http.StatusBadRequest, "MALFORMED_REQUEST", "Request body is malformed.")
		return "", false
	}
	textJSON, ok := body["text"]
	if !ok || len(body) != 1 {
		writeAPIError(w, http.StatusBadRequest, "MALFORMED_REQUEST", "Request body is malformed.")
		return "", false
	}
	var rawText string
	if err := json.Unmarshal(textJSON, &rawText); err != nil {
		writeAPIError(w, http.StatusBadRequest, "MALFORMED_REQUEST", "Request body is malformed.")
		return "", false
	}
	text := strings.TrimSpace(rawText)
	if text == "" {
		writeAPIError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", "Greeting must not be empty.")
		return "", false
	}
	return text, true
}

func writeDBError(w http.ResponseWriter, err error) {
	if isUnavailable(err) {
		writeAPIError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "Service unavailable.")
		return
	}
	writeAPIError(w, http.StatusInternalServerError, "INTERNAL", "Internal server error.")
}

func isUnavailable(err error) bool {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "08006"
}

func writeAPIError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
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
