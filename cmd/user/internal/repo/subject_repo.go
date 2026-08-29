package repo

import (
	"context"
	"fmt"

	"github.com/abozorov/school_online/cmd/user/internal/models"
)

func (r *Repo) CreateSubject(ctx context.Context, name string, description string) (int32, error) {
	var id int32
	err := r.pg.QueryRow(ctx,
		`INSERT INTO subjects (
			name, 
			description) 
		VALUES ($1, $2) RETURNING id`,
		name, description,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("user_repo.CreateSubject: %w", postgresToErrs(err))
	}
	return id, nil
}

func (r *Repo) GetSubjectById(ctx context.Context, id int32) (*models.Subject, error) {
	var subject models.Subject
	row := r.pg.QueryRow(ctx, `SELECT id, 
		name, 
		description 
	FROM subjects 
	WHERE id = $1`, id)
	if err := row.Scan(&subject.ID, &subject.Name, &subject.Description); err != nil {
		return nil, fmt.Errorf("user_repo.GetSubjectById: %w", postgresToErrs(err))
	}
	return &subject, nil
}

func (r *Repo) GetAllSubjects(ctx context.Context) ([]*models.Subject, error) {
	rows, err := r.pg.Query(ctx, `SELECT id, 
		name, 
		description 
	FROM subjects`)
	if err != nil {
		return nil, fmt.Errorf("user_repo.GetAllSubjects: %w", postgresToErrs(err))
	}
	defer rows.Close()

	result := make([]*models.Subject, 0)
	for rows.Next() {
		var subject models.Subject
		err := rows.Scan(&subject.ID, &subject.Name, &subject.Description)
		if err != nil {
			return nil, fmt.Errorf("user_repo.GetAllSubjects: %w", postgresToErrs(err))
		}
		result = append(result, &subject)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user_repo.GetAllSubjects: %w", postgresToErrs(err))
	}
	return result, nil
}

func (r *Repo) UpdateSubject(ctx context.Context, id int32, name string, description string) error {
	_, err := r.pg.Exec(ctx,
		`UPDATE subjects 
		SET name = COALESCE(NULLIF($2, ''), name), 
			description = COALESCE(NULLIF($3, ''), description) 
		WHERE id = $1`,
		id, name, description,
	)
	return fmt.Errorf("user_repo.UpdateSubject: %w", postgresToErrs(err))
}
