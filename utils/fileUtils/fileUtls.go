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

func ReverseRead(filename string, n uint) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0, n)
	scanner := bufio.NewScanner(file)

	// Read lines in reverse order
	for scanner.Scan() {
		lines = append([]string{scanner.Text()}, lines...)
		if len(lines) > int(n) {
			lines = lines[:n]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
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
	err = os.RemoveAll(path)
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

func CreateFileIfNotExists(path string) error {

	// 检查文件是否存在
	_, err := os.Stat(path)
	if err == nil {
		// 文件已经存在，直接返回
		return nil
	}
	if !os.IsNotExist(err) {
		// 其他错误，返回错误信息
		return err
	}

	// 创建文件所在的目录
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		// 创建目录失败，返回错误信息
		return err
	}

	// 创建文件
	_, err = os.Create(path)
	if err != nil {
		// 创建文件失败，返回错误信息
		return err
	}

	// 创建成功，返回 nil
	return nil
}

func CreateDirIfNotExists(filepath string) {
	if !Exists(filepath) {
		CreateDir(filepath)
	}
}

func Copy(srcPath, outFileDir string) error {

	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	// 如果源文件是目录，则递归复制目录
	if srcInfo.IsDir() {
		// 创建目标目录（如果不存在）
		err = os.MkdirAll(outFileDir, srcInfo.Mode())
		if err != nil {
			return err
		}
		// 遍历源目录中的所有文件和子目录，并递归复制它们
		srcDir, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer srcDir.Close()

		files, err := srcDir.Readdir(-1)
		if err != nil {
			return err
		}

		for _, file := range files {
			srcFilePath := filepath.Join(srcPath, file.Name())
			outFilePath := filepath.Join(outFileDir, filepath.Base(srcPath))
			err = Copy(srcFilePath, outFilePath)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return copyHelper(srcPath, outFileDir)
}

func copyHelper(srcPath, outFileDir string) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标目录（如果不存在）
	err = os.MkdirAll(outFileDir, 0755)
	if err != nil {
		return err
	}

	// 创建目标文件
	outFilePath := filepath.Join(outFileDir, filepath.Base(srcPath))
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 复制数据
	_, err = io.Copy(outFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}
