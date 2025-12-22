package utils

import "github.com/samber/lo"

// TreeNode 树节点接口
type TreeNode interface {
	GetID() int64
	GetParentID() int64
}

// TreeResult 带 Children 的树节点
type TreeResult[T any] struct {
	Node     T                `json:"node"`
	Children []*TreeResult[T] `json:"children,omitempty"`
}

// BuildTree 构建树结构
// items: 原始列表
// rootID: 根节点的 ParentID（通常为 0）
func BuildTree[T TreeNode](items []T, rootID int64) []*TreeResult[T] {
	// 按 ParentID 分组
	grouped := lo.GroupBy(items, func(item T) int64 {
		return item.GetParentID()
	})

	var build func(parentID int64) []*TreeResult[T]
	build = func(parentID int64) []*TreeResult[T] {
		children := grouped[parentID]
		if len(children) == 0 {
			return nil
		}

		return lo.Map(children, func(item T, _ int) *TreeResult[T] {
			return &TreeResult[T]{
				Node:     item,
				Children: build(item.GetID()),
			}
		})
	}

	return build(rootID)
}

// FlatTree 扁平化树结构（树转列表）
func FlatTree[T any](trees []*TreeResult[T]) []T {
	var result []T
	var flat func(nodes []*TreeResult[T])
	flat = func(nodes []*TreeResult[T]) {
		for _, node := range nodes {
			result = append(result, node.Node)
			if len(node.Children) > 0 {
				flat(node.Children)
			}
		}
	}
	flat(trees)
	return result
}

// FindTreeNode 在树中查找节点
func FindTreeNode[T TreeNode](trees []*TreeResult[T], id int64) *TreeResult[T] {
	for _, node := range trees {
		if node.Node.GetID() == id {
			return node
		}
		if found := FindTreeNode(node.Children, id); found != nil {
			return found
		}
	}
	return nil
}

// GetTreeIDs 获取树中所有节点 ID
func GetTreeIDs[T TreeNode](trees []*TreeResult[T]) []int64 {
	var ids []int64
	var collect func(nodes []*TreeResult[T])
	collect = func(nodes []*TreeResult[T]) {
		for _, node := range nodes {
			ids = append(ids, node.Node.GetID())
			collect(node.Children)
		}
	}
	collect(trees)
	return ids
}
