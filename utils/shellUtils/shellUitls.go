package shellUtils

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
	"runtime"
)

// 执行shell命令
func Shell(cmd string) (res string, err error) {
	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("cmd.exe", "/c", cmd)
	} else {
		execCmd = exec.Command("bash", "-c", cmd)
	}
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		log.Println("error: " + err.Error())
	}

	output := ConvertByte2String(stderr.Bytes(), GB18030)
	errput := ConvertByte2String(stdout.Bytes(), GB18030)
	//res = fmt.Sprintf("Output:\n%s\nError:\n%s", stdout.String(), stderr.String())

	// log.Printf("shell exec: %s \nOutput:\n%s\nError:\n%s", cmd, output, errput)
	if errput != "" {
		log.Printf("shell exec: %s Error:\n%s", cmd, errput)
	}
	if output != "" {
		log.Printf("shell exec: %s nOutput:\n%s", cmd, output)
	}

	return stdout.String(), err
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
