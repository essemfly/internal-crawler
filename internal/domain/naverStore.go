package domain

type Store struct {
	Id           int
	NaverUrl     string
	NaverHashId  string
	ReferenceUrl string
	ReferenceId  string
	IsProcessed  bool
	Name         *string
	Thumbnail    *string
}
