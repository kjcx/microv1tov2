/*
 * @Author: kjcx
 * @Date: 2020-08-04 09:46:21
 * @LastEditTime: 2020-08-04 14:37:36
 * @Description: microv1 to microv2 将micro v1 升级到v2 代码批量替换
 * @FilePath:
 */
package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/**
* 读取文件
 */
func Read(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var lines []string
	for {
		line, _ := r.ReadString('\n')
		if line == "" {
			break
		}
		lines = append(lines, strings.Trim(line, "\r\n"))
	}
	return lines, nil
}

/**
* 写入文件
 */
func write(file string, lines []string, filter string, content string) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, line := range lines {

		if strings.Contains(line, filter) {
			fmt.Fprintf(f, "%s\r\n", content)
		} else {
			fmt.Fprintf(f, "%s\r\n", line)
		}
	}
	return nil
}

/**
* 文件file,按filter匹配字符串并替换成content
 */
func replace_file(file, filter string, content string) error {
	lines, err := Read(file)
	if err != nil {
		return err
	}
	return write(file, lines, filter, content)
}

// go.mod 替换下 replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
//项目目录
var Dir = "/ckp"

func main() {
	pwd, _ := os.Getwd()
	//生成新的proto
	GenerateProto(pwd + Dir)
	files := ReadFileSuffix(pwd+Dir, ".go")
	//获取项目下所有文件进行替换
	for _, filename := range files {
		replace_file(filename, "github.com/micro/go-micro/metadata\"", "	\"github.com/micro/go-micro/v2/metadata\"")
		replace_file(filename, "github.com/micro/go-micro/registry\"", "	\"github.com/micro/go-micro/v2/registry\"")
		replace_file(filename, "github.com/micro/go-micro/server\"", "	\"github.com/micro/go-micro/v2/server\"")
		replace_file(filename, "github.com/micro/go-micro\"", "	\"github.com/micro/go-micro/v2\"")
		replace_file(filename, "github.com/micro/go-micro/client\"", "	\"github.com/micro/go-micro/v2/client\"")

	}
}

/**
读取dir目录下所有文件 返回路径
*/
func ReadFileSuffix(dir string, suffix string) []string {
	var files []string
	//获取当前目录下的所有文件或目录信息
	filepath.Walk(dir, func(pathstr string, info os.FileInfo, err error) error {
		fileSuffix := path.Ext(pathstr) //获取文件后缀
		if fileSuffix == suffix {       //按文件后缀读取
			// fmt.Println(pathstr) //打印path信息
			files = append(files, pathstr)
		}
		return nil
	})
	return files
}
