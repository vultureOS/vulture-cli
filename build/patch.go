package build

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

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

type buildOverlay struct {
	Replace map[string]string
}

func (b *Builder) pranaosImportFile() string {
	return filepath.Join(b.basedir, vultureosImportFile)
}

func (b *Builder) overlayFile() string {
	return filepath.Join(b.basedir, overlayFile)
}

func writeOverlayFile(overlayFile, dest, source string) error {
	overlay := buildOverlay{
		Replace: map[string]string{
			dest: source,
		},
	}
	buf, _ := json.Marshal(overlay)
	return os.WriteFile(overlayFile, buf, 0644)
}

func (b *Builder) readGoModule() (*gomodule, error) {
	var buf bytes.Buffer
	cmd := exec.Command(b.gobin(), "mod", "edit", "-json")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var mod gomodule
	err = json.Unmarshal(buf.Bytes(), &mod)
	if err != nil {
		return nil, err
	}

	return &mod, nil
}
