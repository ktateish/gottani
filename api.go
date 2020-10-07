package gottani

import (
	"bytes"
	"fmt"

	"github.com/ktateish/gottani/internal/appinfo"
	"github.com/ktateish/gottani/internal/pkginfo"
)

// Combine returns an application source code created by combining all
// functions, vars, consts, types that are reachable form the given entry
// point of the package in the given dir.
func Combine(dir, entryPointName string) ([]byte, error) {
	pi, err := pkginfo.New(dir)
	if err != nil {
		return nil, fmt.Errorf("laoding package information: %w", err)
	}

	ai := appinfo.NewApplicationInfo(pi, entryPointName)

	app, err := ai.Squash()
	if err != nil {
		return nil, fmt.Errorf("creating combined application: %w", err)
	}

	w := new(bytes.Buffer)
	err = app.Fprint(w)
	if err != nil {
		return nil, fmt.Errorf("formatting: %w", err)
	}

	return w.Bytes(), nil
}
