CREATE TABLE IF NOT EXISTS "classroom" (
  "id" SERIAL NOT NULL,
  "grade_number" INT NOT NULL,
  "letter" VARCHAR(1) NOT NULL,
  "hometown_teacher_id" INT,
  "academic_year" VARCHAR(6) NOT NULL,
  PRIMARY KEY ("id")
);


CREATE TABLE IF NOT EXISTS "schedule" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT NOT NULL,
  "subject_id" INT NOT NULL,
  "teacher_id" INT NOT NULL,
  "day_of_week" INT NOT NULL,
  "lesson_number" INT NOT NULL,
  "room" INT,
  "academic_year" VARCHAR(6) NOT NULL,
  PRIMARY KEY ("id"),
  FOREIGN KEY ("classroom_id") REFERENCES "classroom"("id")
);
