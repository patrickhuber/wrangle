package feed

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/patrickhuber/wrangle/pkg/packages"
)

type GenerateRequest struct {
	Items []*GenerateItem
}

type GenerateItem struct {
	Package   *GeneratePackage
	Platforms []*GeneratePlatform
	Template  string
}

type GeneratePackage struct {
	Name     string
	Versions []string
}

type GeneratePlatform struct {
	Name          string
	Architectures []string
}

type GenerateResponse struct {
	Packages []*packages.Package
}

type GenerateVersion struct {
	Version   string
	Platforms []*GeneratePlatform
}

func Generate(request *GenerateRequest) (*GenerateResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if len(request.Items) == 0 {
		return &GenerateResponse{}, nil
	}
	packageList := []*packages.Package{}
	for _, i := range request.Items {
		templateString := i.Template
		tmpl := template.New("template")
		tmpl, err := tmpl.Parse(templateString)
		if err != nil {
			return nil, err
		}

		versionList := []*packages.Version{}
		for _, v := range i.Package.Versions {

			writer := &bytes.Buffer{}
			data := map[string]any{
				"version":   v,
				"platforms": i.Platforms,
			}
			err := tmpl.Execute(writer, data)
			if err != nil {
				return nil, err
			}
			version := &packages.Version{
				Version: v,
			}
			versionList = append(versionList, version)
		}

		pkg := &packages.Package{
			Name:     i.Package.Name,
			Versions: versionList,
		}
		packageList = append(packageList, pkg)
	}
	return &GenerateResponse{
		Packages: packageList,
	}, nil
}
