package build

import "text/template"

const (
	vultureosModulePath = "github.com/vultureOS/vulture"
	vultureosImportFile = "import_vultureos.go"
	overlayFile         = "overlay.json"
)

var (
	vultureosImportTpl = template.Must(template.New("pranaos").Parse(`
	// + build pranaos
	
	package {{.name}}
	import _ "github.com/vultureOS/vulture"
	`))
)

type gomodule struct {
	Module struct {
		Path string `json:"Path"`
	} `json:"Module"`
	Go      string `json:"Go"`
	Require []struct {
		Path    string `json:"Path"`
		Version string `json:"Version"`
	} `json:"Require"`
	Exclude interface{} `json:"Exclude"`
	Replace interface{} `json:"Replace"`
	Retract interface{} `json:"Retract"`
}
