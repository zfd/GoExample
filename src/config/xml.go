package config

import (
	"encoding/xml"
	"os"
	"fmt"
	"errors"
	"strings"
	"io"
)

type xmlNode struct {
	name  string
	attributes map[string]string
	children []*xmlNode
	value string
}

func newNode() *xmlNode {
	node := new(xmlNode)
	node.children = make([]*xmlNode, 0)
	node.attributes = make(map[string]string)
	return node
}

func (node *xmlNode) String() string {
	str := fmt.Sprintf("<%s", node.name)
	for attrName, attrVal := range node.attributes {
		str += fmt.Sprintf(" %s=\"%s\"", attrName, attrVal)
	}
	str += ">"
	str += node.value
	if len(node.children) != 0 {
		for _, child := range node.children {
			str += fmt.Sprintf("%s", child)
		}
	}
	str += fmt.Sprintf("</%s>", node.name)
	return str
}

func (node *xmlNode) unmarshal(startEl xml.StartElement) error {
	node.name = startEl.Name.Local
	for _, v := range startEl.Attr {
		_, alreadyExists := node.attributes[v.Name.Local]
		if alreadyExists {
			return errors.New("tag '" + node.name + "' has duplicated attribute: '" + v.Name.Local + "'")
		}
		node.attributes[v.Name.Local] = v.Value
	}
	return nil
}

func (node *xmlNode) add(child *xmlNode) {
	if node.children == nil {
		node.children = make([]*xmlNode, 0)
	}
	node.children = append(node.children, child)
}

func (node *xmlNode) hasChildren() bool {
	return node.children != nil && len(node.children) > 0
}

func getNextToken(xmlParser *xml.Decoder) (tok xml.Token, err error) {
	if tok, err = xmlParser.Token(); err != nil {
		if err == io.EOF {
			err = nil
			return
		}
		return
	}
	return
}

func unmarshalNode(xmlParser *xml.Decoder, curToken xml.Token) (node *xmlNode, err error) {
	firstLoop := true
	for {
		var tok xml.Token
		if firstLoop && curToken != nil {
			tok = curToken
			firstLoop = false
		} else {
			tok, err = getNextToken(xmlParser)
			if err != nil || tok == nil {
				return
			}
		}

		switch tt := tok.(type) {
		case xml.SyntaxError:
			err = errors.New(tt.Error())
			return
		case xml.CharData:
			value := strings.TrimSpace(string([]byte(tt)))
			if node != nil {
				node.value += value
			}
		case xml.StartElement:
			if node == nil {
				node = newNode()
				err := node.unmarshal(tt)
				if err != nil {
					return nil, err
				}
			} else {
				childNode, childErr := unmarshalNode(xmlParser, tok)
				if childErr != nil {
					return nil, childErr
				}

				if childNode != nil {
					node.add(childNode)
				} else {
					return
				}
			}
		case xml.EndElement:
			return
		}
	}
}

//分割线

func setXmlMap(mp map[string]string, node *xmlNode, depth int, prefix string) {
	if prefix != "" {
		prefix += ">"
	}
	if depth != 0 {
		prefix += node.name
	}
	if node.hasChildren() {
		depth++
		for _, v := range node.children {
			setXmlMap(mp, v, depth, prefix)
		}
	} else {
		if node.value != "" {
			mp[prefix] = node.value
		}
		prefix += ">"
		for k, v := range node.attributes {
			mp[prefix+k] = v
		}
	}
}

func readXml(mp map[string]string, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic("读取配置文件出错：" + err.Error())
	}
	defer file.Close()
	xmlParser := xml.NewDecoder(file)
	node := newNode()
	node, err = unmarshalNode(xmlParser, nil)
	if err != nil {
		panic("处理配置文件出错：" + err.Error())
	}
	setXmlMap(mp, node, 0, "")
}
