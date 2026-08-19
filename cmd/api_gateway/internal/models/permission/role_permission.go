package permission

var RolePermission = map[string]map[string]struct{}{
	Staff: {
		UserView:              {},
		UserList:              {},
		UserCreate:            {},
		UserUpdate:            {},
		UserDelete:            {},
		SubjectCreate:         {},
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
	},

	Parent: {
		UserViewMe:            {},
		ScheduleClassroomView: {},
		JournalStudentViewMe:  {},
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
	},

	User: {},
}
