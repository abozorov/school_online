package permission

const (
	UserView   = "user:view"
	UserViewMe = "user:view-me"
	UserList   = "user:list"
	UserCreate = "user:create"
	UserUpdate = "user:update"
	UserDelete = "user:Delete"

	SubjectCreate = "subject:create"

	ClassroomView   = "classroom:view"
	ClassroomList   = "classroom:list"
	ClassroomCreate = "classroom:create"
	ClassroomUpdate = "classroom:update"
	ClassroomDelete = "classroom:delete"

	ScheduleClassroomView = "schedule:classroom-view"
	ScheduleTeacherView   = "schedule:teacher-view"
	ScheduleCreate        = "schedule:create"
	ScheduleUpdate        = "schedule:update"
	ScheduleDelete        = "schedule:delete"

	JournalStudentViewMe    = "journalStudent:view-me"
	JournalStudentView      = "journalStudent:view"
	JournalClassroomView    = "journalClassroom:view"
	JournalGradeUpdate      = "journalGrade:update"
	JournalAttendanceUpdate = "journalAttendance:update"
	JournalHomeworkUpdate   = "journalHomework:update"
)
