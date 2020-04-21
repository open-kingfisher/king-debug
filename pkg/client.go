package pkg

import (
	"context"
	"os"
)

func DockerClient() (*Docker, error) {
	// 从容器中获取相关环境变量，实例化Docker结构体
	image := os.Getenv("DEBUG_IMAGE")        // 使用的debug镜像，里面可以携带各种工具包
	entryPoint := os.Getenv("ENTRY_POINT")   // debug镜像的entryPoint
	containerId := os.Getenv("CONTAINER_ID") // 要debug的容器
	return NewDocker(image, containerId, entryPoint, context.TODO())
}
