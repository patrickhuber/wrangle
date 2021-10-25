package config

import (
	"fmt"

	"github.com/imdario/mergo"
)

type multiReader struct {
	readers []Reader
}

func NewMultiReader(readers ...Reader) Reader {
	return &multiReader{
		readers: readers,
	}
}

func (r *multiReader) Get() (*Config, error) {

	if len(r.readers) == 0 {
		return nil, fmt.Errorf("must specify at least one reader")
	}
	cfg, err := r.readers[0].Get()
	if err != nil {
		return nil, err
	}
	if len(r.readers) == 1 {
		return cfg, nil
	}
	for i := 1; i < len(r.readers); i++ {

		src, err := r.readers[i].Get()
		if err != nil {
			return nil, err
		}

		err = mergo.Merge(&cfg, src)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
