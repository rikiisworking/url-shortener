package storage

import (
	"context"
	"fmt"
	"time"

	"gitgub.com/rikiisworking/url-shortener/internal/config"
	"gitgub.com/rikiisworking/url-shortener/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(cfg config.Config) (*PostgresRepo, error) {
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// Good production defaults
	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create table if not exists
	createTable := `
		CREATE TABLE IF NOT EXISTS urls (
			id BIGSERIAL PRIMARY KEY,
			original_url TEXT NOT NULL,
			short_code VARCHAR(10) UNIQUE NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			expires_at TIMESTAMPTZ,
			click_count BIGINT DEFAULT 0
		);
	`
	_, err = pool.Exec(context.Background(), createTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresRepo{pool: pool}, nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}

func (r *PostgresRepo) Create(ctx context.Context, url *model.URL) error {
	query := `
		INSERT INTO urls (original_url, short_code)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	row := r.pool.QueryRow(ctx, query, url.OriginalURL, url.ShortCode)

	return row.Scan(&url.ID, &url.CreatedAt)
}

func (r *PostgresRepo) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query := `
		SELECT id, original_url, short_code, created_at, expires_at, click_count
		FROM urls
		WHERE short_code = $1
	`
	u := &model.URL{}
	err := r.pool.QueryRow(ctx, query, shortCode).Scan(
		&u.ID, &u.OriginalURL, &u.ShortCode, &u.CreatedAt, &u.ExpiresAt, &u.ClickCount,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresRepo) IncrementClick(ctx context.Context, shortCode string) error {
	query := `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`
	_, err := r.pool.Exec(ctx, query, shortCode)
	return err
}
