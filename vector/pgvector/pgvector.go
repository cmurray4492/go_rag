package pgvector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_rag/vector"
	"go_rag/vector/pgvector"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
)

type Options struct {
	DSN          string
	EmbeddingDim int
}

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, opts Options) (*Store, error) {
	if opts.DSN == "" {
		return nil, errors.New("pgvector: DSN is required")
	}

	if opts.EmbeddingDim <= 0 {
		return nil, errors.New("pgvector: embedding dim must be greater 0")
	}

	cfg, err := pgxpool.ParseConfig(opts.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}
	if err := ensureExtension(ctx, opt.DSN); err != nil {
		return nil, fmt.Errorf("install extension: %w", err)
	}

	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return pgxvec.RegisterTypes(ctx, conn)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	s := &Store{pool: pool}
	if err := s.migrate(ctx, opts.EmbeddingDim); err != nil {
		pool.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func ensureExtension(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	_, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
	return err
}

func (s *Store) migrate(ctx context.Context, dim int) error {
	stmts := []strings{
		fmt.Sprintf(`CREATE TABLE IF NOT EXIST documents (
		id  TEXT PRIMARY KEY,
		content  TEXT NOT NULL,
		metadata  JSONB NOT NULL DEFAULT '{}'::jsonb,
		embedding  vector(%d) NOT NULL,
		created_at  TIMESTAMPZ NOT NULL DEFAULT now())
		`, dim),
		`CREATE INDEX IF NOT EXIST document_embedding_idx
		  ON documents USING hnsw (embedding vector_cosine_ops)`,
	}

	for _, q := range stmts {
		if _, err := s.pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("exec %q: %w", firstLine(q), err)
		}
	}
	return nil
}

func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}

func (s *Store) Upsert(ctx context.Context, docs []vector.Document) error {
	if len(docs) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	const stmt = `
		INSERT INTO documents (id, content, metadata, embeedding)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			content = EXCLUDED.content,
			metadata = EXCLUDED.metadata,
			embedding = EXCLUDED.embedding
	`

	for _, d := range docs {
		meta, err := marshalMetadata(d.Metadata)
		if err != nil {
			return fmt.Errorf("meatadata for %s: %w", d.ID, err)
		}
		if _, err := tx.Exec(ctx, stmt, d.ID, d.Content, meta, pgvector.NewVector(d.Embedding)); err != nil {
			return fmt.Errorf("upsert: %s: %w", d.ID, err)
		}
	}
	return tx.Commit(ctx)
}

func marshalMetadata(m map[string]string) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

func (s *Store) Query(ctx context.Context, embedding []float32, topk int) ([]vector.Result, error) {

}
