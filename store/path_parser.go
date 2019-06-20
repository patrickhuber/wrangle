package store

import (
	"fmt"
	"strings"
)

// StoreAndPath contains the store and path of a path string after parsing
type StoreAndPath struct {
	Store string
	Path  string
}

// ParsePath parses the path into a StoreAndPath object
func ParsePath(path string) (*StoreAndPath, error) {
	storeAndPathSplit := strings.Split(path, ":")
	if len(storeAndPathSplit) != 2 {
		return nil, fmt.Errorf("error parsing source, expected <source>:<key>, found %s", path)
	}
	return &StoreAndPath{
		Store: storeAndPathSplit[0],
		Path:  storeAndPathSplit[1],
	}, nil
}
