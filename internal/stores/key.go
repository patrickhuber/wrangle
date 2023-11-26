package stores

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/patrickhuber/wrangle/internal/dataptr"
	"github.com/patrickhuber/wrangle/internal/dataptr/parse"
)

type Value struct {
	Secret *Secret
	Path   *dataptr.DataPointer
}

type Secret struct {
	Name    string
	Version Version
}

type Version struct {
	Major    int
	Minor    int
	Revision int
	Latest   bool
}

func Parse(str string) (*Value, error) {
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
		return &Value{
			Secret: secret,
			Path:   &dataptr.DataPointer{},
		}, nil
	}
	str = str[1:]

	ptr, err := parse.Parse(str)
	if err != nil {
		return nil, err
	}
	return &Value{
		Secret: secret,
		Path:   ptr,
	}, nil
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}
func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
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
	if eat(str, 'v') {
		str = str[1:]
	}

	major, str, err := parseInteger(str)
	if err != nil {
		return
	}

	if !eat(str, '.') {
		err = fmt.Errorf("unable to parse version")
		return
	}
	str = str[1:]

	minor, str, err := parseInteger(str)
	if err != nil {
		return
	}

	if !eat(str, '.') {
		err = fmt.Errorf("unable to parse version")
		return
	}
	str = str[1:]

	revision, str, err := parseInteger(str)
	if err != nil {
		return
	}

	return Version{
		Major:    major,
		Minor:    minor,
		Revision: revision,
	}, str, nil
}

func parseInteger(str string) (int, string, error) {
	i := 0
	for {
		r, size := utf8.DecodeRuneInString(str[i:])
		if !isNumber(r) {
			break
		}
		i += size
	}
	integer, err := strconv.Atoi(str[0:i])
	if err != nil {
		return 0, "", err
	}
	return integer, str[i:], nil
}

func eat(str string, ch rune) bool {
	if len(str) == 0 {
		return false
	}
	r, _ := utf8.DecodeRuneInString(str)
	return r == ch
}
