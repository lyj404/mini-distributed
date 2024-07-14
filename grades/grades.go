package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID int
	FirstName string
	LastName string
	Grades []Grade
}

func (s Student) average() float32{
	var result float32
	for _, grade := range s.Grades {
		result += grade.Srocde
	}

	return result / float32(len(s.Grades))
}

type Students []Student

var (
	students Students
	studentMutex sync.Mutex
)

func (ss Students) GetByID(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}

	return nil, fmt.Errorf("student with ID %d not found", id)
}

type GradeType string

const (
	GradeQuiz = GradeType("quiz")
	GradeTest = GradeType("test")
	GradeExam = GradeType("exam")
)

type Grade struct {
	Title string
	Type GradeType
	Srocde float32
}
