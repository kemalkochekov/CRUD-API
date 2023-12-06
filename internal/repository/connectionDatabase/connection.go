package connectionDatabase

import (
	"CRUD_Go_Backend/internal/config"
	"context"
	"database/sql"

	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBops interface {
	GetPool(_ context.Context) *pgxpool.Pool
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	ExecQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Database struct {
	cluster *pgxpool.Pool
}

func newDatabase(cluster *pgxpool.Pool) *Database {
	return &Database{cluster: cluster}
}

func (db Database) GetPool(ctx context.Context) *pgxpool.Pool {
	return db.cluster
}

func (db Database) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, db.cluster, dest, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

func (db Database) ExecQuery(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return db.cluster.Query(ctx, query, args...)
}
func (db Database) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, db.cluster, dest, query, args...)
}

func GenerateDsn(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
}
func NewDB(ctx context.Context, cfg config.DatabaseConfig) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, GenerateDsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %v", err)
	}

	return newDatabase(pool), nil
}
func MigrationUp(cfg config.DatabaseConfig) error {
	db, err := sql.Open("postgres", GenerateDsn(cfg))
	if err != nil {
		return err
	}

	if err := goose.Up(db, "./internal/repository/migrations"); err != nil {
		return fmt.Errorf("goose migration up failed: %v", err)
	}

	return db.Close()
}

func MigrationDownAndCloseSql(cfg config.DatabaseConfig) error {
	db, err := sql.Open("postgres", GenerateDsn(cfg))
	if err != nil {
		return err
	}

	if err := goose.Down(db, "./internal/repository/migrations"); err != nil {
		return fmt.Errorf("goose migration down failed: %v", err)
	}

	return db.Close()
}
