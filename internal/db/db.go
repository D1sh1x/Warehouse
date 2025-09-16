package db

import (
	"context"
	"database/sql"
	"time"
	"warehouse/internal/models"

	_ "github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

var r = retry.Strategy{Attempts: 3, Delay: 300 * time.Millisecond, Backoff: 2}

type DB struct {
	DB *dbpg.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := dbpg.New(dsn, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func (d *DB) CreateItem(ctx context.Context, item *models.Item) error {
	_, err := d.DB.ExecWithRetry(ctx, r, "INSERT INTO items (name, count) VALUES ($1, $2) RETURNING id", item.Name, item.Count)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetItems(ctx context.Context) ([]models.Item, error) {
	rows, err := d.DB.QueryWithRetry(ctx, r, "SELECT id, name, count FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (d *DB) UpdateItem(ctx context.Context, id int, item *models.Item) error {
	_, err := d.DB.ExecWithRetry(ctx, r, "UPDATE items SET name=$1, count=$2 WHERE id=$3", item.Name, item.Count, id)
	return err
}

func (d *DB) DeleteItem(ctx context.Context, id int) error {
	_, err := d.DB.ExecWithRetry(ctx, r, "DELETE FROM items WHERE id=$1", id)
	return err
}

func (d *DB) GetHistory(ctx context.Context) ([]models.History, error) {
	rows, err := d.DB.QueryWithRetry(ctx, r, "SELECT id, item_id, action, changed_by, timestamp, old_data, new_data FROM history ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.History
	for rows.Next() {
		var h models.History
		if err := rows.Scan(&h.ID, &h.ItemID, &h.Action, &h.ChangedBy, &h.Timestamp, &h.OldData, &h.NewData); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

func (d *DB) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	rows, err := d.DB.QueryWithRetry(ctx, r,
		"SELECT id, username, role FROM users WHERE username=$1 AND password=$2 LIMIT 1",
		username, password,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var u models.User
	if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
		return nil, err
	}
	return &u, rows.Err()
}
