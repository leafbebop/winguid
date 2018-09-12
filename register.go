package guid

var (
	GUIDMap = make(map[string]GUID)
	NameMap = make(map[GUID]string)
)

// Register lets caller register a GUID with a name, and checks if any duplicates
// appear.
func Register(name string,g GUID) error {
	if v:=GUIDMap[name]; v!=g  {
		return MultipleGUIDError{name,g,v}
	}
	if v:=NameMap[g]; v != name {
		return DuplicateGUIDError{g,name,v}
	}

	GUIDMap[name] = g
	NameMap[g] = name

	return nil
}

type MultipleGUIDError struct {
	Name string
	New,Orig GUID
}

type DuplicateGUIDError struct {
	GUID GUID
	New,Orig string
}

func (m MultipleGUIDError) Error() string {
	return m.Name + " has multiple GUID: \n" + m.New.String() + "\n" + m.Orig.String()
}

func (d DuplicateGUIDError) Error() string {
	return d.GUID.String() + " is assigned to " + d.New + " and " + d.Orig
}

//GUIDOf returns the registered GUID of given name. If there is no name registered, the bool is set to be false
func GUIDOf(name string) (GUID,bool) {
	g,ok := GUIDMap[name]
	return g,ok
}

//NameOf returns the registered name of given GUID. If there is no GUID registered, the bool is set to be false
func NameOf(g GUID) (string,bool) {
	name,ok := NameMap[g]
	return name,ok
}
