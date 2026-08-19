package repo

import (
	"context"

	"github.com/abozorov/school_online/cmd/user/internal/models"
	"github.com/abozorov/school_online/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}

func (r *Repo) Get(ctx context.Context, id int32) (*models.User, error) {
	return r.getUserByQuery(ctx, `SELECT id, name, username, email, verify_email, password_hash, phone_number, role, birthday FROM users WHERE id = $1`, id)
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getUserByQuery(ctx, `SELECT id, name, username, email, verify_email, password_hash, phone_number, role, birthday FROM users WHERE email = $1`, email)
}

func (r *Repo) GetAll(ctx context.Context) ([]*models.User, error) {
	rows, err := r.pg.Query(ctx, `SELECT id FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.User, 0)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		user, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repo) Create(ctx context.Context, user *models.User) (int32, error) {
	if user == nil {
		return 0, models.ErrInvalidID
	}

	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var id int32
	err = tx.QueryRow(ctx, `INSERT INTO users (name, username, email, verify_email, phone_number, password_hash, refresh_token, role, birthday) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		user.Name,
		user.Username,
		user.Email,
		user.VerifyEmail,
		user.PhoneNumber,
		user.PasswordHash,
		user.RefreshToken,
		user.Role,
		user.Birthday,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	if user.StudentRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO students (classroom_id, user_id) VALUES ($1, $2)`, user.StudentRole.ClassroomID, id)
		if err != nil {
			return 0, err
		}
	}
	if user.ParentRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO parents (students_id, user_id) VALUES ($1, $2)`, user.ParentRole.StudentsID, id)
		if err != nil {
			return 0, err
		}
	}
	if user.StaffRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO staffs (position, user_id, experience) VALUES ($1, $2, $3)`, user.StaffRole.Position, id, user.StaffRole.Experience)
		if err != nil {
			return 0, err
		}
	}
	if user.TeacherRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO teachers (subjects_id, user_id, experience) VALUES ($1, $2, $3)`, user.TeacherRole.SubjectsID, id, user.TeacherRole.Experience)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) Update(ctx context.Context, user *models.User) error {
	if user == nil {
		return models.ErrInvalidID
	}

	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE users SET name = $2, username = $3, phone_number = $4, role = $5, birthday = $6 WHERE id = $1`,
		user.ID,
		user.Name,
		user.Username,
		user.PhoneNumber,
		user.Role,
		user.Birthday,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM students WHERE user_id = $1`, user.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM parents WHERE user_id = $1`, user.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM staffs WHERE user_id = $1`, user.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM teachers WHERE user_id = $1`, user.ID)
	if err != nil {
		return err
	}

	if user.StudentRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO students (classroom_id, user_id) VALUES ($1, $2)`, user.StudentRole.ClassroomID, user.ID)
		if err != nil {
			return err
		}
	}
	if user.ParentRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO parents (students_id, user_id) VALUES ($1, $2)`, user.ParentRole.StudentsID, user.ID)
		if err != nil {
			return err
		}
	}
	if user.StaffRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO staffs (position, user_id, experience) VALUES ($1, $2, $3)`, user.StaffRole.Position, user.ID, user.StaffRole.Experience)
		if err != nil {
			return err
		}
	}
	if user.TeacherRole != nil {
		_, err = tx.Exec(ctx, `INSERT INTO teachers (subjects_id, user_id, experience) VALUES ($1, $2, $3)`, user.TeacherRole.SubjectsID, user.ID, user.TeacherRole.Experience)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repo) Delete(ctx context.Context, id int32) error {
	if id <= 0 {
		return models.ErrInvalidID
	}

	tx, err := r.pg.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM students WHERE user_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM parents WHERE user_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM staffs WHERE user_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM teachers WHERE user_id = $1`, id); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) CreateSubject(ctx context.Context, name string, description string) (int32, error) {
	var id int32
	err := r.pg.QueryRow(ctx, `INSERT INTO subjects (name, description) VALUES ($1, $2) RETURNING id`, name, description).Scan(&id)
	return id, err
}

func (r *Repo) UpdateSubject(ctx context.Context, id int32, name string, description string) error {
	if name == "" && description == "" {
		return nil
	}
	_, err := r.pg.Exec(ctx, `UPDATE subjects SET name = COALESCE(NULLIF($2, ''), name), description = COALESCE(NULLIF($3, ''), description) WHERE id = $1`, id, name, description)
	return err
}

func (r *Repo) getUserByQuery(ctx context.Context, query string, arg any) (*models.User, error) {
	var user models.User
	var passwordHash string
	row := r.pg.QueryRow(ctx, query, arg)
	if err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.VerifyEmail, &passwordHash, &user.PhoneNumber, &user.Role, &user.Birthday); err != nil {
		return nil, err
	}
	user.PasswordHash = passwordHash

	studentRole, err := r.getStudentRole(ctx, user.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	user.StudentRole = studentRole

	parentRole, err := r.getParentRole(ctx, user.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	user.ParentRole = parentRole

	staffRole, err := r.getStaffRole(ctx, user.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	user.StaffRole = staffRole

	teacherRole, err := r.getTeacherRole(ctx, user.ID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	user.TeacherRole = teacherRole

	return &user, nil
}

func (r *Repo) getStudentRole(ctx context.Context, userID int32) (*models.StudentRole, error) {
	var role models.StudentRole
	err := r.pg.QueryRow(ctx, `SELECT classroom_id FROM students WHERE user_id = $1 LIMIT 1`, userID).Scan(&role.ClassroomID)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repo) getParentRole(ctx context.Context, userID int32) (*models.ParentRole, error) {
	var role models.ParentRole
	err := r.pg.QueryRow(ctx, `SELECT students_id FROM parents WHERE user_id = $1 LIMIT 1`, userID).Scan(&role.StudentsID)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repo) getStaffRole(ctx context.Context, userID int32) (*models.StaffRole, error) {
	var role models.StaffRole
	err := r.pg.QueryRow(ctx, `SELECT position, experience FROM staffs WHERE user_id = $1 LIMIT 1`, userID).Scan(&role.Position, &role.Experience)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repo) getTeacherRole(ctx context.Context, userID int32) (*models.TeacherRole, error) {
	var role models.TeacherRole
	err := r.pg.QueryRow(ctx, `SELECT subjects_id, experience FROM teachers WHERE user_id = $1 LIMIT 1`, userID).Scan(&role.SubjectsID, &role.Experience)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
