package postgres

import (
    "database/sql"
    "fmt"
    "taxiAPI/config"
    _ "github.com/lib/pq"
)
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
    dbCfg := cfg.Database
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Name, dbCfg.SSLMode,
    )
    
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database connection: %w", err)
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}