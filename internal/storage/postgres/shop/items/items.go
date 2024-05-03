package items

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
)

type Storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewStorage(db *sqlx.DB, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) GetAll(ctx context.Context) ([]model.ShopItem, error) {
	op := "items.GetAll"

	log := s.log.With(slog.String("op", op))

	log.Info("getting all items from storage")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return nil, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	var items []dbShopItem
	if err := conn.SelectContext(ctx, &items, "SELECT * FROM shop_items"); err != nil {
		log.Error("failed to get all items", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("got all items from storage")

	var shopItems []model.ShopItem
	for _, item := range items {
		shopItems = append(shopItems, model.ShopItem(item))
	}

	return shopItems, nil
}

func (s *Storage) GetByID(ctx context.Context, id int) (model.ShopItem, error) {
	op := "items.GetByID"

	log := s.log.With(slog.String("op", op))

	log.Info("getting item from storage", slog.Int("id", id))
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return model.ShopItem{}, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	var item dbShopItem
	if err := conn.GetContext(ctx, &item, "SELECT * FROM shop_items WHERE id = $1", id); err != nil {
		log.Error("failed to get item", slog.String("error", err.Error()))
		return model.ShopItem{}, err
	}

	log.Info("got item from storage", slog.Int("id", id))

	return model.ShopItem(item), nil
}

func (s *Storage) Add(ctx context.Context, item model.ShopItem) (int, error) {
	op := "items.Add"

	log := s.log.With(slog.String("op", op))

	log.Info("adding item to storage")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return 0, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `INSERT INTO shop_items (name, description, price, in_stock) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id`

	var id int
	err = conn.QueryRowxContext(ctx, query, item.Name, item.Description, item.Price, item.InStock).Scan(&id)
	if err != nil {
		log.Error("failed to add item", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("added item to storage")

	return id, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbShopItem struct {
	ID          int     `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float64 `db:"price"`
	InStock     int     `db:"in_stock"`
}
