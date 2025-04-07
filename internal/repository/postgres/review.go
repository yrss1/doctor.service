package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/yrss1/doctor.service/internal/domain/review"
	"github.com/yrss1/doctor.service/pkg/store"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) List(ctx context.Context) (dest []review.Entity, err error) {
	query := `
		SELECT 
			id, doctor_id, user_id, rating, comment
		FROM reviews;
		`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *ReviewRepository) Add(ctx context.Context, data review.Entity) (id string, err error) {
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM appointments
			WHERE user_id = $1 AND doctor_id = $2 AND status = 'completed'
		);
	`
	if err = r.db.GetContext(ctx, &exists, checkQuery, data.UserID, data.DoctorID); err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("user has not completed an appointment with this doctor")
	}

	// Проверка: уже оставлял отзыв
	// var alreadyLeft bool
	// checkDup := `
	// 	SELECT EXISTS (
	// 		SELECT 1 FROM reviews
	// 		WHERE user_id = $1 AND doctor_id = $2
	// 	);
	// `
	// if err = r.db.GetContext(ctx, &alreadyLeft, checkDup, data.UserID, data.DoctorID); err != nil {
	// 	return "", err
	// }
	// if alreadyLeft {
	// 	return "", errors.New("review already exists for this doctor by this user")
	// }

	query := `
		INSERT INTO reviews (
			doctor_id, user_id, rating, comment
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id;
	`

	args := []any{
		data.DoctorID,
		data.UserID,
		data.Rating,
		data.Comment,
	}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ReviewRepository) Get(ctx context.Context, id string) (dest review.Entity, err error) {
	query := `
	SELECT 
		id, doctor_id, user_id, rating, comment
	FROM reviews
	WHERE id = $1;
	`
	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ReviewRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM reviews 
	WHERE id = $1;
	RETURNING id;
	`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}
