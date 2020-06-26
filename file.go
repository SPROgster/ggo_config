package ggo

import (
	"bufio"
	"errors"
	"os"
	"sort"
	"strings"
)

type Config struct {
	Fields       map[string]ConfigEntryInterface
	MultipleList map[string]bool
}

func NewConfig() *Config {
	c := new(Config)
	c.MultipleList = make(map[string]bool)
	c.Fields = make(map[string]ConfigEntryInterface)
	return c
}

func (f *Config) CopyScheme() *Config {
	c := NewConfig()
	c.MultipleList = make(map[string]bool, len(f.MultipleList))
	for k, v := range f.MultipleList {
		c.MultipleList[k] = v
	}
	c.Fields = make(map[string]ConfigEntryInterface)

	return c
}

func (f *Config) SetKeyMultiple(name string, isMultple bool) {
	if isMultple {
		f.MultipleList[name] = isMultple
	} else if _, exists := f.MultipleList[name]; exists {
		delete(f.MultipleList, name)
	}
}

func (f *Config) isMultiple(name string) bool {
	if isMultiple, exists := f.MultipleList[name]; exists {
		return isMultiple
	}
	return false
}

func (f *Config) setWhileParsing(e *ConfigEntry) {
	name := e.Name()

	if another, exists := f.Fields[name]; exists {
		f.Fields[name] = another.ChooseActiveOrReduce(e)
	} else {
		if f.isMultiple(name) {
			f.Fields[name] = e.MakeMultiple()
		} else {
			f.Fields[name] = e
		}
	}
}

func (f *Config) Set(e *ConfigEntry) {
	if e == nil {
		return
	}

	name := e.Name()

	if _, exists := f.Fields[name]; exists {
		f.Fields[name] = e
	} else {
		if f.isMultiple(name) {
			f.Fields[name] = e.MakeMultiple()
		} else {
			f.Fields[name] = e
		}
	}
}

func (f *Config) Get(name string) ConfigEntryInterface {
	return f.Fields[name]
}

func (f *Config) Delete(name string) ConfigEntryInterface {
	r, exists := f.Fields[name]

	if exists {
		delete(f.Fields, name)
		return r
	} else {
		return nil
	}
}

func (f *Config) DeleteValue(name string, value string) *ConfigEntry {
	e, exists := f.Fields[name]
	if !exists {
		return nil
	}

	switch v := e.(type) {
	case *ConfigEntry:
		delete(f.Fields, name)
		return v
	case *ConfigMultiEntry:
		return v.Delete(value)
	}

	return nil
}

func (f *Config) FromFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		e := ParseString(line)
		if e == nil {
			continue
		}
		f.setWhileParsing(e)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (f *Config) FromString(str string) {
	f.Fields = make(map[string]ConfigEntryInterface)
	for _, v := range strings.Split(str,"\n") {
		e := ParseString(v)
		if e == nil {
			continue
		}
		f.setWhileParsing(e)
	}
}

func (f *Config) FromStrings(strs []string) {
	f.Fields = make(map[string]ConfigEntryInterface, len(strs))
	for _, v := range strs {
		e := ParseString(v)
		if e == nil {
			continue
		}
		f.setWhileParsing(e)
	}
}

func (f *Config) ParseConfig(data interface{}) error {
	var err error = nil

	switch v := data.(type) {
	case []byte:
		f.FromString(string(v))
	case string:
		f.FromString(v)
	case []string:
		f.FromStrings(v)
	case *os.File:
		err = f.FromFile(v)
	default:
		err = errors.New("invalid data type")
	}

	return err
}

func (f *Config) mergeSchemes(configs ...*Config) {
	for _, c := range configs {
		if c == nil {
			continue
		}
		for k, m := range c.MultipleList {
			if m {
				f.SetKeyMultiple(k, true)
			}
		}
	}
}

func (f *Config) mergeSingleKeys(configs ...*Config) {
	for _, conf := range configs {
		if conf == nil {
			continue
		}
		for k := range conf.Fields {
			if f.isMultiple(k) {
				continue
			}
			f.Set(conf.Get(k).(*ConfigEntry))
		}
	}
}

func (f *Config) mergeMultiKeys(configs ...*Config) {
	for k := range f.MultipleList {
		v := new(ConfigMultiEntry)
		for _, c := range configs {
			if c == nil {
				continue
			}
			v.Merge(c.Fields[k])
		}
		if v.Entries != nil {
			f.Fields[k] = v
		}
	}
}

// Arguments from less to most specific
// most specific overwrite less specific fields
func Merge(configs ...*Config) *Config {
	if len(configs) == 0 {
		return nil
	}

	haveSome := false
	for _, c := range configs {
		if c != nil {
			haveSome = true
			break
		}
	}

	if !haveSome {
		return nil
	}

	f := NewConfig()
	f.mergeSchemes(configs...)
	f.mergeSingleKeys(configs...)
	f.mergeMultiKeys(configs...)

	return f
}

func MergeSchemes(configs ...*Config) *Config {
	if len(configs) == 0 {
		return nil
	}

	haveSome := false
	for _, c := range configs {
		if c != nil {
			haveSome = true
			break
		}
	}

	if !haveSome {
		return nil
	}

	f := NewConfig()
	f.mergeSchemes(configs...)

	return f
}

func (f *Config) Write(Filename string) error {
	file, err := os.Create(Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	keys := make([]string, 0, len(f.Fields))
	for k, _ := range f.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(file)
	for _, k := range keys {
		v := f.Fields[k]
		w.WriteString(v.StringLn())
	}
	return w.Flush()
}

func (f *Config) String() string {
	var res string
	for _, v := range f.Fields {
		res += v.StringLn()
	}
	if len(res) > 1 {
		res = res[:len(res) - 1]
	}
	return res
}

func (f *Config) Len() int {
	return len(f.Fields)
}
