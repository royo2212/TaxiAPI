package postgres

import (
	"database/sql"
	"fmt"
	"taxiAPI/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
    dbCfg := cfg.Database
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Name, dbCfg.SSLMode,
    )
    connConfig, err := pgx.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse DSN: %w", err)
    }
    conn := stdlib.GetConnector(*connConfig)
    db := sql.OpenDB(conn)
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
