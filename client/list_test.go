package client

// 必须导入testing模块，并且方法的接受者为(t *testing.T)
import (
	"fmt"
	"testing"
)

// 测试1: 判断kv是否已经存在
func TestList(t *testing.T) {
	fmt.Println("the rectangular name's result is ok")
}
