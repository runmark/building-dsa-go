package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var total float32
	for _, g := range s.Grades {
		total += g.Score
	}
	return total / float32(len(s.Grades))
}

var (
	students   Students
	gradeMutex sync.Mutex
)

type Students []Student

func (ss Students) GetByID(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}

	return nil, fmt.Errorf("Student %v not found", id)
}

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}

const (
	GradeTest     = GradeType("Test")
	GradeHomework = GradeType("Homework")
	GradeQuiz     = GradeType("Quiz")
)

type GradeType string
