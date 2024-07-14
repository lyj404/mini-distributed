package grades

func init() {
	students = []Student{
		{
			FirstName: "James",
			LastName: "Bond",
			ID: 1,
			Grades: []Grade{
				{
					Title: "Math",
					Srocde: 100,
					Type: GradeExam,
				},
				{
					Title: "English",
					Srocde: 100,
					Type: GradeExam,
				},
			},
		},
	}
}