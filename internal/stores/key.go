package stores

import (
	"unicode/utf8"

	"github.com/patrickhuber/wrangle/internal/dataptr"
	"github.com/patrickhuber/wrangle/internal/dataptr/parse"
)

type Key struct {
	Secret *Secret
	Path   *dataptr.DataPointer
}

type Secret struct {
	Name    string
	Version Version
}

type Version struct {
	Value  string
	Latest bool
}

func Parse(str string) (*Key, error) {
	secret := &Secret{}
	// name
	// name@v1.0.0
	name, str, err := parseName(str)
	if err != nil {
		return nil, err
	}
	secret.Name = name

	if eat(str, '@') {
		str = str[1:]
		var version Version
		version, str, err = parseVersion(str)
		if err != nil {
			return nil, err
		}
		secret.Version = version
	} else {
		secret.Version = Version{Latest: true}
	}

	if !eat(str, '/') {
		return &Key{
			Secret: secret,
			Path:   &dataptr.DataPointer{},
		}, nil
	}
	str = str[1:]

	ptr, err := parse.Parse(str)
	if err != nil {
		return nil, err
	}
	return &Key{
		Secret: secret,
		Path:   ptr,
	}, nil
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}

func parseName(str string) (capture string, rest string, err error) {
	i := 0
	for {
		r, size := utf8.DecodeRuneInString(str[i:])
		if !isLetter(r) {
			break
		}
		i += size
	}
	capture = str[0:i]
	rest = str[i:]
	err = nil
	return
}

func parseVersion(str string) (version Version, rest string, err error) {
	i := 0
	for {
		if len(str[i:]) == 0 || str[i] == '/' {
			break
		}
		i++
	}
	capture := str[0:i]
	return Version{
		Value:  capture,
		Latest: capture == "",
	}, str[i:], nil
}

func eat(str string, ch rune) bool {
	if len(str) == 0 {
		return false
	}
	r, _ := utf8.DecodeRuneInString(str)
	return r == ch
}
