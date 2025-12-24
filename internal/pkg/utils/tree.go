package utils

import (
	"sort"

	"github.com/bytedance/sonic"
	"github.com/spf13/cast"
)

const (
	treeDefaultId       = "id"
	treeDefaultPid      = "pid"
	treeDefaultChildren = "children"
	treeBeginCut        = "tr_"
	treeEndCut          = " "
)

type GenOption struct {
	IdField       string
	PidField      string
	ChildrenField string
}

func GenTree(obj any) ([]map[string]any, error) {
	menus, err := toMaps(obj)
	if err != nil {
		return nil, err
	}
	return GenTreeWithField(menus, GenOption{
		IdField:       treeDefaultId,
		PidField:      treeDefaultPid,
		ChildrenField: treeDefaultChildren,
	}), nil
}

func GenTreeWithField(obj any, op GenOption) []map[string]any {
	menus, err := toMaps(obj)
	if err != nil || len(menus) == 0 {
		return nil
	}

	// 初始化根节点PID为最小值（通常为0）
	minPid := GetMinPid(menus, op.PidField)

	// 构建ID索引表，便于快速查找节点
	formatMenu := make(map[int]map[string]any)
	var rootMenus []map[string]any // 根节点列表

	// 第一轮遍历：建立索引并收集根节点
	for _, m := range menus {
		id := cast.ToInt(m[op.IdField])
		formatMenu[id] = m // 按ID存储节点

		// 当PID等于最小PID时视为根节点
		if cast.ToInt(m[op.PidField]) == minPid {
			rootMenus = append(rootMenus, m)
		}
	}

	// 第二轮遍历：为每个节点找到父节点并添加到子节点列表
	for _, m := range menus { // 注意：遍历原始切片保证顺序（与输入顺序一致）
		pid := cast.ToInt(m[op.PidField])
		parent, exists := formatMenu[pid]
		if !exists {
			continue // 父节点不存在时跳过（可能为无效数据）
		}

		// 安全获取或初始化子节点切片（保证按遍历顺序追加到尾部）
		children, ok := parent[op.ChildrenField].([]map[string]any)
		if !ok {
			children = make([]map[string]any, 0)
		}
		children = append(children, m) // 按原始顺序追加到尾部
		parent[op.ChildrenField] = children
	}

	return rootMenus
}

// sortChildren 对子节点进行排序（保持原有顺序可注释此函数）
func sortChildren(menus []map[string]any, childrenField, idField string) {
	for _, m := range menus {
		children, ok := m[childrenField].([]map[string]any)
		if !ok {
			continue
		}

		// 按ID升序排序（如需保持原始顺序，删除此排序逻辑）
		sort.Slice(children, func(i, j int) bool {
			return cast.ToInt(children[i][idField]) < cast.ToInt(children[j][idField])
		})

		// 递归排序子节点的子节点
		sortChildren(children, childrenField, idField)
	}
}

func GetMinPid(menus []map[string]any, pidField string) int {
	index := -1
	for _, m := range menus {
		pid := cast.ToInt(m[pidField])
		if index == -1 {
			index = pid
			continue
		}
		if pid < index {
			index = pid
		}
	}
	return max(index, 0) // 确保返回非负数（避免-1作为PID）
}

// 辅助函数：取最大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// toMaps 使用 sonic.Unmarshal 将输入的对象转换为 map 切片
func toMaps(obj any) ([]map[string]any, error) {
	// 使用 sonic.Marshal 将输入对象序列化为字节数组
	data, err := sonic.Marshal(obj) // 这里不需要使用 v2，直接使用 json.Marshal
	if err != nil {
		return nil, err
	}

	// 使用 sonic.Unmarshal 解析字节数组为 map 切片
	var result []map[string]any
	err = sonic.Unmarshal(data, &result) // 同样的，直接使用 json.Unmarshal
	if err != nil {
		return nil, err
	}

	return result, nil
}
