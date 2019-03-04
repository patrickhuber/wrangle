package items

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/patrickhuber/wrangle/store"
)

type yamlRenderer struct {
	redactValues bool
}

func (r *yamlRenderer) RenderItems(itemList []store.Item, writer io.Writer) error {
	return r.renderHierarchy(itemList, writer)
}

func (r *yamlRenderer) renderHierarchy(itemList []store.Item, writer io.Writer) error {
	hierarchy := map[string]interface{}{}
	for _, item := range itemList {
		segments := strings.Split(item.Name(), "/")
		node := hierarchy
		for s, segment := range segments {
			isLeaf := s == len(segments)-1
			if isLeaf {
				break
			}
			segmentNode, ok := node[segment]
			if !ok {
				segmentNode = map[string]interface{}{}
				node[segment] = segmentNode
			}
			node = segmentNode.(map[string]interface{})
		}

		// try to handle this in some way, like encoding the value under something else
		if len(node) > 0 {
			return fmt.Errorf("item at path '%s' is a duplicate or conflicts with an other item in the same path", item.Name())
		}

		// node is set to the last element now
		if r.redactValues {
			node["value"] = "<redacted>"
		} else {
			node["value"] = item.Value()
		}
		node["type"] = item.ItemType()

	}
	data, err := yaml.Marshal(&hierarchy)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (r *yamlRenderer) renderList(itemList []store.Item, writer io.Writer) error {
	structureItems := []interface{}{}
	for _, item := range itemList {
		structureItem := map[string]interface{}{}
		if !r.redactValues {
			structureItem["value"] = item.Value()
		} else {
			structureItem["value"] = "<redacted>"
		}
		structureItems = append(structureItems, structureItem)
	}
	data, err := yaml.Marshal(&structureItems)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (r *yamlRenderer) Name() string {
	return "yaml"
}

func NewYamlRenderer(redactValues bool) Renderer {
	return &yamlRenderer{
		redactValues: redactValues,
	}
}
