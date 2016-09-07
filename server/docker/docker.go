package docker

import (
	"fmt"
	"log"

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

func SetInstallationStatuses(installations []*types.Installation) error {

	for _, installation := range installations {
		err := SetInstallationStatus(installation)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetInstallationStatus(inst *types.Installation) error {
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

func RunApp(appName string, appDesc types.AppDescriptor) error {

	appContainer, err := c.ContainerCreate(context.TODO(),
		&container.Config{Image: appDesc.Services.App.Image},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		appName)

	if err != nil {
		return err
	}

	err = c.ContainerStart(context.TODO(), appContainer.ID, ctypes.ContainerStartOptions{})
	if err != nil {
		return err
	}

	webContainer, err := c.ContainerCreate(context.TODO(),
		&container.Config{Image: appDesc.Services.Web.Image},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		fmt.Sprintf("%s_web", appName))

	if err != nil {
		rerr := removeContainer(c, appContainer.ID)
		if rerr != nil {
			return err
		}
		return err
	}

	err = c.ContainerStart(context.TODO(), webContainer.ID, ctypes.ContainerStartOptions{})

	if err != nil {
		rerr := removeContainer(c, appContainer.ID)
		if rerr != nil {
			return err
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
