package core

import (
	"context"

	"github.com/moby/moby/client"
)

type ContainerStat struct {
	ID     string
	Name   string
	Image  string
	State  string
	Status string
}

func FetchContainers() ([]ContainerStat, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var stats []ContainerStat
	for _, c := range containers.Items {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		stats = append(stats, ContainerStat{
			ID:     c.ID[:12],
			Name:   name,
			Image:  c.Image,
			State:  string(c.State),
			Status: c.Status,
		})
	}
	return stats, nil
}
