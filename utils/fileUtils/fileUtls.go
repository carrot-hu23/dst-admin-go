package fileUtils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func CreateDir(dirName string) bool {
	if dirName == "" {
		return false
	}
	if Exists(dirName) {
		return false
	}
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CreateFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return err
}

func WriterTXT(filename, content string) error {
	// 写入文件
	// 判断文件是否存在
	var file *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {

		file, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		//O_APPEND
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err2 := w.WriteString(content)
	if err2 != nil {
		return err2
	}
	w.Flush()
	file.Sync()
	return nil

}

func ReadLnFile(filePath string) ([]string, error) {

	//打开文件
	fi, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	buf := bufio.NewScanner(fi)
	// 循环读取
	var lineArr []string
	for {
		if !buf.Scan() {
			break //文件读完了,退出for
		}
		line := buf.Text() //获取每一行
		lineArr = append(lineArr, line)
	}

	return lineArr, nil
}

func ReadFile(filePath string) (string, error) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error: ", err)
		return "", err
	}
	return string(data), err
}

func WriterLnFile(filename string, lines []string) error {

	var file *os.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {

		file, err = os.Create(filename)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		//O_APPEND
		file, err = os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	for _, v := range lines {
		fmt.Fprintln(w, v)
	}
	return w.Flush()
}

func ReverseRead(name string, lineNum uint) ([]string, error) {
	//打开文件
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件大小
	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fs.Size()

	var offset int64 = -1   //偏移量，初始化为-1，若为0则会读到EOF
	char := make([]byte, 1) //用于读取单个字节
	lineStr := ""           //存放一行的数据
	buff := make([]string, 0, 100)
	for (-offset) <= fileSize {
		//通过Seek函数从末尾移动游标然后每次读取一个字节
		file.Seek(offset, io.SeekEnd)
		_, err := file.Read(char)
		if err != nil {
			return buff, err
		}
		if char[0] == '\n' {
			offset--  //windows跳过'\r'
			lineNum-- //到此读取完一行
			buff = append(buff, lineStr)
			lineStr = ""
			if lineNum == 0 {
				return buff, nil
			}
		} else {
			lineStr = string(char) + lineStr
		}
		offset--
	}
	buff = append(buff, lineStr)
	return buff, nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		log.Printf("remove "+path+" : %v\n", err)
		return err
	}
	return nil
}

func DeleteDir(path string) (err error) {
	err = os.RemoveAll("./file1.txt")
	if err != nil {
		log.Printf("removeAll "+path+" : %v\n", err)
	}
	return
}

func Rename(filePath, newName string) (err error) {
	err = os.Rename(filePath, newName)
	return
}

func FindWorldDirs(rootPath string) ([]string, error) {
	var dirs []string

	// 遍历目录并列出满足条件的目录
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果是目录且名称包含 master 或 caves（不区分大小写）
		if info.IsDir() && (strings.Contains(strings.ToLower(info.Name()), "master") || strings.Contains(strings.ToLower(info.Name()), "caves")) {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func ListDirectories(root string) ([]string, error) {
	var dirs []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(dirs, func(i, j int) bool {
		fi, err := os.Stat(dirs[i])
		if err != nil {
			return false
		}
		fj, err := os.Stat(dirs[j])
		if err != nil {
			return false
		}
		return fi.ModTime().Before(fj.ModTime())
	})

	return dirs, nil
}
