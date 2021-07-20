package ggo

import (
	"errors"
	"os"
	"strings"
)

type ConfigEntry struct {
	IsActive bool
	name     string
	Value    string
	Comment  string
}

var (
	CommentOnly = errors.New("comment only value")
)

type ConfigEntryInterface interface {
	Name() string
	String() string
	StringLn() string

	ChooseActiveOrReduce(v *ConfigEntry) ConfigEntryInterface
	Copy() ConfigEntryInterface
}

func (e *ConfigEntry) Name() string {
	return e.name
}

func removeEmpty(strs []string) []string {
	result := make([]string, 0, len(strs))

	for i := 0; i < len(strs); i++ {
		str := strs[i]
		if str != "" {
			result = append(result, []string{str}...)
		}
	}
	return result
}

func getValue(line []string) (*string, []string, error) {
	if line[0][0] == '#' {
		return nil, line, nil
	}

	if line[0][0] != '"' {
		// Comment
		if len(line) > 1 && line[1][0] != '#' {
			return nil, nil, CommentOnly
		}
		return &line[0], line[1:], nil
	}

	for i, v := range line {
		l := len(v)
		if l > 2 && v[l-2:] == "\\\"" {
			continue
		}

		if l > 0 && v[l-1] == '"' {
			v := strings.Join(line[:i+1], " ")
			return &v, line[i+1:], nil
		}
	}
	return nil, nil, os.ErrInvalid
}

func ParseString(str string) *ConfigEntry {
	e := new(ConfigEntry)

	strTrimmed := strings.TrimSpace(strings.ReplaceAll(str, "\t", " "))

	if strTrimmed == "#" || strTrimmed == "" {
		return nil
	}

	line := strings.Split(strTrimmed, " ")
	line = removeEmpty(line)

	// is_active
	e.IsActive = true
	for ; line[0] == "#" || line[0][0] == '#' && len(line) > 0; {
		e.IsActive = false
		if line[0] == "#" {
			line = line[1:]
		} else {
			line[0] = line[0][1:]
		}
	}

	if len(line) > 0 {
		e.name = strings.TrimSpace(line[0])

		if len(line) > 1 {
			vPtr, otherTokens, err := getValue(line[1:])
			if err != nil {
				if err == CommentOnly {
					return nil
				}

				// TODO return error
				return nil
			}

			if vPtr != nil {
				e.Value = *vPtr
			}

			if len(otherTokens) > 0 && otherTokens[0][0] == '#' {
				otherTokens[0] = otherTokens[0][1:]
				e.Comment = strings.TrimSpace(strings.Join(otherTokens, " "))
			}

			return e
		}
	} else {
		return nil
	}

	return e
}

func (e *ConfigEntry) String() string {
	var res string
	if !e.IsActive {
		res = "# "
	}
	res += e.Name()
	if len(e.Value) > 0 {
		res += " " + e.Value
	}

	if len(e.Comment) > 0 {
		res += " # " + e.Comment
	}

	return res
}

func (e *ConfigEntry) StringLn() string {
	return e.String() + "\n"
}

func (e *ConfigEntry) ChooseActiveOrReduce(e1 *ConfigEntry) ConfigEntryInterface {
	if !e.IsActive {
		return e1
	}
	if e1.IsActive {
		return e1
	}
	return e
}

func (e *ConfigEntry) MakeMultiple() *ConfigMultiEntry {
	res := new(ConfigMultiEntry)
	res.name = e.Name()
	res.Entries = make(map[string]*ConfigEntry)
	res.Entries[e.Value] = e
	return res
}

func (e *ConfigEntry) Copy() ConfigEntryInterface {
	res := new(ConfigEntry)
	*res = *e
	return res
}