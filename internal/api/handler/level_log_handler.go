package handler

import (
	"bufio"
	"bytes"
	"context"
	clusterContext "dst-admin-go/internal/pkg/context"
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LevelLogHandler struct {
	archive *archive.PathResolver
}

func NewLevelLogHandler(archive *archive.PathResolver) *LevelLogHandler {
	return &LevelLogHandler{
		archive: archive,
	}
}
func (h *LevelLogHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/game/log/stream", h.Stream)
	router.GET("/api/game/level/server/log", h.GetServerLog)
	router.GET("/api/game/level/server/download", h.DownloadServerLog)
}

// Stream 服务器日志流
// @Summary 服务器日志流
// @Description 获取指定世界的实时日志流 (SSE)
// @Tags log
// @Accept text/event-stream
// @Produce text/event-stream
// @Param clusterName query string false "集群名称"
// @Param levelName query string true "世界名称"
// @Success 200 {string} string "SSE 格式的日志流"
// @Router /api/game/log/stream [get]
func (h *LevelLogHandler) Stream(c *gin.Context) {
	clusterName := clusterContext.GetClusterName(c)
	levelName := c.Query("levelName")
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
	serverLogPath := h.archive.ServerLogPath(clusterName, levelName)
	lines, err := reader.Snapshot(serverLogPath, 100)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	base, baseOK := inferDSTLogBaseTime(lines, time.Now())
	for _, line := range lines {
		writeSSE(w, "log", formatDSTLogLine(line, base, baseOK))
	}
	flusher.Flush()

	// 2️⃣ follow
	ch, cancel, err := reader.Follow(serverLogPath)
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
			if m := hermesLogMarkerRE.FindStringSubmatch(line); len(m) == 3 {
				if wall, err := time.ParseInLocation("2006-01-02 15:04:05", m[1], time.Local); err == nil {
					base = wall.Add(-time.Duration(parseElapsedStringSeconds(m[2])) * time.Second)
					baseOK = true
				}
			}
			writeSSE(w, "log", formatDSTLogLine(line, base, baseOK))
			flusher.Flush()

		case <-heartbeat.C:
			writeSSE(w, "ping", "")
			flusher.Flush()
		}
	}
}

// GetServerLog 获取服务器日志 swagger 注释
// @Summary 获取服务器日志
// @Description 获取指定世界的服务器日志（默认最近100行）
// @Tags log
// @Produce application/json
// @Param clusterName query string false "集群名称"
// @Param levelName query string true "世界名称"
// @Param lines query string false "返回日志行数，默认为100"
// @Success 200 {object} response.Response{data=[]string} "服务器日志列表"
// @Router /api/game/level/server/log [get]
func (h *LevelLogHandler) GetServerLog(ctx *gin.Context) {
	clusterName := clusterContext.GetClusterName(ctx)
	levelName := ctx.Query("levelName")
	lines := ctx.DefaultQuery("lines", "100")
	if clusterName == "" || levelName == "" {
		ctx.JSON(400, gin.H{"error": "cluster and level required"})
		return
	}
	serverLogPath := h.archive.ServerLogPath(clusterName, levelName)
	linesInt, err := strconv.Atoi(lines)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "lines must be a number"})
		return
	}
	read, err := fileUtils.ReverseRead(serverLogPath, uint(linesInt))
	if err != nil {
		ctx.JSON(200, response.Response{
			Code: 500,
			Msg:  "failed to read server log: " + err.Error(),
			Data: nil,
		})
		return
	}
	formatted := formatDSTLogLines(read, time.Now())
	ctx.JSON(200, response.Response{
		Code: 200,
		Data: formatted,
		Msg:  "success",
	})
}

// DownloadServerLog 下载服务器日志 swagger 注释
// @Summary 下载服务器日志
// @Description 下载指定世界的完整服务器日志文件
// @Tags log
// @Produce application/octet-stream
// @Param clusterName query string false "集群名称"
// @Param levelName query string true "世界名称"
// @Success 200 {file} file "服务器日志文件"
// @Router /api/game/level/server/download [get]
func (h *LevelLogHandler) DownloadServerLog(ctx *gin.Context) {
	clusterName := clusterContext.GetClusterName(ctx)
	levelName := ctx.Query("levelName")
	if clusterName == "" || levelName == "" {
		ctx.JSON(400, gin.H{"error": "cluster and level required"})
		return
	}
	serverLogPath := h.archive.ServerLogPath(clusterName, levelName)
	content, err := os.ReadFile(serverLogPath)
	if err != nil {
		ctx.JSON(200, response.Response{Code: 500, Msg: "failed to read server log: " + err.Error(), Data: nil})
		return
	}
	formatted := formatDSTLogWithWallClock(strings.Split(string(content), "\n"), time.Now())
	filename := fmt.Sprintf("%s_%s_server_log_%s.txt", sanitizeDownloadName(clusterName), sanitizeDownloadName(levelName), time.Now().Format("20060102_150405"))
	ctx.Header("Content-Type", "text/plain; charset=utf-8")
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.String(http.StatusOK, formatted)
}

var dstLogPrefixRE = regexp.MustCompile(`^\[(\d+):(\d+):(\d+)\]:?\s*`)
var dstCurrentTimeRE = regexp.MustCompile(`^\[(\d+):(\d+):(\d+)\].*?Current time:\s*(.+)$`)
var hermesLogMarkerRE = regexp.MustCompile(`HERMES_LOG_MARKER\s+.*wall=(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+elapsed=([0-9:]+)`)

func sanitizeDownloadName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "unknown"
	}
	return regexp.MustCompile(`[^A-Za-z0-9._-]+`).ReplaceAllString(name, "_")
}

func formatDSTLogWithWallClock(lines []string, now time.Time) string {
	return strings.Join(formatDSTLogLines(lines, now), "\n")
}

func formatDSTLogLines(lines []string, now time.Time) []string {
	base, ok := inferDSTLogBaseTime(lines, now)
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, formatDSTLogLine(line, base, ok))
	}
	return out
}

func inferDSTLogBaseTime(lines []string, now time.Time) (time.Time, bool) {
	for _, line := range lines {
		m := hermesLogMarkerRE.FindStringSubmatch(line)
		if len(m) == 3 {
			wall, err := time.ParseInLocation("2006-01-02 15:04:05", m[1], time.Local)
			if err == nil {
				return wall.Add(-time.Duration(parseElapsedStringSeconds(m[2])) * time.Second), true
			}
		}
	}
	for _, line := range lines {
		m := dstCurrentTimeRE.FindStringSubmatch(line)
		if len(m) != 5 {
			continue
		}
		current, err := time.Parse(time.ANSIC, strings.TrimSpace(m[4]))
		if err != nil {
			continue
		}
		elapsed := dstLogElapsedSeconds(m[1], m[2], m[3])
		return current.Add(-time.Duration(elapsed) * time.Second), true
	}
	maxElapsed := -1
	for _, line := range lines {
		m := dstLogPrefixRE.FindStringSubmatch(line)
		if len(m) == 4 {
			elapsed := dstLogElapsedSeconds(m[1], m[2], m[3])
			if elapsed > maxElapsed {
				maxElapsed = elapsed
			}
		}
	}
	if maxElapsed >= 0 {
		return now.Add(-time.Duration(maxElapsed) * time.Second), true
	}
	return time.Time{}, false
}

func dstLogElapsedSeconds(hh, mm, ss string) int {
	h, _ := strconv.Atoi(hh)
	m, _ := strconv.Atoi(mm)
	s, _ := strconv.Atoi(ss)
	return h*3600 + m*60 + s
}

func parseElapsedStringSeconds(elapsed string) int {
	parts := strings.Split(strings.TrimSpace(elapsed), ":")
	if len(parts) == 3 {
		return dstLogElapsedSeconds(parts[0], parts[1], parts[2])
	}
	if len(parts) == 2 {
		m, _ := strconv.Atoi(parts[0])
		s, _ := strconv.Atoi(parts[1])
		return m*60 + s
	}
	if len(parts) == 1 {
		s, _ := strconv.Atoi(parts[0])
		return s
	}
	return 0
}

func formatDSTLogLine(line string, base time.Time, ok bool) string {
	m := dstLogPrefixRE.FindStringSubmatch(line)
	if len(m) != 4 {
		return line
	}
	elapsed := dstLogElapsedSeconds(m[1], m[2], m[3])
	body := line[len(m[0]):]
	if ok {
		return fmt.Sprintf("[%s] %s", base.Add(time.Duration(elapsed)*time.Second).Format("2006-01-02 15:04:05"), body)
	}
	return line
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
	serverLogPath string,
	n int,
) ([]string, error) {

	path := serverLogPath

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
	serverLogPath string,
) (<-chan string, func(), error) {

	path := serverLogPath

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
