package zip

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func zipDir(dirPath string, zipWriter *zip.Writer, basePath string) error {
	fileInfos, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		path := filepath.Join(dirPath, fileInfo.Name())
		if fileInfo.IsDir() {
			err := zipDir(path, zipWriter, basePath)
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			relPath, err := filepath.Rel(basePath, path)
			if err != nil {
				return err
			}

			zipEntry, err := zipWriter.Create(filepath.ToSlash(relPath))
			if err != nil {
				return err
			}

			_, err = io.Copy(zipEntry, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Zip(sourceDir, targetZip string) error {
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	basePath := filepath.Dir(sourceDir)

	err = zipDir(sourceDir, zipWriter, basePath)
	if err != nil {
		return err
	}

	return nil
}

func Unzip(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			return err
		}

		zipEntry, err := file.Open()
		if err != nil {
			return err
		}
		defer zipEntry.Close()

		targetFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, zipEntry)
		if err != nil {
			return err
		}
	}

	return nil
}

func Unzip3(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	var clusterIniPath string
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "cluster.ini") {
			clusterIniPath = file.Name
			break
		}
	}

	if clusterIniPath == "" {
		return errors.New("压缩包中缺少 cluster.ini 文件")
	}

	clusterPath := filepath.Dir(clusterIniPath)
	clusterPath = strings.TrimSuffix(clusterPath, "/")
	clusterPath = strings.TrimSuffix(clusterPath, "\\")
	log.Println("clusterPath", clusterPath)
	for _, file := range reader.File {
		// 获取相对路径
		relativePath := strings.TrimPrefix(file.Name, clusterPath)
		relativePath = strings.TrimPrefix(relativePath, "/")
		relativePath = strings.TrimPrefix(relativePath, "\\")

		extractedFilePath := filepath.Join(destination, relativePath)
		log.Println("extractedFilePath", extractedFilePath)
		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm); err != nil {
			return err
		}

		writer, err := os.Create(extractedFilePath)
		if err != nil {
			return err
		}
		defer writer.Close()

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		if _, err := io.Copy(writer, reader); err != nil {
			return err
		}
	}

	return nil
}

func Unzip2(zipFile, destDir, newName string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		parts := strings.Split(f.Name, string(filepath.Separator)) // 使用路径分隔符拆分路径
		parts = parts[1:]                                          // 去掉第一个部分，即一级目录
		newPath := filepath.Join(parts...)                         // 组合新的路径

		// 构建解压后的文件路径
		extractedFilePath := filepath.Join(destDir, newName, newPath)
		log.Println(">>> ", destDir, newName, newPath)
		if f.FileInfo().IsDir() {
			// 创建目录
			os.MkdirAll(extractedFilePath, os.ModePerm)
			continue
		}

		// 创建解压后的文件
		if err := os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		// 打开压缩文件中的文件
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 将压缩文件中的内容复制到解压文件中
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
