package docker

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"

	ctypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/mobyos/mobyos-admin-app/server/types"
	"golang.org/x/net/context"
)

var c *client.Client

func init() {
	var err error
	c, err = client.NewEnvClient()
	if err != nil {
		// this wont happen if daemon is offline, only for some critical errors
		log.Fatal("Cannot initialize docker client")
	}

}

func SetContainerStates(installations []*types.Application) error {

	for _, installation := range installations {
		err := SetContainerState(installation)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetContainerState(inst *types.Application) error {
	container, err := c.ContainerInspect(context.TODO(), inst.Id)

	if err != nil {
		if !client.IsErrNotFound(err) {
			return err
		}
	} else {
		inst.IsRunning = container.State.Running
		for _, v := range container.NetworkSettings.Ports {
			inst.RemotePort = v[0].HostPort
			break
		}
	}
	return nil
}

func createAndStart(appName string, process types.Process) (string, error) {
	cconfig := &container.Config{Image: process.Image, Cmd: process.Command, ExposedPorts: nat.PortSet{}}
	hconfig := &container.HostConfig{PublishAllPorts: true}

	for _, portNum := range process.Ports {
		port, err := nat.NewPort("tcp", portNum)
		if err != nil {
			return "", err
		}
		cconfig.ExposedPorts[port] = struct{}{}
	}

	if process.Ui {
		hconfig.Binds = []string{"/tmp/.X11-unix/:/tmp/.X11-unix/"}
		cconfig.Env = []string{fmt.Sprintf("DISPLAY=unix%s", os.Getenv("DISPLAY"))}
	}

	if process.Sound {
		hconfig.Devices = []container.DeviceMapping{{"/dev/snd", "/dev/snd", "rwm"}}
		hconfig.Devices = []container.DeviceMapping{{"/dev/video0", "/dev/video0", "rwm"}}
		//hconfig.Devices = []container.DeviceMapping{{"/dev/vchiq", "/dev/vchiq", "rwm"}}
		hconfig.Devices = []container.DeviceMapping{{"/dev/dri", "/dev/dri", "rwm"}}
		hconfig.Binds = append(hconfig.Binds, "/run/user/1000/pulse:/run/pulse:ro")
	}

	if process.Input {
		hconfig.Binds = []string{"/dev/input:/dev/input"}
		cconfig.Tty = true
	}

	hash := fnv.New32a()
	for _, volume := range process.Volumes {
		hash.Write([]byte(volume))
		volumeHash := hash.Sum32()
		hconfig.Binds = append(hconfig.Binds, fmt.Sprintf("%s_%s:%s", appName, strconv.Itoa(int(volumeHash)), volume))
		hash.Reset()
	}

	container, err := c.ContainerCreate(context.TODO(),
		cconfig,
		hconfig,
		&network.NetworkingConfig{},
		appName)

	if err != nil {
		return "", err
	}

	err = c.ContainerStart(context.TODO(), container.ID, ctypes.ContainerStartOptions{})
	if err != nil {
		rerr := removeContainer(c, container.ID)
		if rerr != nil {
			return "", rerr
		}
		return "", err
	}

	return container.ID, nil
}

func RunApp(app *types.Application) error {

	appId, err := createAndStart(app.Id, app.Descriptor.Services.App)
	if err != nil {
		return err
	}

	if app.Descriptor.Services.Remote.Image != "" {
		_, err = createAndStart(fmt.Sprintf("%s_web", app.Id), app.Descriptor.Services.Remote)

		if err != nil {
			rerr := removeContainer(c, appId)
			if rerr != nil {
				return rerr
			}
			return err
		}

	}

	// get app status before returning it to the user
	err = SetContainerState(app)
	if err != nil {
		return err
	}

	return nil

}

func StopApp(app *types.Application) error {

	// try to remove both containers even though the first remove raises an error
	err := c.ContainerRemove(context.TODO(), app.Id, ctypes.ContainerRemoveOptions{Force: true})

	var errWeb error
	if app.Descriptor.Services.Remote.Image != "" {
		errWeb = c.ContainerRemove(context.TODO(), fmt.Sprintf("%s_web", app.Id), ctypes.ContainerRemoveOptions{Force: true})
	}

	if err != nil {
		return err
	} else if errWeb != nil {
		return errWeb
	}

	return nil

}

func removeContainer(c *client.Client, id string) error {
	return c.ContainerRemove(context.TODO(), id, ctypes.ContainerRemoveOptions{Force: true})
}
