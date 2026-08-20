package repo

import (
	"context"

	"github.com/abozorov/school_online/cmd/classroom/internal/models"
	"github.com/abozorov/school_online/pkg/postgres"
)

type Repo struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		pg: pg,
	}
}

func (r *Repo) Get(ctx context.Context, id int32) (*models.Classroom, error) {
	var c models.Classroom
	row := r.pg.QueryRow(ctx, `SELECT id, grade_number, letter, hometown_teacher_id, academic_year FROM classroom WHERE id = $1 LIMIT 1`, id)
	if err := row.Scan(&c.ID, &c.GradeNumber, &c.Letter, &c.HometownTeacherID, &c.AcademicYear); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) List(ctx context.Context, page, limit int32) ([]*models.Classroom, error) {
	if limit <= 0 {
		limit = 50
	}
	off := int((page - 1) * limit)
	rows, err := r.pg.Query(ctx, `SELECT id FROM classroom ORDER BY id LIMIT $1 OFFSET $2`, limit, off)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Classroom, 0)
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		c, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) Create(ctx context.Context, req models.ClassroomRequest) (int32, error) {
	var id int32
	err := r.pg.QueryRow(ctx, `INSERT INTO classroom (grade_number, letter, hometown_teacher_id, academic_year) VALUES ($1, $2, $3, $4) RETURNING id`,
		req.GradeNumber, req.Letter, req.HometownTeacherID, req.AcademicYear).Scan(&id)
	return id, err
}

func (r *Repo) Update(ctx context.Context, req models.Classroom) (int32, error) {
	_, err := r.pg.Exec(ctx, `UPDATE classroom SET grade_number = COALESCE(NULLIF($2,0), grade_number), letter = COALESCE(NULLIF($3,''), letter), hometown_teacher_id = COALESCE(NULLIF($4,0), hometown_teacher_id), academic_year = COALESCE(NULLIF($5,''), academic_year) WHERE id = $1`,
		req.ID, req.GradeNumber, req.Letter, req.HometownTeacherID, req.AcademicYear)
	if err != nil {
		return 0, err
	}
	return req.ID, nil
}

func (r *Repo) Delete(ctx context.Context, id int32) error {
	_, err := r.pg.Exec(ctx, `DELETE FROM classroom WHERE id = $1`, id)
	return err
}

// Schedule operations
func (r *Repo) GetScheduleByClassroom(ctx context.Context, classroomId int32) ([]*models.Schedule, error) {
	rows, err := r.pg.Query(ctx, `SELECT id, classroom_id, subject_id, teacher_id, day_of_week, lesson_number, room, academic_year FROM schedule WHERE classroom_id = $1 ORDER BY id`, classroomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Schedule, 0)
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.ClassroomID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek, &s.LessonNumber, &s.Room, &s.AcademicYear); err != nil {
			return nil, err
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) GetScheduleByTeacher(ctx context.Context, teacherId int32) ([]*models.Schedule, error) {
	rows, err := r.pg.Query(ctx, `SELECT id, classroom_id, subject_id, teacher_id, day_of_week, lesson_number, room, academic_year FROM schedule WHERE teacher_id = $1 ORDER BY id`, teacherId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*models.Schedule, 0)
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(&s.ID, &s.ClassroomID, &s.SubjectID, &s.TeacherID, &s.DayOfWeek, &s.LessonNumber, &s.Room, &s.AcademicYear); err != nil {
			return nil, err
		}
		res = append(res, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) CreateSchedule(ctx context.Context, req models.ScheduleRequest) (int32, error) {
	var id int32
	err := r.pg.QueryRow(ctx, `INSERT INTO schedule (classroom_id, subject_id, teacher_id, day_of_week, lesson_number, room, academic_year) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		req.ClassroomID, req.SubjectID, req.TeacherID, req.DayOfWeek, req.LessonNumber, req.Room, req.AcademicYear).Scan(&id)
	return id, err
}

func (r *Repo) UpdateSchedule(ctx context.Context, req models.Schedule) error {
	_, err := r.pg.Exec(ctx, `UPDATE schedule SET classroom_id = COALESCE(NULLIF($2,0), classroom_id), subject_id = COALESCE(NULLIF($3,0), subject_id), teacher_id = COALESCE(NULLIF($4,0), teacher_id), day_of_week = COALESCE(NULLIF($5,0), day_of_week), lesson_number = COALESCE(NULLIF($6,0), lesson_number), room = COALESCE(NULLIF($7,0), room), academic_year = COALESCE(NULLIF($8,''), academic_year) WHERE id = $1`,
		req.ID, req.ClassroomID, req.SubjectID, req.TeacherID, req.DayOfWeek, req.LessonNumber, req.Room, req.AcademicYear)
	return err
}

func (r *Repo) DeleteSchedule(ctx context.Context, id int32) error {
	_, err := r.pg.Exec(ctx, `DELETE FROM schedule WHERE id = $1`, id)
	return err
}
