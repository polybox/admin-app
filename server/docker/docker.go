package docker

import (
	"fmt"
	"log"
	"os"

	"github.com/docker/engine-api/client"
	ctypes "github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/network"
	"github.com/mobyos/admin-app/server/types"
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

func SetInstallationStatuses(installations []*types.Application) error {

	for _, installation := range installations {
		err := SetInstallationStatus(installation)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetInstallationStatus(inst *types.Application) error {
	container, err := c.ContainerInspect(context.TODO(), inst.Id)

	if err != nil {
		if client.IsErrNotFound(err) {
			inst.Status = "Not running"
		} else {
			return err
		}
	} else {
		inst.Status = container.State.Status
	}
	return nil
}

func createAndStart(appName string, process types.Process) (string, error) {
	cconfig := &container.Config{Image: process.Image}
	hconfig := &container.HostConfig{}
	if process.Ui {
		hconfig.Binds = []string{"/tmp/.X11-unix/:/tmp/.X11-unix/"}
		cconfig.Env = []string{fmt.Sprintf("DISPLAY=unix%s", os.Getenv("DISPLAY"))}
	}

	if process.Sound {
		hconfig.Devices = []container.DeviceMapping{{"/dev/snd", "/dev/snd", "rwm"}}
		hconfig.Devices = []container.DeviceMapping{{"/dev/dri", "/dev/dri", "rwm"}}
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

func RunApp(appName string, appDesc types.AppDescriptor) error {

	appId, err := createAndStart(appName, appDesc.Services.App)
	if err != nil {
		return err
	}

	_, err = createAndStart(fmt.Sprintf("%s_web", appName), appDesc.Services.Remote)

	if err != nil {
		rerr := removeContainer(c, appId)
		if rerr != nil {
			return rerr
		}
		return err
	}

	return nil

}

func StopApp(appName string) error {

	// try to remove both containers even though the first remove raises an error
	err := c.ContainerRemove(context.TODO(), appName, ctypes.ContainerRemoveOptions{Force: true})
	errWeb := c.ContainerRemove(context.TODO(), fmt.Sprintf("%s_web", appName), ctypes.ContainerRemoveOptions{Force: true})

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
