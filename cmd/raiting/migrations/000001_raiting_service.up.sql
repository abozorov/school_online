CREATE TABLE IF NOT EXISTS "journals" (
  "id" SERIAL NOT NULL,
  "class_id" INT,
  "subject_id" INT,
  "teacher_id" INT,
  "date" DATETIME,
  "lesson_number" INT,
  "room" INT,
  "attendence" BOOLEAN,
  "student_id" INT,
  "grade" INT,
  "homework" TEXT,
  PRIMARY KEY ("id")
);
