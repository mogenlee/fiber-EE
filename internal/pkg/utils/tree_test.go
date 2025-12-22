package utils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/samber/lo"
)

// 1. 实体实现 TreeNode 接口
type Menu struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name"`
	Type     int64  `json:"type"`
}

func (m Menu) GetID() int64       { return m.ID }
func (m Menu) GetParentID() int64 { return m.ParentID }

func TestTree(t *testing.T) {
	// 2. 构建树
	menus := []Menu{
		{ID: 1, ParentID: 0, Name: "系统管理"},
		{ID: 2, ParentID: 1, Name: "用户管理"},
		{ID: 3, ParentID: 1, Name: "角色管理"},
		{ID: 4, ParentID: 0, Name: "内容管理"},
		{ID: 5, ParentID: 2, Name: "内容管理2"},
		{ID: 6, ParentID: 5, Name: "内容管理5"},
	}

	tree := BuildTree(menus, 0)

	bytes, err := json.Marshal(tree)
	fmt.Println(err)
	fmt.Println(string(bytes))

	flat := FlatTree(tree) // 树转列表
	fmt.Println(flat)
	node := FindTreeNode(tree, 2) // 查找节点
	fmt.Println(node)
	ids := GetTreeIDs(tree) // 获取所有 ID
	fmt.Println(ids)

}
func Test3(t *testing.T) {
	status := lo.Ternary(true, "成年", "未成年") // "成年"
	fmt.Println(status)
}
