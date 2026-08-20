package repo

import (
	"context"

	"github.com/abozorov/school_online/cmd/raiting/internal/models"
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

func (r *Repo) GetStudentJournal(ctx context.Context, studentId int32, startDate, endDate string) ([]models.Journal, error) {
	rows, err := r.pg.Query(ctx, `SELECT id, classroom_id, subject_id, teacher_id, date, lesson_number, room, attendance, student_id, grade, homework FROM journals WHERE student_id = $1 AND date BETWEEN to_date($2,'DD.MM.YYYY') AND to_date($3,'DD.MM.YYYY') ORDER BY date`, studentId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Journal, 0)
	for rows.Next() {
		var j models.Journal
		var attendance *bool
		var grade *int32
		var homework *string
		var room *int32
		if err := rows.Scan(&j.ID, &j.ClassroomID, &j.SubjectID, &j.TeacherID, &j.Date, &j.LessonNumber, &room, &attendance, &j.StudentID, &grade, &homework); err != nil {
			return nil, err
		}
		if room != nil {
			j.Room = *room
		}
		j.Attendance = attendance
		j.Grade = grade
		j.Homework = homework
		res = append(res, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) GetClassroomJournal(ctx context.Context, classroomId int32, startDate, endDate string) ([]models.Journal, error) {
	rows, err := r.pg.Query(ctx, `SELECT id, classroom_id, subject_id, teacher_id, date, lesson_number, room, attendance, student_id, grade, homework FROM journals WHERE classroom_id = $1 AND date BETWEEN to_date($2,'DD.MM.YYYY') AND to_date($3,'DD.MM.YYYY') ORDER BY date`, classroomId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]models.Journal, 0)
	for rows.Next() {
		var j models.Journal
		var attendance *bool
		var grade *int32
		var homework *string
		var room *int32
		if err := rows.Scan(&j.ID, &j.ClassroomID, &j.SubjectID, &j.TeacherID, &j.Date, &j.LessonNumber, &room, &attendance, &j.StudentID, &grade, &homework); err != nil {
			return nil, err
		}
		if room != nil {
			j.Room = *room
		}
		j.Attendance = attendance
		j.Grade = grade
		j.Homework = homework
		res = append(res, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) UpsertJournal(ctx context.Context, j models.Journal, dateStr string) error {
	// Try update first
	tag, err := r.pg.Exec(ctx, `UPDATE journals SET attendance = $1, grade = $2, homework = $3 WHERE classroom_id = $4 AND subject_id = $5 AND date = to_date($6,'DD.MM.YYYY') AND lesson_number = $7 AND student_id = $8`, j.Attendance, j.Grade, j.Homework, j.ClassroomID, j.SubjectID, dateStr, j.LessonNumber, j.StudentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() > 0 {
		return nil
	}
	// No rows updated, insert
	var id int32
	err = r.pg.QueryRow(ctx, `INSERT INTO journals (classroom_id, subject_id, teacher_id, date, lesson_number, room, attendance, student_id, grade, homework) VALUES ($1,$2,$3,to_date($4,'DD.MM.YYYY'),$5,$6,$7,$8,$9,$10) RETURNING id`, j.ClassroomID, j.SubjectID, j.TeacherID, dateStr, j.LessonNumber, j.Room, j.Attendance, j.StudentID, j.Grade, j.Homework).Scan(&id)
	if err != nil {
		return err
	}
	j.ID = id
	return nil
}
