package api

//easyjson:json
type Forum struct {
	Slug    string
	Posts   int64
	Threads int
	Title   string
	Creator string
}
