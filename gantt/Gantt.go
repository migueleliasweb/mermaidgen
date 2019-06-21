package gantt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

////////// AxisFormat //////////////////////////////////////////////////////////

type axisFormat string

// Format definitions for x axis scale as described at
// https://mermaidjs.github.io/gantt.html#scale.
// The default is no axisFormat statement which results in FormatDateOnly.
const (
	FormatDateOnly               axisFormat = `%Y-%m-%d`
	FormatTimeOnly24h            axisFormat = `%H:%M`
	FormatTimeOnly24hWithSeconds axisFormat = `%H:%M:%S`
)

////////// Gantt ///////////////////////////////////////////////////////////////

// Gantt objects are the entrypoints to this package, the whole diagram is
// constructed around a Gantt object. Create an instance of Gantt via
// Gantt's constructor NewGantt, do not create instances directly.
type Gantt struct {
	sectionsMap map[string]*Section // lookup table for existing Sections
	sections    []*Section          // Section items for ordered rendering
	tasksMap    map[string]*Task    // lookup table for existing Tasks
	tasks       []*Task             // Section-less Task items
	Title       string              // Title of the Gantt diagram
	AxisFormat  axisFormat          // Optional time format for x axis
}

// NewGantt is the constructor used to create a new Gantt object.
// This object is the entrypoint for any further interactions with your diagram.
// Always use the constructor, don't create Gantt objects directly.
func NewGantt() (newGantt *Gantt) {
	g := &Gantt{}
	g.sectionsMap = make(map[string]*Section)
	g.tasksMap = make(map[string]*Task)
	return g
}

// String recursively renders the whole diagram to mermaid code lines.
func (g *Gantt) String() (renderedElement string) {
	text := "gantt\ndateFormat YYYY-MM-DDTHH:mm:ssZ\n"
	if g.AxisFormat != "" {
		text += fmt.Sprintln("axisFormat", g.AxisFormat)
	}
	if g.Title != "" {
		text += fmt.Sprintln("title", g.Title)
	}
	for _, t := range g.tasks {
		text += t.String()
	}
	for _, s := range g.sections {
		text += s.String()
	}
	return text
}

// Structs for JSON encode
type mermaidJSON struct {
	Theme string `json:"theme"`
}
type dataJSON struct {
	Code    string      `json:"code"`
	Mermaid mermaidJSON `json:"mermaid"`
}

// LiveURL renders the Gantt and generates a view URL for
// https://mermaidjs.github.io/mermaid-live-editor from it.
func (g *Gantt) LiveURL() (url string) {
	liveURL := `https://mermaidjs.github.io/mermaid-live-editor/#/edit/`
	data, _ := json.Marshal(dataJSON{
		Code: g.String(), Mermaid: mermaidJSON{Theme: "default"}})
	return liveURL + base64.URLEncoding.EncodeToString(data)
}

// ViewInBrowser uses the URL generated by Gantt's LiveURL method and opens
// that URL in the OS's default browser. It starts the browser command
// non-blocking and eventually returns any error occured.
func (g *Gantt) ViewInBrowser() (err error) {
	switch runtime.GOOS {
	case "openbsd", "linux":
		return exec.Command("xdg-open", g.LiveURL()).Start()
	case "darwin":
		return exec.Command("open", g.LiveURL()).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler",
			g.LiveURL()).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

////////// add Items ///////////////////////////////////////////////////////////

// AddSection ...
func (g *Gantt) AddSection(id string) (newSection *Section) {
	_, alreadyExists := g.sectionsMap[id]
	if alreadyExists {
		return nil
	}
	s := &Section{id: id, gantt: g}
	g.sectionsMap[id] = s
	g.sections = append(g.sections, s)
	return s
}

// AddTask ...
func (g *Gantt) AddTask(id string, init ...interface{}) (newTask *Task, err error) {
	_, alreadyExists := g.tasksMap[id]
	if alreadyExists {
		return nil, fmt.Errorf("id already exists")
	}
	t, e := taskNew(id, g, nil, init)
	if e != nil {
		return nil, e
	}
	g.tasksMap[id] = t
	g.tasks = append(g.tasks, t)
	return t, nil
}

////////// get Items ///////////////////////////////////////////////////////////

// GetSection ...
func (g *Gantt) GetSection(id string) (existingSection *Section) {
	// if not found -> nil
	return g.sectionsMap[id]
}

// GetTask ...
func (g *Gantt) GetTask(id string) (existingTask *Task) {
	// if not found -> nil
	return g.tasksMap[id]
}

////////// list Items //////////////////////////////////////////////////////////

// ListSections ...
func (g *Gantt) ListSections() (allSections []*Section) {
	values := make([]*Section, 0, len(g.sectionsMap))
	for _, v := range g.sectionsMap {
		values = append(values, v)
	}
	return values
}

// ListTasks ...
func (g *Gantt) ListTasks() (allTasks []*Task) {
	values := make([]*Task, 0, len(g.tasksMap))
	for _, v := range g.tasksMap {
		values = append(values, v)
	}
	return values
}