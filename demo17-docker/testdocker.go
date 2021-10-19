package check

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	XtraBackup23  string = "xtrabackup2_3"     //镜像名称
	XtraBackup24  string = "xtrabackup2_4"     //镜像名称
	XtraBackup80  string = "xtrabackup8_0"     //镜像名称
	containerName string = "xtrabackup-latest" //容器名称
	indexName     string = "/" + containerName //容器索引名称，用于检查该容器是否存在是使用
	CMD           string = "/bin/bash"         //运行的cmd命令，用于启动container中的程序
	WorkDir       string = "/tmp/"             //container工作目录
	ContainerDir  string = "/home"             //容器挂载目录
	HostDir       string = "/tmp"              //容器挂在到宿主机的目录
	n             int    = 5                   //每5s检查一个容器是否在运行
)

// 获取docker镜像列表
func ListImages() []string {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	var imagesLib []string
	for _, image := range images {
		imagesLib = append(imagesLib, image.RepoTags[0])
	}
	return imagesLib
}

// 容器列表
func ListContainer() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	for _, container := range containers {
		fmt.Printf("%s\n", container.Image)
	}
}

func CreateContainer(images string) string {
	//创建容器
	//cli, err := client.NewClientWithOpts(client.FromEnv)
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	cont, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      images,
		Cmd:        []string{CMD}, //docker 容器中执行的命令
		Tty:        true,          //docker run命令中的-t选项
		OpenStdin:  true,          //docker run命令中的-i选项
		WorkingDir: WorkDir,       //docker容器中的工作目录
		//Volumes: map[string]struct{}{"/tmp":""},
	}, &container.HostConfig{ // 挂载
		//Mounts: m,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: HostDir,
				Source: "/Users/sevck/Desktop/",
				Target: ContainerDir,
			},
		},
	}, nil, nil, "")
	if err != nil {
		log.Error(err)
		//stopContainer()
		os.Exit(0)
	}
	log.Info(cont.ID)
	return cont.ID
}

//启动容器
func StartContainer(containerID string) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		log.Error("启动错误:", err)
	}
}

//停止容器
func StopContainer(containerID string) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := cli.ContainerStop(ctx, containerID, nil)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("stop", containerID)
	}
}

func RemoveContainer(containerID string) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		//RemoveLinks: true,
		//RemoveVolumes: true,
		Force: true,
	})
	if err != nil {
		log.Error(err)
	}
}

// 执行容器里面的命令 docker exec -it xxxx command
func ExecContainer(containerID string, command []string) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	execOpts := types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          command,
	}

	resp, err := cli.ContainerExecCreate(ctx, containerID, execOpts)

	if err != nil {
		log.Error(err)
	}

	attach, err := cli.ContainerExecAttach(ctx, resp.ID, types.ExecStartCheck{})
	if err != nil {
		log.Error(err)
	}
	defer attach.Close()
	data, _ := ioutil.ReadAll(attach.Reader)
	log.Info("正在解压...\n", string(data))
	//fmt.Println(string(data))
	//err = cli.ContainerExecStart(ctx, resp.ID, types.ExecStartCheck{})
	//if err != nil {
	//log.Error(err)
	//}
	//running := true
	//for running {
	//ret, err := cli.ContainerExecInspect(ctx, resp.ID)
	//if err != nil {
	//   log.Info(err)
	//}
	//if ret.Running {
	//   running = false
	//}
	//time.Sleep(250 * time.Millisecond)
	//}
}
