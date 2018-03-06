package pack

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/Azure/draft/pkg/osutil"
)

const (
	// ChartfileName is the default Chart file name.
	ChartfileName = "Chart.yaml"
	// ValuesfileName is the default values file name.
	ValuesfileName = "values.yaml"
	// IgnorefileName is the name of the Helm ignore file.
	IgnorefileName = ".helmignore"
	// DeploymentName is the name of the deployment file.
	DeploymentName = "deployment.yaml"
	// ServiceName is the name of the service file.
	ServiceName = "service.yaml"
	// IngressName is the name of the ingress file.
	IngressName = "ingress.yaml"
	// NotesName is the name of the NOTES.txt file.
	NotesName = "NOTES.txt"
	// HelpersName is the name of the helpers file.
	HelpersName = "_helpers.tpl"
	// TemplatesDir is the relative directory name for templates.
	TemplatesDir = "templates"
	// ChartsDir is the directory name for the packaged chart.
	// This also doubles as the directory name for chart dependencies.
	ChartsDir = "charts"
	// HelmignoreName is the name of the Helm ignore file
	HelmignoreName = ".helmignore"
)

// Pack defines a Draft Starter Pack.
type Pack struct {
	// Chart is the Helm chart to be installed with the Pack.
	Charts []*chart.Chart
	// Files are the files inside the Pack that will be installed.
	Files map[string]io.ReadCloser
}

// SaveDir saves a pack as files in a directory.
func (p *Pack) SaveDir(dest string) error {
	// Create the chart directory
	chartPath := filepath.Join(dest, ChartsDir)
	if err := os.Mkdir(chartPath, 0755); err != nil {
		return fmt.Errorf("Could not create %s: %s", chartPath, err)
	}
	for _, chart := range p.Charts {
		if err := chartutil.SaveDir(chart, chartPath); err != nil {
			return err
		}
	}

	// save the rest of the files
	for relPath, f := range p.Files {
		path := filepath.Join(dest, relPath)
		exists, err := osutil.Exists(path)
		if err != nil {
			return err
		}
		if !exists {
			newfile, err := os.Create(path)
			if err != nil {
				return err
			}
			defer newfile.Close()
			defer f.Close()
			io.Copy(newfile, f)
		}
	}

	return nil
}
