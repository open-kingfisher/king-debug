package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	dockerClient "github.com/docker/docker/client"
	"io"
	"kingfisher/kf/common/log"
	"os"
	"time"
)

var (
	DockerEndpoint      = "unix:///var/run/docker.sock"
	DockerStartTimeout  = 30 * time.Second
	DockerDeleteTimeout = 5 * time.Second
)

type Docker struct {
	Client      *dockerClient.Client
	image       string
	entryPoint  string
	containerId string
	Context     context.Context
}

func NewDocker(image, containerId, entryPoint string, context context.Context) (*Docker, error) {
	client, err := dockerClient.NewClient(DockerEndpoint, "", nil, nil)
	if err != nil {
		return nil, err
	}
	return &Docker{
		Client:      client,
		image:       image,
		containerId: containerId,
		entryPoint:  entryPoint,
		Context:     context,
	}, nil
}

// 拉取镜像
func (d *Docker) PullImage() error {
	reader, err := d.Client.ImagePull(context.Background(), d.image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	// 把输出复制到标准输出
	if _, err := io.Copy(os.Stdout, reader); err != nil {
		return err
	}
	return nil
}

func (d *Docker) RunContainer() (string, error) {
	if id, err := d.CreateContainer(); err != nil {
		return "", err
	} else {
		if err := d.StartContainer(id); err != nil {
			return "", err
		}
		return id, nil
	}
}

func (d *Docker) StartContainer(id string) error {
	ctx, cancel := context.WithTimeout(d.Context, DockerStartTimeout)
	defer cancel()
	if err := d.Client.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *Docker) CreateContainer() (string, error) {
	config := &container.Config{
		Image:     d.image, // 通过环境变量获取此次debug使用的image
		Tty:       true,
		OpenStdin: true,
		StdinOnce: true,
	}
	// 自定义Debug容器的启动命令
	if d.entryPoint != "" {
		config.Entrypoint = strslice.StrSlice([]string{d.entryPoint})
	}
	containerID := fmt.Sprintf("container:%s", d.containerId)
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(containerID),                     // 网络模式与要debug的容器相同
		UsernsMode:  container.UsernsMode(containerID),                      // 用户命名空间与要debug的容器相同
		IpcMode:     container.IpcMode(containerID),                         // IPC(提供命名的共享内存段)命名空间与要debug的容器相同
		PidMode:     container.PidMode(containerID),                         // PID与要debug的容器相同
		CapAdd:      strslice.StrSlice([]string{"SYS_PTRACE", "SYS_ADMIN"}), // container拥有一系列的内核修改权限 docker run --cap-add=SYS_PTRACE http://man7.org/linux/man-pages/man7/capabilities.7.html
	}
	body, err := d.Client.ContainerCreate(context.TODO(), config, hostConfig, nil, "")
	if err != nil {
		return "", err
	}
	return body.ID, err
}

func (d *Docker) CleanContainer(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), DockerDeleteTimeout)
	defer cancel()
	// 等待容器优雅退出
	status, err := d.Client.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	var removeErr error
	select {
	case err := <-err:
		if err != nil {
			log.Errorf("error waiting container exit, kill with --force")
			// 强制删除
			removeErr = d.RmContainer(id, true)
		}
	case <-status:
		removeErr = d.RmContainer(id, false)
	}
	if removeErr != nil {
		log.Errorf("remove container error: %s \n", id)
	}
}

func (d *Docker) RmContainer(id string, force bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), DockerDeleteTimeout)
	defer cancel()
	if err := d.Client.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: force}); err != nil {
		return err
	}
	return nil
}
