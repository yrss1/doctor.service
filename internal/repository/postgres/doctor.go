package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/yrss1/doctor.service/internal/domain/doctor"
	"github.com/yrss1/doctor.service/pkg/store"

	"github.com/jmoiron/sqlx"
)

type DoctorRepository struct {
	db *sqlx.DB
}

func NewDoctorRepository(db *sqlx.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) ListWithSchedules(ctx context.Context) ([]doctor.Entity, error) {
	query := `
	SELECT
		d.id,
		d.name,
		d.specialization,
		d.experience,
		d.price,
		d.rating,
		d.address,
		d.phone,
		d.photo_url,
		c.name AS clinic_name,
		COALESCE(json_agg(
			json_build_object(
				'schedule_id', s.id,
				'slot_start', s.slot_start,
				'slot_end', s.slot_end
			)
		) FILTER (WHERE s.is_available = TRUE), '[]') AS available_schedules
	FROM doctors d
	LEFT JOIN clinic c ON d.clinic_id = c.id
	LEFT JOIN schedule s ON s.doctor_id = d.id
	GROUP BY d.id, c.name
	ORDER BY d.name;`

	var rawList []doctor.EntityWithRaw
	err := r.db.SelectContext(ctx, &rawList, query)
	if err != nil {
		return nil, err
	}

	var result []doctor.Entity
	for _, row := range rawList {
		var schedules []doctor.ScheduleSlot
		if err := json.Unmarshal(row.AvailableSchedulesRaw, &schedules); err != nil {
			return nil, err
		}
		row.Entity.AvailableSchedules = schedules
		result = append(result, row.Entity)
	}

	return result, nil
}

func (r *DoctorRepository) GetWithSchedules(ctx context.Context, id string) (doctor.Entity, error) {
	query := `
	SELECT
		d.id,
		d.name,
		d.specialization,
		d.experience,
		d.price,
		d.rating,
		d.address,
		d.phone,
		d.photo_url,
		c.name AS clinic_name,
		COALESCE(json_agg(
			json_build_object(
				'schedule_id', s.id,
				'slot_start', to_json(s.slot_start),
				'slot_end', to_json(s.slot_end)
			)
		) FILTER (WHERE s.is_available = TRUE), '[]') AS available_schedules
	FROM doctors d
	LEFT JOIN clinic c ON d.clinic_id = c.id
	LEFT JOIN schedule s ON s.doctor_id = d.id
	WHERE d.id = $1
	GROUP BY d.id, c.name
	LIMIT 1;
	`

	var raw doctor.EntityWithRaw
	err := r.db.GetContext(ctx, &raw, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return doctor.Entity{}, store.ErrorNotFound
		}
		return doctor.Entity{}, err
	}

	var schedules []doctor.ScheduleSlot
	if err := json.Unmarshal(raw.AvailableSchedulesRaw, &schedules); err != nil {
		return doctor.Entity{}, err
	}

	raw.Entity.AvailableSchedules = schedules

	return raw.Entity, nil
}

func (r *DoctorRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
	DELETE FROM doctors 
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

func (r *DoctorRepository) SearchWithSchedules(ctx context.Context, filter doctor.Entity) ([]doctor.Entity, error) {
	query := `
	SELECT
		d.id,
		d.name,
		d.specialization,
		d.experience,
		d.price,
		d.rating,
		d.address,
		d.phone,
		d.photo_url,
		c.name AS clinic_name,
		COALESCE(json_agg(
			json_build_object(
				'schedule_id', s.id,
				'slot_start', to_json(s.slot_start),
				'slot_end', to_json(s.slot_end)
			)
		) FILTER (WHERE s.is_available = TRUE), '[]') AS available_schedules
	FROM doctors d
	LEFT JOIN clinic c ON d.clinic_id = c.id
	LEFT JOIN schedule s ON s.doctor_id = d.id
	WHERE 1=1
	`

	sets, args := r.prepareSearchArgs(filter)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}
	query += " GROUP BY d.id, c.name ORDER BY d.name;"

	var rawList []doctor.EntityWithRaw
	err := r.db.SelectContext(ctx, &rawList, query, args...)
	if err != nil {
		return nil, err
	}

	var result []doctor.Entity
	for _, row := range rawList {
		var schedules []doctor.ScheduleSlot
		if err := json.Unmarshal(row.AvailableSchedulesRaw, &schedules); err != nil {
			return nil, err
		}
		row.Entity.AvailableSchedules = schedules
		result = append(result, row.Entity)
	}

	return result, nil
}

func (r *DoctorRepository) prepareSearchArgs(data doctor.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, "%"+*data.Name+"%")
		sets = append(sets, fmt.Sprintf("d.name ILIKE $%d", len(args)))
	}

	if data.Specialization != nil {
		args = append(args, "%"+*data.Specialization+"%")
		sets = append(sets, fmt.Sprintf("d.specialization ILIKE $%d", len(args)))
	}

	if data.ClinicName != nil {
		args = append(args, "%"+*data.ClinicName+"%")
		sets = append(sets, fmt.Sprintf("c.name ILIKE $%d", len(args)))
	}

	return
}
