package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func main() {
	err := runContainer("alpine")
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func runContainer(imageRef string) error {
	ctx := context.Background()

	dc, err := client.NewClientWithOpts()
	if err != nil {
		return errors.Wrap(err, "error creating docker client")
	}

	reader, err := dc.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrapf(err, "error pulling docker image %s", imageRef)
	}
	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return errors.Wrap(err, "error reading docker pull output")
	}

	resp, err := dc.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
	}, nil, nil, "")
	if err != nil {
		return errors.Wrap(err, "error creating container")
	}

	err = dc.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "error starting container")
	}

	return nil
}
