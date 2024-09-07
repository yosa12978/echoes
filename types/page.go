package types

type Page[T interface{}] struct {
	HasNext  bool
	Size     int
	NextPage int
	Content  []T
	Total    int
}
