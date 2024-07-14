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
					Score: 100,
					Type: GradeExam,
				},
				{
					Title: "English",
					Score: 100,
					Type: GradeExam,
				},
			},
		},
		{
			FirstName: "Tom",
			LastName: "Bond",
			ID: 2,
			Grades: []Grade{
				{
					Title: "Math",
					Score: 90,
					Type: GradeExam,
				},
				{
					Title: "English",
					Score: 70,
					Type:GradeExam,
				},
			},
		},
	}
}