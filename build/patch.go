package build

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	vultureosModulePath = "github.com/vultureOS/vulture"
	vultureosImportFile = "import_vultureos.go"
	overlayFile         = "overlay.json"
)

var (
	vultureosImportTpl = template.Must(template.New("vultureos").Parse(`
	//+build vultureos

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

func (b *Builder) vultureosImportFile() string {
	return filepath.Join(b.basedir, vultureosImportFile)
}

func (b *Builder) overlayFile() string {
	return filepath.Join(b.basedir, overlayFile)
}

func (b *Builder) buildPrepare() error {
	var err error

	if !b.modHasvultureos() {
		log.Printf("vultureos not found in go.mod")
		err = b.editGoMod()
		if err != nil {
			return err
		}
	}

	err = b.writeImportFile(b.vultureosImportFile())
	if err != nil {
		return err
	}

	err = writeOverlayFile(b.overlayFile(), vultureosImportFile, b.vultureosImportFile())
	if err != nil {
		return err
	}
	return nil
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

func (b *Builder) readGomodule() (*gomodule, error) {
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

func (b *Builder) modHasvultureos() bool {
	if b.currentModulePath() == vultureosModulePath {
		return true
	}

	mods, err := b.readGomodule()
	if err != nil {
		panic(err)
	}
	for _, mod := range mods.Require {
		if mod.Path == vultureosModulePath {
			return true
		}
	}
	return false
}

func (b *Builder) editGoMod() error {
	getPath := vultureosModulePath
	if b.cfg.vultureOSVersion != "" {
		getPath = getPath + "@" + b.cfg.vultureOSVersion
	}
	log.Printf("go get %s", getPath)
	env := []string{
		"GOOS=linux",
		"GOARCH=amd64",
	}
	env = append(env, os.Environ()...)
	cmd := exec.Command(b.gobin(), "get", getPath)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (b *Builder) currentPkgName() string {
	out, err := exec.Command(b.gobin(), "list", "-f", `{{.Name}}`).CombinedOutput()
	if err != nil {
		log.Panicf("get current package name:%s", out)
	}
	return strings.TrimSpace(string(out))
}

func (b *Builder) currentModulePath() string {
	out, err := exec.Command(b.gobin(), "list", "-f", `{{.Module.Path}}`).CombinedOutput()
	if err != nil {
		log.Panicf("get current module path:%s", out)
	}
	return strings.TrimSpace(string(out))

}

func (b *Builder) writeImportFile(fname string) error {
	pkgname := b.currentPkgName()
	var rawFile bytes.Buffer
	err := vultureosImportTpl.Execute(&rawFile, map[string]interface{}{
		"name": pkgname,
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(fname, rawFile.Bytes(), 0644)
	return err
}
