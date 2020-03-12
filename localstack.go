package localstack

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type LocalStack struct {
	dockerClient *client.Client
}

const abstractionIp = "0.0.0.0"

func New() (*LocalStack, error) {

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	return &LocalStack{dockerClient: dockerClient}, nil
}

func (l *LocalStack) Create(ctx context.Context) error {

	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"8080/tcp": {
				{HostIP: abstractionIp, HostPort: "8080"},
			},
			"4572/tcp": {
				{HostIP: abstractionIp, HostPort: "4572"},
			},
			"4576/tcp": {
				{HostIP: abstractionIp, HostPort: "4576"},
			},
		},
	}

	cont, err := l.dockerClient.ContainerCreate(ctx, &container.Config{
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env: []string{"DOCKER_HOST=unix:///var/run/docker.sock", "SERVICES=sqs,s3",
			"DEFAULT_REGION=us-east-1", "DATA_DIR=/tmp/localstack/data", "USE_SSL=false"},
		Image: "localstack/localstack",
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
			"4572/tcp": struct{}{},
			"4576/tcp": struct{}{},
		},
		Volumes: map[string]struct{}{
			"/var/run/docker.sock:/var/run/docker.sock": {},
			"/private${TMPDIR}:/tmp/localstack":         {},
		},
	}, hostConfig, nil, "localstack")

	if err != nil {
		return err
	}

	return save(cont.ID)
}

func (l *LocalStack) Start(ctx context.Context) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	if err := l.dockerClient.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// TODO: workaround
	fmt.Println("Waiting for services stay up!!")
	time.Sleep(30 * time.Second)

	return nil
}

func (l *LocalStack) Stop(ctx context.Context, timeout *time.Duration) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	return l.dockerClient.ContainerStop(ctx, containerId, timeout)
}

func (l *LocalStack) Remove(ctx context.Context, removeVolumes bool) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	return l.dockerClient.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: removeVolumes,
	})
}

func (l *LocalStack) Pull(ctx context.Context) error {

	output, err := l.dockerClient.ImagePull(ctx, "docker.io/localstack/localstack", types.ImagePullOptions{})

	if err != nil {
		return err
	}

	defer output.Close()

	_, err = io.Copy(os.Stdout, output)

	if err != nil {
		return err
	}

	return nil
}

func (l *LocalStack) Run(ctx context.Context) error {

	exist, err := l.ContainerExists(ctx)

	if err != nil {
		return err
	}

	// TODO: create a method
	if !exist {

		err := l.Create(ctx)

		if !client.IsErrImageNotFound(err) {

			return err
		} else {

			if err := l.Pull(ctx); err != nil {
				return err
			}

			if err := l.Create(ctx); err != nil {
				return err
			}
		}
	}

	isRunning, err := l.IsRunning(ctx)

	if err != nil {
		return err
	}

	if !isRunning {
		return l.Start(ctx)
	}

	return nil
}

func (l *LocalStack) Logs(ctx context.Context) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	logs, err := l.dockerClient.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})

	if err != nil {
		return err
	}

	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)

	if err != nil {
		return err
	}

	return nil
}

func (l *LocalStack) ContainerExists(ctx context.Context) (bool, error) {

	containers, err := l.dockerClient.ContainerList(ctx, types.ContainerListOptions{All: true})

	if err != nil {
		return false, err
	}

	for _, cont := range containers {

		if strings.EqualFold(cont.Names[0], "/localstack") {

			if err := save(cont.ID); err != nil {
				return false, err
			}

			return true, nil
		}
	}

	return false, nil
}

func (l *LocalStack) IsRunning(ctx context.Context) (bool, error) {

	containerId, err := load()

	if err != nil {
		return false, err
	}

	containerInspected, err := l.dockerClient.ContainerInspect(ctx, containerId)

	if err != nil {
		return false, err
	}

	if containerInspected.State.Status != "running" {
		return false, nil
	}

	return true, nil
}

func save(containerId string) error {

	err := ioutil.WriteFile("localstack.out", []byte(strings.TrimSpace(containerId)), os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func load() (string, error) {

	file, err := ioutil.ReadFile("localstack.out")

	if err != nil {

		if errors.Is(err, os.ErrNotExist) {

			_, err := New()

			if err != nil {
				return "", err
			}

		}

		return "", err
	}

	return string(file), nil
}
