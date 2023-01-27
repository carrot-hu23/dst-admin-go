package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
)

var adminlist_txt_path = constant.HOME_PATH + constant.SINGLE_SLASH + constant.DST_ADMIN_LIST_PATH
var blocklist_txt_path = constant.HOME_PATH + constant.SINGLE_SLASH + constant.DST_PLAYER_BLOCK_LIST_PATH

// var adminlist_txt_path = "C:/Users/xm/Desktop/dst-admin-go/dst/adminlist.txt"
// var blocklist_txt_path = "C:/Users/xm/Desktop/dst-admin-go/dst/blocklist.txt"

func GetDstAdminList() (str []string) {
	if !fileUtils.Exists(adminlist_txt_path) {
		return
	}
	str, err := fileUtils.ReadLnFile(adminlist_txt_path)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func GetDstBlcaklistPlayerList() (str []string) {
	if !fileUtils.Exists(blocklist_txt_path) {
		return
	}
	str, err := fileUtils.ReadLnFile(blocklist_txt_path)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func SaveDstAdminList(adminlist []string) {

	err := fileUtils.WriterLnFile(adminlist_txt_path, adminlist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}

func SaveDstBlacklistPlayerList(blacklist []string) {
	err := fileUtils.WriterLnFile(blocklist_txt_path, blacklist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}
