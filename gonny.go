package gonnydepp

import (
	"fmt"
	"log"
	"os/exec"
)

const (
	must   int = 1
	should int = 10
	might  int = 100
)

var (
	groups       = make([]*Group, 0)
	execNames    = make(map[string]string)
	execNotFound = make([]string, 0)
)

//Group is a way to categorize processes
type Group struct {
	gid         int           //Internal group id
	level       int           //Internal group level (must,might,should,cant)
	executables []*Executable //internal list of executables
}

//Executable constist of information about an executable
type Executable struct {
	eid   int    //internal executable id
	ename string //short name
	epath string //full path to bin
}

//NewGroup returns a new group with a default group id
func NewGroup() *Group {
	g := Group{len(groups) + 1, 0, nil}
	groups = append(groups, &g)
	return &g //Just a placeholder return
}

//Executable adds a new short name for an executable
func (g *Group) Executable(bin string) {
	e := Executable{}
	e.eid = -1
	e.ename = bin
	g.executables = append(g.executables, &e)
}

//Must indicates that all Executables MUST exist in PATH
func (g *Group) Must() {
	g.level = 1
}

//Should indicates that at least one of the executables MUST be in the PATH
func (g *Group) Should() {
	g.level = 10
}

//Might indicates that any executable in the group can exist (support exists)
func (g *Group) Might() {
	g.level = 100
}

func (g *Group) parse() {
	for _, e := range g.executables {
		if path, err := exec.LookPath(e.ename); err == nil {
			log.Printf("Setting epath to %s on %s", path, e.ename)
			e.epath = path
			execNames[e.ename] = e.epath
			log.Printf("epath=%s", e.epath)
			continue
		}
		execNotFound = append(execNotFound, e.ename)
	}
}

func (g *Group) must() error {
	for _, e := range g.executables {
		log.Printf("musting %s (%s)", e.ename, e.epath)
		if !e.exists() {
			return fmt.Errorf("%s doesn't exist in PATH", e.ename)
		}
	}
	return nil
}

func (g *Group) might() error {
	iFound := 0
	var sarr []string
	for _, e := range g.executables {
		if e.exists() {
			iFound++
		} else {
			sarr = append(sarr, e.ename)
		}
	}
	if iFound > 0 {
		return nil
	}
	return fmt.Errorf("no binary in the collection %s was found in PATH", sarr)
}

//Parse parses all groups and returns values and stuff
func Parse() error {
	for _, g := range groups {
		log.Printf("Parsing group %d", g.gid)
		g.parse()
	}
	for _, g := range groups {
		if g.level < should {
			if err := g.must(); err == nil {
				continue
			} else {
				return err
			}
		}

		if g.level < might {
			if err := g.might(); err == nil {
				continue
			} else {
				return err
			}
		}
	}
	return nil
}

func (e *Executable) exists() bool {
	return (len(e.epath) > 0)
}

//Missing returns an array of missing executables
func Missing() []string {
	return execNotFound
}

//Existing returns an array of the existing executables (short name only)
func Existing() []string {
	var sval []string
	for _, s := range execNames {
		sval = append(sval, s)
	}
	return sval
}

//Path returns the path to a given (existing) binary
func Path(short string) string {
	if v, prs := execNames[short]; prs {
		return v
	}
	return ""
}
