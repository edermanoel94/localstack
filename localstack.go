package localstack

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

type LocalStack struct {
	client   *client.Client
	services []Service

	// dumpingContainer dump all config in a json file
	dumpingContainer bool
}

const abstractIP = "0.0.0.0"

func New(makeInspect bool, services ...Service) (*LocalStack, error) {

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		return nil, err
	}

	if len(services) == 0 {
		services = all
	}

	return &LocalStack{
		client:           dockerClient,
		services:         services,
		dumpingContainer: makeInspect,
	}, nil
}

func (l *LocalStack) create(ctx context.Context) error {

	portBindings := l.mountPortBindings()

	hostConfig := l.mountHostConfig(portBindings)

	servicesEnv := l.mountServicesEnv()

	exposedPorts := l.mountExposedPorts()

	volumes := l.mountVolumes()

	cont, err := l.client.ContainerCreate(ctx, &container.Config{
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		Env: []string{"DOCKER_HOST=unix:///var/run/docker.sock", servicesEnv,
			"DEFAULT_REGION=us-east-1", "DATA_DIR=/tmp/localstack/data", "USE_SSL=false"},
		Image:        "localstack/localstack",
		ExposedPorts: exposedPorts,
		Volumes:      volumes,
	}, hostConfig, nil, "localstack")

	if err != nil {
		return err
	}

	return save(cont.ID)
}

func (l *LocalStack) start(ctx context.Context) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	if err := l.client.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
		return err
	}

	// TODO: this is a workaround
	fmt.Println("Waiting for services stay up!!")
	time.Sleep(30 * time.Second)

	return nil
}

func (l *LocalStack) pull(ctx context.Context) error {

	output, err := l.client.ImagePull(ctx, "docker.io/localstack/localstack", types.ImagePullOptions{})

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

func (l *LocalStack) containerExists(ctx context.Context) (bool, error) {

	containers, err := l.client.ContainerList(ctx, types.ContainerListOptions{All: true})

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

func (l *LocalStack) isRunning(ctx context.Context) (bool, error) {

	containerInspected, err := l.inspect(ctx)

	if err != nil {
		return false, err
	}

	if containerInspected.State.Status != "running" {
		return false, nil
	}

	return true, nil
}

func (l *LocalStack) inspect(ctx context.Context) (types.ContainerJSON, error) {

	containerId, err := load()

	if err != nil {
		return types.ContainerJSON{}, err
	}

	containerInspected, err := l.client.ContainerInspect(ctx, containerId)

	if err != nil {
		return types.ContainerJSON{}, err
	}

	return containerInspected, nil
}

func (l *LocalStack) logs(ctx context.Context) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	logs, err := l.client.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{
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

func (l *LocalStack) dumpingInspectedContainer(ctx context.Context) error {

	containerJson, err := l.inspect(ctx)

	if err != nil {
		return err
	}

	containerJsonBytes, err := json.Marshal(&containerJson)

	if err != nil {
		return err
	}

	return ioutil.WriteFile("dumping_inspected.json", containerJsonBytes, os.ModePerm)
}

func (l *LocalStack) mountHostConfig(portBindings map[nat.Port][]nat.PortBinding) *container.HostConfig {

	networkMode := "default"

	if runtime.GOOS == "linux" {
		networkMode = "host"
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		NetworkMode:  container.NetworkMode(networkMode),
	}

	return hostConfig
}

func (l *LocalStack) mountServicesEnv() string {

	services := strings.Builder{}

	services.WriteString("SERVICES=")

	for _, service := range l.services {

		// Ignore admin service
		if strings.EqualFold(service.Name(), "admin") {
			continue
		}

		services.WriteString(fmt.Sprintf("%s,", service.Name()))
	}

	result := services.String()

	suffix := ","

	if strings.HasSuffix(result, suffix) {
		return result[:len(result)-len(suffix)]
	}

	return result
}

func (l *LocalStack) mountExposedPorts() nat.PortSet {

	exposedPorts := make(nat.PortSet)

	for _, service := range l.services {
		exposedPorts[concatWithTCP(service.NatPort())] = struct{}{}
	}

	return exposedPorts
}

func (l *LocalStack) mountPortBindings() map[nat.Port][]nat.PortBinding {

	portBindings := make(map[nat.Port][]nat.PortBinding)

	for _, service := range l.services {
		portBindings[concatWithTCP(service.NatPort())] = []nat.PortBinding{
			{HostIP: abstractIP, HostPort: service.NatPort().Port()},
		}
	}

	return portBindings
}

func (l *LocalStack) mountVolumes() map[string]struct{} {

	volumes := make(map[string]struct{})

	volumes["/var/run/docker.sock:/var/run/docker.sock"] = struct{}{}

	if runtime.GOOS == "darwin" {
		volumes["/private${TMPDIR}:/tmp/localstack"] = struct{}{}
	}

	return volumes
}

func (l *LocalStack) Stop(ctx context.Context, timeout *time.Duration) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	return l.client.ContainerStop(ctx, containerId, timeout)
}

func (l *LocalStack) Remove(ctx context.Context, removeLinks, removeVolumes bool) error {

	containerId, err := load()

	if err != nil {
		return err
	}

	return l.client.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   removeLinks,
		RemoveVolumes: removeVolumes,
	})
}

func (l *LocalStack) Run(ctx context.Context) error {

	exist, err := l.containerExists(ctx)

	if err != nil {
		return err
	}

	// TODO: create a method
	if !exist {

		// try to create container if have a image
		err := l.create(ctx)

		// if not have image, download and then try create container again and if successful, start!
		if !client.IsErrImageNotFound(err) {

			return err
		} else {

			if err := l.pull(ctx); err != nil {
				return err
			}
		}

		// retry to create
		if err := l.create(ctx); err != nil {
			return err
		}

		return l.start(ctx)
	}

	isRunning, err := l.isRunning(ctx)

	if err != nil {
		return err
	}

	if l.dumpingContainer {
		if err := l.dumpingInspectedContainer(ctx); err != nil {
			return err
		}
	}

	if !isRunning {
		return l.start(ctx)
	}

	return nil
}

func save(containerId string) error {

	err := ioutil.WriteFile("localstack.out", []byte(strings.TrimSpace(containerId)), os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func load() (string, error) {

	// TODO: check if delete is gonna be a problem
	file, err := ioutil.ReadFile("localstack.out")

	if err != nil {
		return "", err
	}

	return string(file), nil
}
