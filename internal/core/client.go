package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/wibu-gaptek/oppai-desu/internal/utils"
)

type DockerOppai struct {
	ctx    context.Context
	client *client.Client

	registryAuthString string
	registryAuthMap    map[string]registry.AuthConfig
}

func NewOppai(ctx context.Context, cli *client.Client) DockerOppai {
	return DockerOppai{
		ctx:    ctx,
		client: cli,
	}
}

func (d *DockerOppai) Login(username, password, server string) string {
	authConfig := registry.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: server,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(encodedJSON)

}

func (d *DockerOppai) ImagesList(opt image.ListOptions) ([]image.Summary, error) {
	img, err := d.client.ImageList(d.ctx, opt)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (d *DockerOppai) ImagesHistory(image string) ([]image.HistoryResponseItem, error) {
	img, err := d.client.ImageHistory(d.ctx, image)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (d *DockerOppai) ImagesBuild(opt types.ImageBuildOptions, src, dockerFile string, f func(rd io.Reader) error) error {
	tar, err := utils.TarWithOpt(src)
	if err != nil {
		return err
	}

	res, err := d.client.ImageBuild(d.ctx, tar, opt)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Hook function
	return f(res.Body)
}

func (d *DockerOppai) ImagesPush(image string, opt image.PushOptions, f func(rd io.Reader) error) error {

	res, err := d.client.ImagePush(d.ctx, image, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	// Hook function
	return f(res)
}

func (d *DockerOppai) ImagesPull(ref string, opt image.PullOptions, f func(rd io.Reader) error) error {

	res, err := d.client.ImagePull(d.ctx, ref, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	// Hook function
	return f(res)
}

func (d *DockerOppai) ImagesImport(ref string, src image.ImportSource, opt image.ImportOptions, f func(rd io.Reader) error) error {

	res, err := d.client.ImageImport(d.ctx, src, ref, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	return f(res)
}
