package server

type PostParams struct {
	Action   string
	Property Property
}

type Property struct {
	Name string
}
