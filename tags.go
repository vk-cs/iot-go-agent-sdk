package iot_go_agent_sdk

import "github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/models"

var (
	StatusTagPath          = []string{"$state", "$status"}
	ConfigVersionTagPath   = []string{"$state", "$config", "$version"}
	ConfigUpdatedAtTagPath = []string{"$state", "$config", "$updated_at"}
)

// FindTagByPath returns tag by path from given agent or device tags subtree
func FindTagByPath(tags []*models.TagConfigObject, path []string) (*models.TagConfigObject, bool) {
	for _, tag := range tags {
		if *tag.Name == path[0] {
			if len(path) == 1 {
				return tag, true
			} else {
				return FindTagByPath(tag.Children, path[1:])
			}
		}
	}

	return nil, false
}

type TagNode struct {
	Tag      models.TagConfigObject
	Children map[string]TagNode
}

func (node TagNode) GetPath(path []string) (TagNode, bool) {
	if len(path) == 0 {
		return TagNode{}, false
	}

	child, ok := node.Children[path[0]]
	if !ok {
		return TagNode{}, false
	}

	if len(path) == 1 {
		return child, true
	} else {
		return child.GetPath(path[1:])
	}
}

func NewTagNode(tag models.TagConfigObject) TagNode {
	childrenMap := make(map[string]TagNode, len(tag.Children))
	for _, child := range tag.Children {
		childrenMap[*tag.Name] = NewTagNode(*child)
	}

	return TagNode{
		Tag:      tag,
		Children: childrenMap,
	}
}

type TagTree struct {
	root TagNode
}

func (tree *TagTree) GetRoot() TagNode {
	return tree.root
}

func (tree *TagTree) GetPath(path []string) (TagNode, bool) {
	return tree.root.GetPath(path)
}

func (tree *TagTree) GetStatusTag() (TagNode, bool) {
	return tree.root.GetPath(StatusTagPath)
}

func (tree *TagTree) GetConfigVersionTag() (TagNode, bool) {
	return tree.root.GetPath(ConfigVersionTagPath)
}

func (tree *TagTree) ConfigUpdatedAtTagPath() (TagNode, bool) {
	return tree.root.GetPath(ConfigUpdatedAtTagPath)
}

func NewTagTree(root models.TagConfigObject) *TagTree {
	return &TagTree{
		root: NewTagNode(root),
	}
}
