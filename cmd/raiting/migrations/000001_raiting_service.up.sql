CREATE TABLE IF NOT EXISTS "journals" (
  "id" SERIAL NOT NULL,
  "classroom_id" INT NOT NULL,
  "subject_id" INT NOT NULL,
  "teacher_id" INT NOT NULL,
  "date" DATE NOT NULL,
  "lesson_number" INT NOT NULL,
  "room" INT,
  "attendance" BOOLEAN DEFAULT true,
  "student_id" INT NOT NULL,
  "grade" INT,
  "homework" TEXT,
  PRIMARY KEY ("id")
);