package permission

var RolePermission = map[string]map[string]struct{}{
	Staff: {
		UserView:              {},
		UserList:              {},
		UserCreate:            {},
		UserUpdate:            {},
		UserDelete:            {},
		SubjectCreate:         {},
		SubjectUpdate:         {},
		SubjectList:           {},
		ClassroomView:         {},
		ClassroomList:         {},
		ClassroomCreate:       {},
		ClassroomUpdate:       {},
		ClassroomDelete:       {},
		ScheduleClassroomView: {},
		ScheduleTeacherView:   {},
		ScheduleCreate:        {},
		ScheduleUpdate:        {},
		ScheduleDelete:        {},
		JournalStudentView:    {},
		JournalClassroomView:  {},
	},

	Student: {
		UserViewMe:            {},
		ScheduleClassroomView: {},
		JournalStudentViewMe:  {},
		SubjectList:           {},
	},

	Parent: {
		UserViewMe:            {},
		ScheduleClassroomView: {},
		JournalStudentViewMe:  {},
		SubjectList:           {},
	},

	Teacher: {
		UserViewMe:              {},
		ClassroomView:           {},
		ClassroomList:           {},
		ScheduleClassroomView:   {},
		ScheduleTeacherView:     {},
		JournalGradeUpdate:      {},
		JournalHomeworkUpdate:   {},
		JournalAttendanceUpdate: {},
		SubjectList:             {},
	},

	User: {},
}
