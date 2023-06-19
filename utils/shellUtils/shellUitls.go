package shellUtils

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
)

// ExecuteCommand 执行给定的 Shell 命令，并返回输出和错误（如果有的话）。
func ExecuteCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("执行命令时发生错误: %v, 命令输出: %s", err, stderr.String())
	}

	return out.String(), nil
}

// 执行shell命令
func Shell(cmd string) (res string, err error) {
	//var execCmd *exec.Cmd
	//if runtime.GOOS == "windows" {
	//	execCmd = exec.Command("cmd.exe", "/c", cmd)
	//} else {
	//	execCmd = exec.Command("bash", "-c", cmd)
	//}
	//var (
	//	stdout bytes.Buffer
	//	stderr bytes.Buffer
	//)
	//
	//execCmd.Stdout = &stdout
	//execCmd.Stderr = &stderr
	//err = execCmd.Run()
	//if err != nil {
	//	log.Println("error: " + err.Error())
	//}
	//
	//output := ConvertByte2String(stderr.Bytes(), GB18030)
	//errput := ConvertByte2String(stdout.Bytes(), GB18030)
	////res = fmt.Sprintf("Output:\n%s\nError:\n%s", stdout.String(), stderr.String())
	//
	//// log.Printf("shell exec: %s \nOutput:\n%s\nError:\n%s", cmd, output, errput)
	//if errput != "" {
	//	log.Printf("shell exec: %s Error:\n%s", cmd, errput)
	//}
	//if output != "" {
	//	log.Printf("shell exec: %s nOutput:\n%s", cmd, output)
	//}
	//
	//return stdout.String(), err

	return ExecuteCommand(cmd)
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func Chmod(filePath string) error {
	// 获取文件的当前权限
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	// 添加可执行权限
	newMode := fileInfo.Mode() | 0100

	// 更改文件权限
	err = os.Chmod(filePath, newMode)
	if err != nil {
		return err
	}

	return nil
}
