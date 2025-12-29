package api

import (
	"bufio"
	"bytes"
	"context"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstUtils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogApi struct {
}

func (h *LogApi) Stream(c *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(c)
	levelName := c.Query("levelName")
	clusterName := cluster.ClusterName
	if clusterName == "" || levelName == "" {
		c.JSON(400, gin.H{"error": "cluster and level required"})
		return
	}

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "streaming unsupported"})
		return
	}

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // nginx

	ctx := c.Request.Context()

	// 1️⃣ snapshot
	lines, err := reader.Snapshot(clusterName, levelName, 100)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, line := range lines {
		writeSSE(w, "log", line)
	}
	flusher.Flush()

	// 2️⃣ follow
	ch, cancel, err := reader.Follow(clusterName, levelName)
	if err != nil {
		writeSSE(w, "error", err.Error())
		flusher.Flush()
		return
	}
	defer cancel()

	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case line, ok := <-ch:
			if !ok {
				return
			}
			writeSSE(w, "log", line)
			flusher.Flush()

		case <-heartbeat.C:
			writeSSE(w, "ping", "")
			flusher.Flush()
		}
	}
}

func writeSSE(w io.Writer, event, data string) {
	if event != "" {
		fmt.Fprintf(w, "event: %s\n", event)
	}

	// data 可能包含换行，必须逐行写
	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		fmt.Fprintf(w, "data: %s\n", scanner.Text())
	}

	fmt.Fprint(w, "\n")
}

var reader = NewFileLogReader()

type FileLogReader struct {
	interval time.Duration
}

func NewFileLogReader() *FileLogReader {
	return &FileLogReader{
		interval: time.Second,
	}
}

func (r *FileLogReader) Snapshot(
	cluster, level string,
	n int,
) ([]string, error) {

	path := dstUtils.GetLevelServerLogPath(cluster, level)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	var (
		size   = stat.Size()
		offset = size
		lines  []string
		buf    []byte
	)

	for offset > 0 && len(lines) < n {
		readSize := int64(4096)
		if offset < readSize {
			readSize = offset
		}

		offset -= readSize
		chunk := make([]byte, readSize)

		_, err := f.ReadAt(chunk, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}

		buf = append(chunk, buf...)

		for {
			idx := bytes.LastIndexByte(buf, '\n')
			if idx < 0 {
				break
			}

			line := strings.TrimRight(string(buf[idx+1:]), "\r")
			lines = append(lines, line)
			buf = buf[:idx]

			if len(lines) >= n {
				break
			}
		}
	}

	// 反转
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines, nil
}

func (r *FileLogReader) Follow(
	cluster, level string,
) (<-chan string, func(), error) {

	path := dstUtils.GetLevelServerLogPath(cluster, level)

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	out := make(chan string, 100)
	ctx, cancel := context.WithCancel(context.Background())

	offset := stat.Size()

	go func() {
		defer close(out)
		defer f.Close()

		reader := bufio.NewReader(f)

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(r.interval):
				stat, err := f.Stat()
				if err != nil {
					continue
				}

				// 文件被 truncate
				if stat.Size() < offset {
					offset = 0
					f.Seek(0, io.SeekStart)
					reader.Reset(f)
				}

				if stat.Size() == offset {
					continue
				}

				f.Seek(offset, io.SeekStart)
				reader.Reset(f)

				for {
					line, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					offset += int64(len(line))
					out <- strings.TrimRight(line, "\r\n")
				}
			}
		}
	}()

	return out, cancel, nil
}
