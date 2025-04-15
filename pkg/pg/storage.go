package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kon3gor/joba/pkg"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	DBName   string `yaml:"db-name"`
}

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, c Config) (pkg.Storage, func(), error) {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", c.User, c.Password, c.Address, c.DBName))
	if err != nil {
		return nil, nil, err
	}

	return &Storage{
		db: pool,
	}, pool.Close, nil
}

func (s *Storage) Intersect(ctx context.Context, aid string, ids []string) ([]string, error) {
	query := `
		SELECT id FROM known_jobs
		WHERE alert_id = $1 AND id = ANY($2)
	`
	res, err := s.db.Query(ctx, query, aid, ids)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ids, nil
		}

		return nil, err
	}
	defer res.Close()

	knwonIds := make([]string, 0, len(ids))
	var id string
	for res.Next() {
		err := res.Scan(&id)
		if err != nil {
			return nil, err
		}
		knwonIds = append(knwonIds, id)
	}

	unknownIds := make([]string, 0, len(ids))
	for _, id := range ids {
		if slices.Contains(knwonIds, id) {
			continue
		}

		unknownIds = append(unknownIds, id)
	}

	return unknownIds, nil
}

func (s *Storage) Save(ctx context.Context, aid string, ids []string) error {
	// who made this????????
	query := squirrel.Insert("known_jobs").Columns("alert_id", "id").PlaceholderFormat(squirrel.Dollar)
	for _, id := range ids {
		query = query.Values(aid, id)
	}
	queryStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = s.db.Exec(ctx, queryStr, args...)
	return err
}
