package service

import (
	"fast-learn/service/dto"
	"fmt"
)

var hostService *HostService

type HostService struct {
	BaseService
}

func NewHostService() *HostService {
	if hostService == nil {
		hostService = &HostService{}
	}

	return hostService
}

func (m *HostService) Shutdown(iShutdownHostDTO dto.ShutdownHostDTO) error {
	var errResult error
	stHostIP := iShutdownHostDTO.HostIP
	fmt.Println("stHostIP", stHostIP)

	// 关机处理相关逻辑代码
	// ...

	return errResult
}
