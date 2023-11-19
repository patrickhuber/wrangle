package stores

import (
	"fmt"
	"strings"
)

type Key struct {
	Store  string
	Secret *Secret
	Path   []PathItem
}

type Secret struct {
	Name    string
	Value   string
	Version string
}

type PathItem interface {
	_pathItem()
}

type Name struct {
	Value string
	PathItem
}

type Criteria struct {
	Key   string
	Value string
	PathItem
}

// Parse parses the given key into a key object
//
//		key
//		   = store
//		   | store '/' secret_path
//
//		store
//		    = name
//
//		secret_path
//			= secret
//			| secret '/' path
//
//		secret
//			= name
//			| name '@' version
//
//		path
//			= path_item
//			|  path_item '/' path
//
//	 path_item
//			= name
//			| criteria
func Parse(key string) (*Key, error) {
	split := strings.Split(key, "/")

	if len(split) == 0 {
		return nil, fmt.Errorf("error parsing key. expected at least one segment")
	}

	// parse store
	store := split[0]

	// parse secret identifier
	var secret *Secret
	if len(split) > 1 {
		s := split[1]
		v := ""
		if strings.Contains(s, "@") {
			sv := strings.Split(split[1], "@")
			s = sv[0]
			v = sv[1]
		}
		secret = &Secret{
			Name:    s,
			Version: v,
		}
	}
	// parse the path items
	// Name = {name}
	// Criteria = {name}={value}
	var items []PathItem
	for i := 2; i < len(split); i++ {
		s := split[i]
		kv := strings.Split(s, "=")
		if len(kv) == 1 {
			items = append(items, Name{Value: kv[0]})
		} else if len(kv) == 2 {
			items = append(items, Criteria{ Key: kv[0], Value: kv[1]})
		}
	}
	return &Key{
		Store:  store,
		Secret: secret,
		Path:   items,
	}, nil
}
