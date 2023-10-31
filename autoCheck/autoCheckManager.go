package autoCheck

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"log"
	"sync"
	"time"
)

type AutoCheckManager struct {
	AutoChecks []model.AutoCheck
	statusMap  map[string]chan int
	mutex      sync.Mutex
}

func (m *AutoCheckManager) Start(autoChecks []model.AutoCheck) {

	m.AutoChecks = autoChecks
	m.statusMap = make(map[string]chan int)

	for i := range autoChecks {
		taskId := autoChecks[i].ClusterName + autoChecks[i].LevelName
		m.statusMap[taskId] = make(chan int)
	}

	for i := range autoChecks {
		go func(index int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(r)
				}
			}()
			taskId := autoChecks[index].ClusterName + autoChecks[index].LevelName
			m.run(autoChecks[index], m.statusMap[taskId])
		}(i)
	}

}

func (m *AutoCheckManager) run(task model.AutoCheck, stop chan int) {
	for {
		select {
		case <-stop:
			return
		default:
			m.check(task)
		}
	}
}
func (m *AutoCheckManager) GetAutoCheck(clusterName, levelName string) *model.AutoCheck {
	db := database.DB
	autoCheck := model.AutoCheck{}
	db.Where("cluster_name = ? and level_name", clusterName, levelName).Find(&autoCheck)
	if autoCheck.Interval == 0 {
		autoCheck.Interval = 10
	}
	return &autoCheck
}

// TODO 这里要修改
func (m *AutoCheckManager) check(task model.AutoCheck) {
	autoCheck := m.GetAutoCheck(task.ClusterName, task.LevelName)
	if autoCheck.EnableModUpdate != 1 || autoCheck.EnableDownCheck != 1 {
		time.Sleep(10 * time.Second)
	} else {

		checkInterval := time.Duration(autoCheck.Interval) * time.Minute
		//if !m.checkFunc(m.clusterName, m.levelName, m.bin, m.beta) {
		//	log.Println(m.clusterName, m.name, " is not running, waiting for ", checkInterval)
		//	time.Sleep(checkInterval)
		//	if !m.checkFunc(m.clusterName, m.levelName, m.bin, m.beta) {
		//		log.Println(m.clusterName, m.name, "has not started, starting it...")
		//		err := m.startFunc(m.clusterName, m.levelName, m.bin, m.beta)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//	}
		//}
		time.Sleep(checkInterval)
	}
}

func (m *AutoCheckManager) AddAutoCheckTasks(task model.AutoCheck) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	taskId := task.ClusterName + task.LevelName
	if oldChan, ok := m.statusMap[taskId]; ok {
		close(oldChan) // 关闭旧的通道
	}
	m.statusMap[taskId] = make(chan int)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		m.run(task, m.statusMap[taskId])
	}()

}

func (m *AutoCheckManager) DeleteAutoCheck(taskId string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if ch, ok := m.statusMap[taskId]; ok {
		close(ch)                   // 关闭通道
		delete(m.statusMap, taskId) // 从 statusMap 中删除键值对
	}

}
