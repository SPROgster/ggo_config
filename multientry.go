package ggo

type ConfigMultiEntry struct {
	name    string
	Entries map[string]*ConfigEntry
}

func (e *ConfigMultiEntry) Name() string {
	return e.name
}

func (e *ConfigMultiEntry) String() string {
	res := e.StringLn()
	if len(res) > 0 {
		res = res[:len(res)-1]
	}
	return res
}

func (e *ConfigMultiEntry) StringLn() string {
	res := ""
	for _, e := range e.Entries {
		res += e.StringLn()
	}
	return res
}

func (e *ConfigMultiEntry) Get(value string) *ConfigEntry {
	return e.Entries[value]
}

func (e *ConfigMultiEntry) Delete(value string) *ConfigEntry {
	v, existed := e.Entries[value]
	if existed {
		delete(e.Entries, value)
		return v
	}
	return nil
}

func (e *ConfigMultiEntry) Replace(v *ConfigEntry) bool {
	value := v.Value
	_, replaces := e.Entries[value]
	e.Entries[value] = v

	return replaces
}

func (e *ConfigMultiEntry) ChooseActiveOrReduce(v *ConfigEntry) ConfigEntryInterface {
	value := v.Value
	another, replaces := e.Entries[value]
	if replaces {
		e.Entries[value] = another.ChooseActiveOrReduce(v).(*ConfigEntry)
	} else {
		e.Entries[value] = v
	}

	return e
}

func (e *ConfigMultiEntry) Copy() ConfigEntryInterface {
	res := new(ConfigMultiEntry)
	res.Entries = make(map[string]*ConfigEntry, len(e.Entries))
	for k, v := range e.Entries {
		res.Entries[k] = v
	}

	return res
}

func (e *ConfigMultiEntry) Merge(e1 ConfigEntryInterface) {
	if e == nil {
		return
	}
	if e1 == nil {
		return
	}

	if e.Entries == nil {
		e.Entries = make(map[string]*ConfigEntry)
	}

	switch v := e1.(type) {
	case *ConfigEntry:
		e.Entries[v.Value] = v

	case *ConfigMultiEntry:
		for _, v := range v.Entries {
			e.Entries[v.Value] = v
		}
	}

}