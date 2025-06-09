package test

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContextTimeout(t *testing.T) {
	go func() {
		// 创建一个5秒的超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 使用上下文进行某些操作，例如：
		select {
		case <-ctx.Done():
			fmt.Println("超时了，操作未能完成。")
			return
		case <-time.After(10 * time.Second):
			fmt.Println("操作使用了过多时间。")
			return
		}
	}()
	time.Sleep(10 * time.Second)
}
