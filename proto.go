/*
 * @Author: kjcx
 * @Date: 2020-08-04 09:46:21
 * @LastEditTime: 2020-08-04 14:38:10
 * @Description: microv1 to microv2 将micro v1 升级到v2 代码批量替换proto生成方法
 * @FilePath:
 */
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Run() {
	var Dir = "/ckp"
	GenerateProto(Dir)
}

/**
*	生成v2 proto
 */
func GenerateProto(dir string) {
	files := ReadFileSuffix(dir, ".proto")
	fmt.Println("files:", files)
	for _, v := range files {

		dir, file := path.Split(v)
		fmt.Println("dddddd", dir, file, v)
		//删除v1生成的proto 文件后缀.pb.go .micro.go
		file_strs := strings.Split(file, ".")
		// fmt.Println("file_names", file_names, v, dir)
		Cmd("rm", []string{"-rf", dir + file_strs[0] + ".pb.go", dir + file_strs[0] + ".micro.go"})
		// 生成新的proto文件
		s, err := CmdAndChangeDir(dir, "protoc", []string{"--proto_path=.", "--micro_out=.", "--go_out=.", file})
		fmt.Println("res:", s, err)
	}
}

/**
* 进入目录执行命令
 */
func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDir", dir, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

/**
* 执行cmd命令
 */
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

/**
* 执行cmd命令字符串
 */
func ExecShell(cmdstring string) string {
	cmd := exec.Command("/bin/bash", "-c", cmdstring)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
		return err.Error()
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return err.Error()
	}
	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("ReadAll Stdout:", err.Error())
		return string(bytes)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("wait:", err.Error(), cmdstring, string(bytes))

		return err.Error()
	}
	return string(bytes)
}
