package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/juancortelezzi/gogsd/pkg/gsdlogger"
)

//go:embed schema.sql
var schemaString string

func Connect(ctx context.Context, logger gsdlogger.Logger, dsl string) (*Queries, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if result, err := db.ExecContext(ctx, schemaString); err != nil {
		formattedError := fmt.Errorf("error running migration: result=%v err=%w", result, err)
		return nil, formattedError
	}

	return New(db), nil
}
