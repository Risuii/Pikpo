package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pikpo2/helpers/exception"
	"pikpo2/models"
)

type (
	UserRepository interface {
		Create(ctx context.Context, user models.User) (int64, error)
		FindByID(ctx context.Context, id int64) (models.User, error)
		Update(ctx context.Context, id int64, user models.User) error
		UpdateStatus(ctx context.Context, id int64, user models.User) error
	}

	userRepositoryImpl struct {
		db              *sql.DB
		tableName       string
		secondTableName string
	}
)

func NewUserRepository(db *sql.DB, tableName string, secondTableName string) UserRepository {
	return &userRepositoryImpl{
		db:              db,
		tableName:       tableName,
		secondTableName: secondTableName,
	}
}

func (repo *userRepositoryImpl) Create(ctx context.Context, user models.User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (username_T3,role,status) VALUES (?,?,?)", repo.tableName)
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		user.UsernameT3,
		user.Role,
		user.Status,
	)
	if err != nil {
		return 0, err
	}

	ID, _ := result.LastInsertId()

	return ID, nil
}

func (repo *userRepositoryImpl) FindByID(ctx context.Context, id int64) (models.User, error) {
	var user models.User
	var tier models.Tier

	// query := fmt.Sprintf(`SELECT id, username_T3, role, status FROM %s WHERE id = ?`, repo.tableName)

	query := fmt.Sprintf("SELECT users.id, username_T3,username_T2, username_T1, role, status, tiers.id, user_id, tier FROM %s JOIN %s ON %s.user_id = %s.id WHERE %s.id = ?", repo.tableName, repo.secondTableName, repo.secondTableName, repo.tableName, repo.tableName)

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		log.Println(err)
		return user, exception.ErrInternalServer
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(
		&user.ID,
		&user.UsernameT3,
		&user.UsernameT2,
		&user.UsernameT1,
		&user.Role,
		&user.Status,
		&tier.ID,
		&tier.UserId,
		&tier.Tier,
	)
	if err != nil {
		log.Println(err)
		return user, exception.ErrNotFound
	}

	// user.ID = int(id)

	user.Tier = append(user.Tier, tier)

	return user, nil
}

func (repo *userRepositoryImpl) Update(ctx context.Context, id int64, user models.User) error {

	if user.UsernameT2 != "" {
		command := fmt.Sprintf(`UPDATE %s SET username_T2 = ? WHERE id = %d`, repo.tableName, id)
		stmt, err := repo.db.PrepareContext(ctx, command)
		if err != nil {
			log.Println(err)
			return exception.ErrInternalServer
		}
		defer stmt.Close()

		result, err := stmt.ExecContext(
			ctx,
			user.UsernameT2,
		)

		if err != nil {
			log.Println(err)
			return exception.ErrInternalServer
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected < 1 {
			return exception.ErrNotFound
		}
	}

	if user.UsernameT1 != "" {
		command := fmt.Sprintf(`UPDATE %s SET username_T1 = ? WHERE id = %d`, repo.tableName, id)
		stmt, err := repo.db.PrepareContext(ctx, command)
		if err != nil {
			log.Println(err)
			return exception.ErrInternalServer
		}
		defer stmt.Close()

		result, err := stmt.ExecContext(
			ctx,
			user.UsernameT1,
		)

		if err != nil {
			log.Println(err)
			return exception.ErrInternalServer
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected < 1 {
			return exception.ErrNotFound
		}
	}

	return nil
}

func (repo *userRepositoryImpl) UpdateStatus(ctx context.Context, id int64, user models.User) error {
	command := fmt.Sprintf(`UPDATE %s SET status = ? WHERE id = %d`, repo.tableName, id)
	stmt, err := repo.db.PrepareContext(ctx, command)
	if err != nil {
		log.Println(err)
		return exception.ErrInternalServer
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		&user.Status,
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
