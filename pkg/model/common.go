package model

import "mmlabel.gitlab.com/mm-printing-backend/pkg/enum"

type UserInfo struct {
	Name string
	ID   string
}
type CourseInfo struct {
	ID     string
	Code   string
	Title  string
	Type   enum.CourseType
	Status enum.CommonStatus
}
type SyllabusInfo struct {
	ID        string
	Title     string
	CourseID  string
	Course    *CourseInfo
	Code      string
	TeacherID string
	Teacher   *UserInfo
	Status    enum.CommonStatus
}
