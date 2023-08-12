package git

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"gopkg.in/yaml.v3"
)

func DecodeYamlFileFromGitTree(tree *object.Tree, name string, out any) error {
	file, err := tree.File(name)
	if err != nil {
		if err == object.ErrFileNotFound {
			return nil
		}
		return err
	}
	return DecodeYamlFromGitFile(file, out)
}

func DecodeYamlFromGitFile(file *object.File, out any) error {
	reader, err := file.Blob.Reader()
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(reader)
	return decoder.Decode(out)
}
