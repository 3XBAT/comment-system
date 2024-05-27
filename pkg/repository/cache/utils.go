package cache

import "github.com/3XBAT/coments-system/graph/model"

func constainsKeyPost(id string, cache map[string]model.Post) bool {
	_, exists := cache[id]
	return exists
}

func constainsKeyComment(id string, cache map[string]model.Comment) bool {
	_, exists := cache[id]
	return exists
}

