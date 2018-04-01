package config

//easyjson:json
type Database struct {
	SchemaFile string
	Timestamp  string
}

//easyjson:json
type Config struct {
	Database Database
}
