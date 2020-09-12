package ymltoerl

import (
	"io/ioutil"

	"fmt"

	"bytes"

	"strings"

	"gopkg.in/yaml.v3"
)

type indentWriter struct {
	indent int
	buf    *bytes.Buffer
}

func (iw indentWriter) Write(p []byte) (n int, err error) {
	pp := []byte(strings.ReplaceAll(string(p), "\n", "\n"+strings.Repeat(" ", iw.indent)))

	return iw.buf.Write(pp)
}

func ConvertFile(yamlFile string) ([]byte, error) {
	b, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return nil, err
	}

	yamlNode := &yaml.Node{}
	if err := yaml.Unmarshal(b, yamlNode); err != nil {
		return nil, err
	}

	if len(yamlNode.Content) > 1 {
		return nil, fmt.Errorf("yaml document content greater than 1: %d", len(yamlNode.Content))
	}

	return ConvertDocument(yamlNode.Content[0])
}

func ConvertDocument(yamlNode *yaml.Node) ([]byte, error) {
	iw := indentWriter{buf: &bytes.Buffer{}}
	if err := convertDocument(iw, yamlNode); err != nil {
		return nil, err
	}

	return iw.buf.Bytes(), nil
}

func convertDocument(w indentWriter, yamlNode *yaml.Node) error {
	if err := convertNode(w, yamlNode); err != nil {
		return err
	}

	fmt.Fprint(w, ".")

	return nil
}

func convertNode(w indentWriter, yamlNode *yaml.Node) error {
	switch yamlNode.Kind {
	case yaml.SequenceNode:
		last := len(yamlNode.Content) - 1
		switch yamlNode.Style {
		case 33:
			fallthrough
		case yaml.FlowStyle:
			fmt.Fprintf(w, arrayStart(yamlNode))
			for i, childYamlMode := range yamlNode.Content {
				if err := convertNode(w, childYamlMode); err != nil {
					return err
				}
				if i != last {
					fmt.Fprintf(w, ", ")
				}
			}
			fmt.Fprintf(w, arrayEnd(yamlNode))
		default:
			fmt.Fprintf(w, arrayStart(yamlNode))
			w.indent += 2
			for i, childYamlMode := range yamlNode.Content {
				fmt.Fprintf(w, "\n")
				if err := convertNode(w, childYamlMode); err != nil {
					return err
				}
				if i != last {
					fmt.Fprintf(w, ",")
				}
			}
			w.indent -= 2
			fmt.Fprintf(w, "\n"+arrayEnd(yamlNode))
		}

		//
		// return n, nil
	case yaml.MappingNode:
		last := len(yamlNode.Content) - 2
		switch yamlNode.Style {
		case yaml.FlowStyle:
			fmt.Fprintf(w, "[")
			for i := 0; i < len(yamlNode.Content); i += 2 {
				fmt.Fprintf(w, "{%s, ", yamlNode.Content[i].Value)
				if err := convertNode(w, yamlNode.Content[i+1]); err != nil {
					return err
				}
				fmt.Fprintf(w, "}")

				if i != last {
					fmt.Fprintf(w, ", ")
				}
			}
			fmt.Fprintf(w, "]")
		default:
			fmt.Fprintf(w, "[")
			w.indent += 2
			for i := 0; i < len(yamlNode.Content); i += 2 {
				fmt.Fprintf(w, "\n{%s, ", yamlNode.Content[i].Value)
				if err := convertNode(w, yamlNode.Content[i+1]); err != nil {
					return err
				}
				fmt.Fprintf(w, "}")

				if i != last {
					fmt.Fprintf(w, ",")
				}
			}
			w.indent -= 2
			fmt.Fprintf(w, "\n]")
		}
	case yaml.ScalarNode:
		switch yamlNode.Tag {
		case "!!null":
			fmt.Fprintf(w, `nil`)
		case "!!bool":
			fmt.Fprintf(w, `%s`, yamlNode.Value)
		case "!!str":
			if strings.Contains(yamlNode.LineComment, "erlang=bin") {
				fmt.Fprintf(w, `<<"%s">>`, yamlNode.Value)
			} else if strings.Contains(yamlNode.LineComment, "erlang=atom") {
				fmt.Fprintf(w, `%s`, yamlNode.Value)
			} else {
				fmt.Fprintf(w, `"%s"`, yamlNode.Value)
			}

		case "!!int":
			fmt.Fprintf(w, yamlNode.Value)
		case "!!float":
			fmt.Fprintf(w, yamlNode.Value)
		case "!!binary":
			fmt.Fprintf(w, `<<"%s">>`, yamlNode.Value)
		case "!atom":
			fmt.Fprintf(w, `%s`, yamlNode.Value)
		}
	case yaml.AliasNode:
		return fmt.Errorf("alias node not supported")
	case yaml.DocumentNode:
		return fmt.Errorf("document node is not expected")
	default:
		return fmt.Errorf(fmt.Sprintf("document unexpected node: %v", yamlNode.Kind))
	}

	// todo remove
	return nil
}

func arrayStart(yamlNode *yaml.Node) string {
	if yamlNode.Tag == "!tuple" ||
		strings.Contains(yamlNode.LineComment, "erlang=tuple") ||
		strings.Contains(yamlNode.HeadComment, "erlang=tuple") ||
		strings.Contains(yamlNode.FootComment, "erlang=tuple") {
		return "{"
	}

	return "["
}

func arrayEnd(yamlNode *yaml.Node) string {
	if yamlNode.Tag == "!tuple" ||
		strings.Contains(yamlNode.LineComment, "erlang=tuple") ||
		strings.Contains(yamlNode.HeadComment, "erlang=tuple") ||
		strings.Contains(yamlNode.FootComment, "erlang=tuple") {
		return "}"
	}

	return "]"
}
