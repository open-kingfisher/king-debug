package main

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"kingfisher/kf/common/log"
	"net/http"
)

var (
	ctx            context.Context
	cli            *client.Client
	DockerEndpoint = "unix:///var/run/docker.sock"
)

func main() {
	dockerClient()
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err, _ := w.Write([]byte("king ok"))
		log.Error("Response write error: ", err)
	})
	http.HandleFunc("/terminal", terminal)
	log.Info("Starting king debug ...")
	log.Fatal(http.ListenAndServe(":9091", nil))
}

func dockerClient() {
	ctx = context.Background()
	newCli, err := client.NewClient(DockerEndpoint, "", nil, nil)
	cli = newCli
	if err != nil {
		log.Infof("Failed to get docker client: ", err)
		return
	}

	cli.NegotiateAPIVersion(ctx)
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func terminal(w http.ResponseWriter, r *http.Request) {
	// http协议升级为websocket
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	// 通过/etc/DEBUG_CONTAINER_ID获取容器ID
	containerId, err := ioutil.ReadFile("/etc/DEBUG_CONTAINER_ID")
	if err != nil {
		log.Infof("Failed to get debug container Id: ", err)
		return
	}
	// 执行exec，获取到容器终端的连接
	log.Infof("Container Id: %s", containerId)
	hr, err := exec(string(containerId), "/bin")
	if err != nil {
		log.Error("Failed to exec container: ", err)
		return
	}
	// 关闭I/O流
	defer hr.Close()
	// 退出进程
	defer func() {
		hr.Conn.Write([]byte("exit\r"))
	}()

	go func() {
		wsWriterCopy(hr.Conn, conn)
	}()
	wsReaderCopy(conn, hr.Conn)
}

func exec(container string, workdir string) (hr types.HijackedResponse, err error) {
	// TODO cmd 应为环境变量里面的ENTRY_POINT值
	// 执行/bin/bash命令
	ir, err := cli.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"/bin/bash"},
		Tty:          true,
	})
	if err != nil {
		log.Error("Failed to exec /bin/bash command: ", err)
		return
	}

	// 附加到上面创建的/bin/bash进程中
	hr, err = cli.ContainerExecAttach(ctx, ir.ID, types.ExecStartCheck{Detach: false, Tty: true})
	if err != nil {
		log.Error("Failed to attach /bin/bash: ", err)
		return
	}
	return
}

func wsWriterCopy(reader io.Reader, writer *websocket.Conn) {
	buf := make([]byte, 8192)
	for {
		nr, err := reader.Read(buf)
		if nr > 0 {
			err := writer.WriteMessage(websocket.BinaryMessage, buf[0:nr])
			if err != nil {
				log.Error("Failed to write message: ", err)
				return
			}
		}
		if err != nil {
			log.Error("Failed to read buffer: ", err)
			return
		}
	}
}

func wsReaderCopy(reader *websocket.Conn, writer io.Writer) {
	for {
		messageType, p, err := reader.ReadMessage()
		if err != nil {
			log.Error("Failed to read message: ", err)
			return
		}
		if messageType == websocket.TextMessage {
			writer.Write(p)
		}
	}
}
