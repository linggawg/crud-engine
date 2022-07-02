package commands

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type EngineCommand struct {
}

func NewEngineCommand() *EngineCommand {
	return &EngineCommand{}
}

func (s *EngineCommand) InsertOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error) {
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func (s *EngineCommand) UpdateOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error) {
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *EngineCommand) DeleteOne(ctx context.Context, db *sqlx.DB, query string, args []interface{}) (err error) {
	_, err = db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
