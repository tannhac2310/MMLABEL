package repository

type SortOrder string

const (
	SortOrderDESC SortOrder = "DESC"
	SortOrderASC  SortOrder = "ASC"
)

var CourseTypeName = map[string]string{
	"DESC": "DESC",
	"ASC":  "ASC",
}

type Sort struct {
	Order SortOrder
	By    string
}

type CountResult struct {
	Count int64 `db:"cnt"`
}
