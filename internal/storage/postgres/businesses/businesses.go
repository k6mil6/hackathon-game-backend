package businesses

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
)

const (
	FarmTypeID      = 1
	FactoryTypeID   = 2
	WarehouseTypeID = 3
	BankTypeID      = 4
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

func (s *Storage) Save(ctx context.Context, business *model.Business) (int, error) {
	op := "businesses.Save"

	log := s.log.With("op", op)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return 0, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `INSERT INTO businesses (name, type_id) VALUES ($1, $2) RETURNING id`

	var id int

	if err := conn.QueryRowContext(
		ctx,
		query,
		business.Name,
		business.TypeID,
	).Scan(&id); err != nil {
		log.Error("failed to save business", slog.String("error", err.Error()))
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetAll(ctx context.Context) ([]model.Business, error) {
	op := "businesses.GetAll"

	log := s.log.With("op", op)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return nil, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `SELECT id, name, type_id, price, owner_id FROM businesses`

	var dbBusinesses []dbBusiness

	if err := conn.SelectContext(ctx, &dbBusinesses, query); err != nil {
		log.Error("failed to get all businesses", slog.String("error", err.Error()))
		return nil, err
	}

	var businesses []model.Business

	for _, dbBusiness := range dbBusinesses {
		if !dbBusiness.OwnerID.Valid {
			dbBusiness.OwnerID = sql.NullInt64{Int64: 0, Valid: true}
		}

		businesses = append(businesses, model.Business{
			ID:      dbBusiness.ID,
			Name:    dbBusiness.Name,
			TypeID:  dbBusiness.TypeID,
			Price:   dbBusiness.Price,
			OwnerID: int(dbBusiness.OwnerID.Int64),
		})
	}

	return businesses, nil
}

type dbBusiness struct {
	ID      int           `db:"id"`
	Name    string        `db:"name"`
	TypeID  int           `db:"type_id"`
	Price   float64       `db:"price"`
	OwnerID sql.NullInt64 `db:"owner_id"`
}

type dbBusinessType struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Profit      int    `db:"profit"`
}
