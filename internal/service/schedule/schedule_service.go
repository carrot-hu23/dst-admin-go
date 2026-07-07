package schedule

import (
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/service/backup"
	"dst-admin-go/internal/service/dstConfig"
	"dst-admin-go/internal/service/game"
	"dst-admin-go/internal/service/update"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ScheduleService struct {
	db          *gorm.DB
	cron        *cron.Cron
	cache       sync.Map
	strategyMap map[string]Strategy
	gameProcess game.Process
}

func NewScheduleService(
	db *gorm.DB,
	gameProcess game.Process,
	backupService *backup.BackupService,
	updateService update.Update,
	dstConfigService dstConfig.Config,
) *ScheduleService {
	c := cron.New()

	// 创建策略上下文
	strategyContext := &StrategyContext{
		GameProcess:   gameProcess,
		BackupService: backupService,
		UpdateService: updateService,
		DstConfig:     dstConfigService,
	}

	service := &ScheduleService{
		db:          db,
		cron:        c,
		gameProcess: gameProcess,
		strategyMap: make(map[string]Strategy),
	}

	// 注册策略
	service.strategyMap["backup"] = &BackupStrategy{context: strategyContext}
	service.strategyMap["update"] = &UpdateStrategy{context: strategyContext}
	service.strategyMap["start"] = &StartStrategy{context: strategyContext}
	service.strategyMap["stop"] = &StopStrategy{context: strategyContext}
	service.strategyMap["restart"] = &RestartStrategy{context: strategyContext}
	service.strategyMap["regenerate"] = &RegenerateStrategy{context: strategyContext}
	service.strategyMap["startGame"] = &StartGameStrategy{context: strategyContext}
	service.strategyMap["stopGame"] = &StopGameStrategy{context: strategyContext}
	service.strategyMap["none"] = &NoneStrategy{context: strategyContext}

	// 初始化数据库中的任务
	service.initDBTasks()

	// 启动调度器
	c.Start()

	return service
}

// Stop 停止调度器
func (s *ScheduleService) Stop() {
	s.cron.Stop()
}

// AddJob 添加定时任务
func (s *ScheduleService) AddJob(task *model.JobTask) error {
	entryID, err := s.cron.AddFunc(task.Cron, func() {
		// 发送公告
		s.sendAnnouncement(task.ClusterName, task.Announcement, task.Sleep, task.Times)

		// 执行策略
		strategy, exists := s.strategyMap[task.Category]
		if exists {
			strategy.Execute(task.ClusterName, task.Uuid)
		} else {
			log.Printf("未找到策略: %s", task.Category)
		}
	})

	if err != nil {
		return err
	}

	s.cache.Store(entryID, task.ID)
	return nil
}

// DeleteJob 删除定时任务
func (s *ScheduleService) DeleteJob(jobID int) error {
	entryID := cron.EntryID(jobID)
	taskID, loaded := s.cache.LoadAndDelete(entryID)

	if !loaded {
		log.Printf("未找到任务: %d", jobID)
		return nil
	}

	s.cron.Remove(entryID)

	// 从数据库删除
	if err := s.db.Delete(&model.JobTask{}, taskID.(uint)).Error; err != nil {
		return err
	}

	return nil
}

// GetJobs 获取所有任务列表
func (s *ScheduleService) GetJobs() []map[string]interface{} {
	var results []map[string]interface{}
	entries := s.cron.Entries()
	log.Println("cron size: ", len(entries))

	for _, entry := range entries {
		taskID, exists := s.cache.Load(entry.ID)
		if !exists {
			continue
		}

		task := s.findDB(taskID.(uint))
		if task == nil {
			continue
		}

		results = append(results, map[string]interface{}{
			"clusterName":  task.ClusterName,
			"levelName":    task.LevelName,
			"uuid":         task.Uuid,
			"jobId":        entry.ID,
			"next":         entry.Next,
			"prev":         entry.Prev,
			"valid":        entry.Valid(),
			"cron":         task.Cron,
			"comment":      task.Comment,
			"category":     task.Category,
			"announcement": task.Announcement,
		})
	}

	return results
}

// GetInstructList 获取指令列表
func (s *ScheduleService) GetInstructList() []map[string]string {
	return []map[string]string{
		{"backup": "备份"},
		{"update": "更新"},
	}
}

// CreateJobTask 创建定时任务
func (s *ScheduleService) CreateJobTask(task *model.JobTask) error {
	// 创建数据库记录
	if err := s.db.Create(task).Error; err != nil {
		return err
	}

	// 添加到调度器
	return s.AddJob(task)
}

// initDBTasks 初始化数据库中的任务
func (s *ScheduleService) initDBTasks() {
	var jobTaskList []model.JobTask
	s.db.Find(&jobTaskList)

	for i := range jobTaskList {
		task := &jobTaskList[i]
		entryID, err := s.cron.AddFunc(task.Cron, func() {
			// 发送公告
			s.sendAnnouncement(task.ClusterName, task.Announcement, task.Sleep, task.Times)

			// 执行策略
			strategy, exists := s.strategyMap[task.Category]
			if exists {
				strategy.Execute(task.ClusterName, task.Uuid)
			} else {
				log.Printf("未找到策略: %s", task.Category)
			}
		})

		if err != nil {
			log.Println("初始化任务失败", err)
			continue
		}

		s.cache.Store(entryID, task.ID)
	}
}

// findDB 从数据库查找任务
func (s *ScheduleService) findDB(taskID uint) *model.JobTask {
	task := &model.JobTask{}
	if err := s.db.Where("ID = ?", taskID).First(task).Error; err != nil {
		return nil
	}
	return task
}

// sendAnnouncement 发送游戏内公告
func (s *ScheduleService) sendAnnouncement(clusterName string, announcement string, sleep int, times int) {
	if announcement == "" {
		return
	}

	for i := 0; i < times; i++ {
		log.Println("开始发送公告")
		lines := strings.Split(announcement, "\n")
		log.Println(lines)

		for j := range lines {
			// 使用 c_announce 命令发送公告
			command := "c_announce(\"" + lines[j] + "\")"
			// 向所有世界发送
			s.gameProcess.Command(clusterName, "Master", command)
			s.gameProcess.Command(clusterName, "Caves", command)
		}

		time.Sleep(time.Duration(sleep) * time.Second)
		log.Println("结束发送公告")
	}
}
