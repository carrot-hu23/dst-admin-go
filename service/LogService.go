package service

import (
	"dst-admin-go/constant"
	optype "dst-admin-go/constant/opType"
	"dst-admin-go/utils/fileUtils"
	"log"
)

func ReadDstLog(opType int, lineNum uint) []string {

	if opType == optype.READ_MASTER_LOG {
		return ReadDstMasterLog(lineNum)
	}

	if opType == optype.READ_CAVES_LOG {
		return ReadDstCavesLog(lineNum)
	}

	return []string{}
}

func ReadDstMasterLog(lineNum uint) []string {
	logPath := constant.GET_DST_MASTER_LOG_PATH()
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dst master log error:", err)
	}
	return logs
}

func ReadDstCavesLog(lineNum uint) []string {

	logPath := constant.GET_DST_CAVES_LOG_PATH()
	logs, err := fileUtils.ReverseRead(logPath, lineNum)
	if err != nil {
		log.Panicln("read dst caves log error:", err)
	}
	return logs
}
