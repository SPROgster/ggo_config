package ggo

import "strings"

type ConfigEntry struct {
	IsActive bool
	name     string
	Value    string
	Comment  string
}

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

func ParseString(str string) *ConfigEntry {
	e := new(ConfigEntry)

	strTrimmed := strings.TrimSpace(strings.ReplaceAll(str, "\t", " "))

	if strTrimmed == "#" || strTrimmed == "" {
		return nil
	}

	line := strings.Split(strTrimmed, " ")
	line = removeEmpty(line)
	for i, s := range line {
		line[i] = strings.TrimSpace(s)
	}

	// is_active
	e.IsActive = true
	for ; line[0] == "#" || line[0][0] == '#' && len(line) > 0; {
		e.IsActive = false
		if (line[0] == "#") {
			line = line[1:]
		} else {
			line[0] = line[0][1:]
		}
	}

	if len(line) > 0 {
		// name
		e.name = line[0]
		line = line[1:]

		if len(line) > 0 {
			v := line[0]
			line = line[1:]

			if v[0] == '#' {
				e.Comment = strings.Join(line, " ")
				return e
			}

			if v[0] == '#' {
				v = v[1:]
				line = append([]string{v}, line...)
				e.Comment = strings.Join(line, " ")
				return e
			}

			e.Value = v

			// Following comment
			if len(line) > 0 {
				v := line[0]
				line = line[1:]

				if v == "#" {
					e.Comment = strings.Join(line, " ")
				} else if v[0] == '#' {
					v = v[1:]
					line = append([]string{v}, line...)
					e.Comment = strings.Join(line, " ")
				} else {
					// Comment string or just string
					return nil
				}
			}
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