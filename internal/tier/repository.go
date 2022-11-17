package tier

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pikpo2/helpers/exception"
	"pikpo2/models"
	"strconv"
)

type (
	TierRepository interface {
		Create(ctx context.Context, tier models.Tier) (int64, error)
		FindByUserID(ctx context.Context, userID int64) (models.Tier, error)
		UpdateTier(ctx context.Context, ID int64, TierNumber int64, Tier models.Tier) error
	}

	tierRepositoryImpl struct {
		db        *sql.DB
		tableName string
	}
)

func NewTierRepository(db *sql.DB, tableName string) TierRepository {
	return &tierRepositoryImpl{
		db:        db,
		tableName: tableName,
	}
}

func (repo *tierRepositoryImpl) Create(ctx context.Context, tier models.Tier) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id,tier) VALUES (?, ?)", repo.tableName)
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		tier.UserId,
		tier.Tier,
	)
	if err != nil {
		return 0, err
	}

	ID, _ := result.LastInsertId()

	return ID, nil
}

func (repo *tierRepositoryImpl) FindByUserID(ctx context.Context, userID int64) (models.Tier, error) {
	var tier models.Tier

	query := fmt.Sprintf("SELECT user_id, tier FROM %s WHERE user_id = ?", repo.tableName)

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return tier, exception.ErrInternalServer
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID)

	err = row.Scan(
		&tier.UserId,
		&tier.Tier,
	)

	if err != nil {
		log.Println(err)
		return tier, exception.ErrNotFound
	}

	return tier, nil
}

func (repo *tierRepositoryImpl) UpdateTier(ctx context.Context, ID int64, TierNumber int64, Tier models.Tier) error {

	command := fmt.Sprintf(`UPDATE %s SET tier = ? WHERE user_id = %d`, repo.tableName, ID)
	stmt, err := repo.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	Tier.Tier = "Tier" + "-" + strconv.Itoa(int(TierNumber))

	result, err := stmt.ExecContext(
		ctx,
		&Tier.Tier,
	)

	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected < 1 {
		return exception.ErrNotFound
	}

	return nil
}
