package hash

import (
	"fmt"
	"math"
)

type HashRing struct {
	servers []Server
}

var activeServer = NewActiveServer()

func NewHashRing() *HashRing {
	return &HashRing{
		servers: make([]Server, 0),
	}
}

func (hr *HashRing) AddServer(name string) error {
	token := GetToken(name)

	if activeServer.IsOccupied(token) || hr.isAlreadyAdded(name) {
		return fmt.Errorf("given server is already added")
	}

	points := GetPointInHashRing(token)

	server := NewServer(name, points, token)

	activeServer.Add(server.token, server.points)

	hr.servers = append(hr.servers, server)

	hr.updateServerRangeWhenServerAdded(server)

	fmt.Println(hr.servers)

	return nil
}

func (hr *HashRing) RemoveServer(name string) error {
	idx := -1

	for i, server := range hr.servers {
		if server.name == name {
			idx = i
		}
	}

	if idx < 0 {
		return fmt.Errorf("given server not found")
	}

	hr.updateServerRangeWhenServerRemoved(hr.servers[idx])

	activeServer.Remove(hr.servers[idx].token)

	hr.servers = append(hr.servers[:idx], hr.servers[idx+1:]...)

	fmt.Println(hr.servers)

	return nil
}

func (hr *HashRing) GetValue(key string) (interface{}, error) {
	token := GetToken(key)

	serverName := hr.getServer(token)

	for _, server := range hr.servers {
		if server.name == serverName {
			if val, ok := server.data[token]; ok {
				return val, nil
			}
		}
	}

	return nil, fmt.Errorf("data not found")
}

func (hr *HashRing) AddData(key string, value interface{}) bool {

	token := GetToken(key)

	serverName := hr.getServer(token)

	for _, server := range hr.servers {
		if server.name == serverName {
			server.data[token] = value
		}
	}

	fmt.Println(hr.servers)

	return true
}

func (hr *HashRing) updateServerRangeWhenServerAdded(newServer Server) {
	for i, server := range hr.servers {
		if server.points[1] == newServer.points[1] && server.name != newServer.name {
			hr.MoveDataWhenServerAdded(server, newServer)
			server.points[1] = newServer.points[0] - 1
			hr.servers[i] = server

			activeServer.Update(server.token, server.points)
		}
	}
}

func (hr *HashRing) updateServerRangeWhenServerRemoved(removedServer Server) {
	for i, server := range hr.servers {
		if server.points[1] == removedServer.points[0]-1 {
			hr.MoveDataWhenServerRemoved(removedServer, server)
			server.points[1] = removedServer.points[1]
			hr.servers[i] = server

			activeServer.Update(server.token, server.points)
		}
	}
}

func (hr *HashRing) isAlreadyAdded(name string) bool {
	for _, server := range hr.servers {
		if server.name == name {
			return true
		}
	}

	return false
}

func (hr *HashRing) getServer(token int) string {

	maxToken := math.MinInt
	maxServerName := ""

	for _, server := range hr.servers {
		if server.points[0] <= token && server.points[1] >= token {
			return server.name
		}

		if server.token > maxToken {
			maxToken = server.token
			maxServerName = server.name
		}
	}

	return maxServerName
}

func (hr *HashRing) MoveDataWhenServerAdded(from, to Server) {

	if to.points[1] > to.points[0] {
		for token, val := range from.data {
			if to.points[0] <= token && to.points[1] >= token {
				to.data[token] = val
				delete(from.data, token)
			}
		}

		return
	}

	for token, val := range from.data {
		if to.token <= token || token < from.points[0] {
			to.data[token] = val
			delete(from.data, token)
		}
	}

}

func (hr *HashRing) MoveDataWhenServerRemoved(from, to Server) {
	for token, val := range from.data {
		to.data[token] = val
	}
}
