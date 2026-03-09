package update

import (
	"dst-admin-go/internal/pkg/utils"
	"dst-admin-go/internal/service/dstConfig"
)

func NewUpdateService(dstConfig dstConfig.Config) Update {
	isWindow := utils.IsWindow()
	if isWindow {
		return NewWindowUpdate(dstConfig)
	}
	return NewLinuxUpdate(dstConfig)
}
