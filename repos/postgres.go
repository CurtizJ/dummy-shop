package repos

import (
	"database/sql"
	"log"
	"os"

	. "github.com/CurtizJ/dummy-shop/errors"
	. "github.com/CurtizJ/dummy-shop/items"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	PG_CONSTRAINT_VIOLATES_ERROR_CODE = "23505"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo() Repo {
	db, err := sql.Open("pgx", os.Getenv("PG_URL"))
	if err != nil {
		log.Panicf("Cannot startup Postgres. err: %v", err)
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS items(
			Id integer PRIMARY KEY,
			name varchar(100),
			category varchar(100))`)

	if err != nil {
		log.Panicf("Cannot startup Postgres. err: %v", err)
	}

	return &PostgresRepo{db}
}

func (repo *PostgresRepo) Get(id uint64) (*Item, error) {
	row := repo.db.QueryRow("SELECT * FROM items WHERE Id = $1", id)
	item, err := FromRow(row)

	if err == sql.ErrNoRows {
		return nil, &ItemNotFoundError{id}
	} else if err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *PostgresRepo) Add(item *Item) error {
	_, err := repo.db.Exec("INSERT INTO items VALUES ($1, $2, $3)", item.Id, item.Name, item.Category)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == PG_CONSTRAINT_VIOLATES_ERROR_CODE {
		return &ItemAlreadyExistsError{item.Id}
	}

	return err
}

func (repo *PostgresRepo) Update(newItem *Item) error {
	row := repo.db.QueryRow(
		`UPDATE items SET name = $1, category = $2 WHERE Id = $3 RETURNING Id`,
		newItem.Name, newItem.Category, newItem.Id)

	var updated uint64
	err := row.Scan(&updated)
	if err == sql.ErrNoRows {
		return &ItemNotFoundError{newItem.Id}
	}

	return err
}

func (repo *PostgresRepo) Delete(id uint64) error {
	row := repo.db.QueryRow(`DELETE FROM items WHERE Id = $1 RETURNING Id`, id)

	var deleted uint64
	err := row.Scan(&deleted)
	if err == sql.ErrNoRows {
		return &ItemNotFoundError{id}
	}

	return err
}

func (repo *PostgresRepo) ListAll() ([]Item, error) {
	rows, err := repo.db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}

	list := make([]Item, 0)
	for rows.Next() {
		item, err := FromRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (repo *PostgresRepo) List(length, offset uint64) ([]Item, error) {
	rows, err := repo.db.Query(
		`SELECT * FROM items 
			ORDER BY Id LIMIT $1 OFFSET $2`, length, offset)
	if err != nil {
		return nil, err
	}

	list := make([]Item, 0)
	for rows.Next() {
		item, err := FromRow(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}
