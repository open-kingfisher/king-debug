package main

import (
	"kingfisher/kf/common/log"
	"kingfisher/king-debug/pkg"
	"os"
)

func main() {
	docker, err := pkg.DockerClient()
	if err != nil {
		log.Errorf("new docker error: %s", err)
		return
	}
	// 1、拉取镜像
	if err := docker.PullImage(); err != nil {
		log.Errorf("pull image error: %s", err)
		return
	}
	// 2、创建、启动Debug容器并获取Debug容器ID
	id, err := docker.RunContainer()
	if err != nil {
		log.Errorf("run container error: %s", err)
		return
	}
	// 3、生成的容器ID存到/etc/DEBUG_CONTAINER_ID中
	f, err := os.Create("/etc/DEBUG_CONTAINER_ID")
	if err != nil {
		log.Errorf("create /etc/DEBUG_CONTAINER_ID error: %s", err)
		return
	}
	_, err = f.Write([]byte(id))
	if err != nil {
		log.Errorf("write /etc/DEBUG_CONTAINER_ID error: %s", err)
		return
	}
}
