package service

import (
	"dst-admin-go/entity"
	"dst-admin-go/vo"
)

func GetModKV() []entity.ModKV {

	db := entity.DB
	modKVList := []entity.ModKV{}
	db.Find(&modKVList)
	return modKVList
}

func SaveModKV(modKVParams []vo.ModKVParam) {

	db := entity.DB

	modKVs := []entity.ModKV{}
	for _, modParam := range modKVParams {
		modKV := entity.ModKV{
			UserId:  modParam.UserId,
			ModId:   modParam.ModId,
			Config:  modParam.ModConfig,
			Version: modParam.Version,
		}
		modKVs = append(modKVs, modKV)
	}
	db.Model(&entity.ModKV{}).CreateInBatches(modKVs, len(modKVs))
}

func UpdateModKV() {

}

func DeleteModKV() {

}
