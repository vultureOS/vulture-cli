package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	workDir        string
	goRoot         string
	baseDir        string
	buildTest      bool
	pranaOSVersion string
	goArgs         []string
}

type Builder struct {
	cfg     Config
	basedir string
}

func newBuilder(cfg Config) *Builder {
	return &Builder{
		cfg: cfg,
	}
}

func (b *Builder) Build() error {
	if b.cfg.baseDir == "" {
		basedir, err := ioutil.TempDir("", "pranaos-build")
		if err != nil {
			return err
		}
		b.basedir = basedir
		defer os.RemoveAll(basedir)
	} else {
		b.basedir = b.cfg.baseDir
	}

	return b.Build()
}

func (b *Builder) gobin() string {
	if b.cfg.goRoot == "" {
		return "go"
	}

	return filepath.Join(b.cfg.goRoot, "bin", "go")
}
