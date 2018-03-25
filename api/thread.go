package api

//easyjson:json
type Thread struct {
	ID      int
	Slug    string
	Author  string
	Created string
	Forum   string
	Message string
	Title   string
	Votes   int
}

//easyjson:json
type ThreadUpdate struct {
	Message string
	Title   string
}

//easyjson:json
type Threads []Thread
