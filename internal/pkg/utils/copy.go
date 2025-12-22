package utils

import (
	"log"

	"github.com/jinzhu/copier"
)

// Copy 拷贝结构体
// 用于 entity -> dto/vo 转换
func Copy(to, from any) {
	if err := copier.Copy(to, from); err != nil {
		log.Printf("[Copy] 拷贝失败: %v", err)
	}
}

// CopyTo 拷贝并返回目标对象（泛型）
func CopyTo[T any](from any) *T {
	to := new(T)
	if err := copier.Copy(to, from); err != nil {
		log.Printf("[CopyTo] 拷贝失败: %v", err)
	}
	return to
}
