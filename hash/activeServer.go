package hash

type ActiveServer struct {
	active map[int][2]int
}

func NewActiveServer() ActiveServer {
	return ActiveServer{
		active: make(map[int][2]int),
	}
}

func (o *ActiveServer) IsOccupied(val int) bool {
	if _, ok := o.active[val]; ok {
		return true
	}

	return false
}

func (o *ActiveServer) Add(token int, points [2]int) {
	o.active[token] = points
}

func (o *ActiveServer) Remove(token int) {
	delete(o.active, token)
}

func (o *ActiveServer) Update(token int, points [2]int) {
	o.active[token] = points
}
