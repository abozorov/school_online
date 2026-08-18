CREATE TABLE IF NOT EXISTS "classroom" (
  "id" SERIAL NOT NULL,
  "grade_number" INT,
  "letter" VARCHAR(1),
  "hometown_teacher_id" INT,
  "academic_year" YEAR,
  PRIMARY KEY ("id")
);


CREATE TABLE IF NOT EXISTS "schedule" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT,
  "subject_id" INT,
  "teacher_id" INT,
  "day_of_week" INT,
  "lesson_number" INT,
  "room" INT,
  "academic_yeart" INT,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("classroom_id") REFERENCES "classroom"("id") ON DELETE CASCADE ON UPDATE CASCADE
);
