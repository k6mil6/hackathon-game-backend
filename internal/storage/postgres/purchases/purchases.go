package purchases

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
	"time"
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

func (s *Storage) GetAll(ctx context.Context) ([]model.Purchase, error) {
	op := "purchases.GetAll"

	log := s.log.With(slog.String("op", op))

	log.Info("getting all purchases from storage")
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

	var purchases []dbPurchase
	if err := conn.SelectContext(ctx, &purchases, "SELECT * FROM purchases"); err != nil {
		log.Error("failed to get all purchases", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("got all purchases from storage")

	var shopPurchases []model.Purchase
	for _, purchase := range purchases {
		shopPurchases = append(shopPurchases, model.Purchase(purchase))
	}

	return shopPurchases, nil
}

func (s *Storage) Add(ctx context.Context, purchase model.Purchase) (int, error) {
	op := "purchases.Add"

	log := s.log.With(slog.String("op", op))

	log.Info("adding purchase to storage")
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

	query := `INSERT INTO purchases (item_id, user_id) VALUES ($1, $2) RETURNING id`

	var id int

	err = conn.QueryRowxContext(ctx, query, purchase.ShopItemID, purchase.BuyerID).Scan(&id)
	if err != nil {
		log.Error("failed to add purchase", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("added purchase to storage")

	return id, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbPurchase struct {
	ID         int
	ShopItemID int
	BuyerID    int
	CreatedAt  time.Time
}
