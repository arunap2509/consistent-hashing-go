package hash

import (
	"consistent-hashing/config"
	"hash/fnv"
	"math"
)

func GetToken(server string) int {

	data := []byte(server)

	modulus := config.HashRingSize

	h := fnv.New32()

	h.Write(data)

	hashValue := h.Sum32()

	result := int(hashValue) % modulus

	return result
}

func GetPointInHashRing(newToken int) [2]int {
	hashSpace := config.HashRingSize

	if len(activeServer.active) == 0 {
		if newToken != 0 {
			return [2]int{newToken, newToken - 1}
		} else {
			return [2]int{newToken, hashSpace}
		}
	}

	var nextServerToken = math.MaxInt

	for token, server := range activeServer.active {
		if token > newToken && token < nextServerToken {
			nextServerToken = server[0]
		}
	}

	if nextServerToken == math.MaxInt {
		for token := range activeServer.active {
			if token < nextServerToken {
				nextServerToken = token
			}
		}
	}

	return [2]int{newToken, nextServerToken - 1}
}
