package hash

type Server struct {
	name   string
	points [2]int
	token  int
	data   map[int]interface{}
}

func NewServer(name string, points [2]int, token int) Server {
	return Server{
		name:   name,
		points: points,
		token:  token,
		data:   make(map[int]interface{}),
	}
}
