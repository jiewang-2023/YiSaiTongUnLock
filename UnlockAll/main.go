package main

import (
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
)

func main() {
	// 获取当前工作目录和可执行文件路径
	path, _ := os.Executable()
	_, selfName := filepath.Split(path)
	str, _ := os.Getwd()

	// 使用 WaitGroup 等待所有任务完成
	var wg sync.WaitGroup

	// 并发获取所有文件
	allFileChan := make(chan string)
	go func() {
		allFile, _ := getAllFileIncludeSubFolder(str)
		for _, path := range allFile {
			allFileChan <- path
		}
		close(allFileChan)
	}()

	// 创建多个 Goroutine 处理文件操作
	for path := range allFileChan {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()

			// 排除不匹配的文件
			if strutil.AfterLast(path, "Unlock.exe") == "" || strutil.AfterLast(path, selfName) == "" {
				return
			}

			dstFilePath := path + ".temp"

			// 复制文件
			if err := copyFile(path, dstFilePath); err != nil {
				log.Printf("文件 %v 复制失败: %v", path, err)
				return
			}

			// 删除原始文件
			if err := os.Remove(path); err != nil {
				log.Printf("文件 %v 删除失败: %v", path, err)
				return
			}

			// 重命名文件
			renameFile(dstFilePath, path)
		}(path)
	}

	// 等待所有 Goroutines 完成
	wg.Wait()

	log.Println("解密完成，按回车键退出")
	fmt.Scanln()
}

func renameFile(sourcePath, dstFilePath string) {
	str, _ := os.Getwd()
	unlockPath := filepath.Join(str, "Unlock.exe")
	arg := fmt.Sprintf(` -sourcePath="%v" -destPath="%v"`, sourcePath, dstFilePath)
	cmd := exec.Command(unlockPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "/c" + arg}
	output, err := cmd.Output()
	if err != nil {
		log.Println("Failed to run command:", err)
	} else {
		info := string(output)
		if info != "" {
			log.Println(info)
		}
	}
}

func copyFile(sourcePath, dstFilePath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, 1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

// getAllFileIncludeSubFolder 获取目录下所有文件（包含子目录）
func getAllFileIncludeSubFolder(folder string) ([]string, error) {
	var result []string
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return result, err
}
