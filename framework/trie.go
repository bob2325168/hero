package framework

import (
	"errors"
	"strings"
)

// Tree 树结构
type Tree struct {
	root *node //根节点
}

// 节点
type node struct {
	isLast   bool                // 这个节点是否可以成为最终的路由规则，该节点是否成成为一个独立的uri，是否自身就是一个终极节点
	segment  string              // uri中的字符串，代表该节点表示的路由中某个段的字符串
	handlers []ControllerHandler // 代表这个节点中包含的控制器，用于最终加载调用
	childs   []*node             // 代表这个节点下的子节点
	parent   *node               // 父节点，双向指针
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{
		root,
	}
}

// 判断一个segment是否是通用的segment，即以:开头
func isWildSegment(seg string) bool {
	return strings.HasPrefix(seg, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(seg string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	//如果segment是通配符，则所有的下一层子节点都满足需求
	if isWildSegment(seg) {
		return n.childs
	}
	nodes := make([]*node, 0, len(n.childs))
	// 过滤所有的下一层子节点
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			//如果下一层子节点有通配符则满足需求
			nodes = append(nodes, cnode)
		} else if cnode.segment == seg {
			//如果下一层节点没有通配符，但是文本完全匹配则满足需求
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

func (n *node) matchNode(uri string) *node {
	//使用分隔符将uri分成2个部分
	segments := strings.SplitN(uri, "/", 2)
	//第一个部分用于匹配下一层子节点
	seg := segments[0]
	if !isWildSegment(seg) {
		seg = strings.ToUpper(seg)
	}
	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(seg)
	//如果当前节点没有一个符合，说明这个URI之前就不存在，直接返回nil
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}
	// 如果只有一个segment，那标记为最后一个
	if len(segments) == 1 {
		//如果segment已经是最后一个节点，判断这些是否有isLast标识
		for _, cn := range cnodes {
			if cn.isLast == true {
				return cn
			}
		}
		//如果都不是最后一个节点返回nil
		return nil
	}
	//如果有2个以上segment，递归每个节点继续查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// AddRouter 增加路由节点
/*
/book/list
/book/:id (冲突)
/book/:id/name
/book/:student/age
/:user/name
/:user/name/:age (冲突)
*/
func (t *Tree) AddRouter(uri string, handlers []ControllerHandler) error {

	n := t.root
	//确认路由是否冲突
	if n.matchNode(uri) != nil {
		return errors.New("router exist: " + uri)
	}
	segments := strings.Split(uri, "/")
	for index, seg := range segments {
		// 最终进入node segment的字段
		if !isWildSegment(seg) {
			seg = strings.ToUpper(seg)
		}
		isLast := index == len(segments)-1
		var objNode *node //标记是否有合适的节点

		childNodes := n.filterChildNodes(seg)
		//如果有匹配的子节点
		if len(childNodes) > 0 {
			for _, cnode := range childNodes {
				if cnode.segment == seg {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			//创建当前node节点
			cnode := newNode()
			cnode.segment = seg
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			// 父节点指针修改
			cnode.parent = n
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}
	return nil
}

// FindHandler 匹配uri
func (t *Tree) FindHandler(uri string) []ControllerHandler {
	//直接复用matchNode函数，uri是不带通用匹配符的地址
	matchNode := t.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}

// 将uri解析为 params
func (n *node) parseParamsFromEndNode(uri string) map[string]string {

	ret := map[string]string{}
	segments := strings.Split(uri, "/")
	cnt := len(segments)
	cur := n
	for i := cnt - 1; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		// 如果是通配符节点
		if isWildSegment(cur.segment) {
			//设置params
			ret[cur.segment[1:]] = segments[i]
		}
		cur = cur.parent
	}
	return ret
}
