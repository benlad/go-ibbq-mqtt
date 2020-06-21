package main

import "encoding/json"

type temperature struct {
	Temperatures []int
}

func (t *temperature) toJson() string {
	j, _ := json.Marshal(t)

	return string(j)
}
