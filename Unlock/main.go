package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var sourceFilePath = flag.String("sourcePath", "", "source file path")
var destFilePath = flag.String("destPath", "", "destination file path")

func main() {
	flag.Parse()

	// 获取命令行参数
	sourceFilePathStr := *sourceFilePath
	destFilePathStr := *destFilePath

	// 使用 WaitGroup 确保所有 Goroutine 完成
	var wg sync.WaitGroup

	// 模拟并发调用，假设需要多次重命名文件
	// 比如多个文件需要被处理
	filePaths := []struct {
		source string
		dest   string
	}{
		{sourceFilePathStr, destFilePathStr},
		// 可以在这里添加更多的源文件和目标文件路径
	}

	// 启动 Goroutines 进行并发处理
	for _, path := range filePaths {
		wg.Add(1)
		go func(source, dest string) {
			defer wg.Done() // 完成后调用 Done

			// 执行文件重命名操作
			err := os.Rename(source, dest)
			if err != nil {
				fmt.Printf("Error: %s ----> %s\n", source, dest)
			} else {
				fmt.Printf("Renamed %s to %s\n", source, dest)
			}
		}(path.source, path.dest) // 将正确的参数传递到 Goroutine
	}

	// 等待所有 Goroutine 完成
	wg.Wait()
}
