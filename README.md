# oppai-desu
Docker orchestrator written in Go.

### `DockerOppai` Struct

The `DockerOppai` struct is a wrapper around the Docker client, with additional fields for context, the client itself, and registry authentication details.

- **`ctx`**: This is the `context.Context` used for managing cancellation and deadlines for operations.
- **`client`**: This is the Docker client used to interact with the Docker engine.
- **`registryAuthString`**: A string for registry authentication, although it's not used in the provided code.
- **`registryAuthMap`**: A map for storing registry authentication configurations, but not used in the provided code.

### `NewOppai` Function

```go
func NewOppai(ctx context.Context, cli *client.Client) DockerOppai {
	return DockerOppai{
		ctx:    ctx,
		client: cli,
	}
}
```

**Parameters**:
- **`ctx`**: The context to use for Docker operations. This helps manage request timeouts and cancellations.
- **`cli`**: An instance of `*client.Client` which is used to interact with the Docker daemon.

**Returns**:
- An instance of `DockerOppai` initialized with the provided context and Docker client.

### `Login` Function

```go
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
```

**Parameters**:
- **`username`**: Docker registry username.
- **`password`**: Docker registry password.
- **`server`**: Docker registry server address.

**Returns**:
- A base64-encoded JSON string of the authentication configuration.

### `ImagesList` Function

```go
func (d *DockerOppai) ImagesList(opt image.ListOptions) ([]image.Summary, error) {
	img, err := d.client.ImageList(d.ctx, opt)
	if err != nil {
		return nil, err
	}
	return img, nil
}
```

**Parameters**:
- **`opt`**: An `image.ListOptions` struct that contains options for filtering the image list.

**Returns**:
- A slice of `image.Summary` structs representing the images.
- An error if the operation fails.

### `ImagesHistory` Function

```go
func (d *DockerOppai) ImagesHistory(image string) ([]image.HistoryResponseItem, error) {
	img, err := d.client.ImageHistory(d.ctx, image)
	if err != nil {
		return nil, err
	}
	return img, nil
}
```

**Parameters**:
- **`image`**: The name of the image to retrieve history for.

**Returns**:
- A slice of `image.HistoryResponseItem` representing the image history.
- An error if the operation fails.

### `ImagesBuild` Function

```go
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
```

**Parameters**:
- **`opt`**: `types.ImageBuildOptions` specifying build options such as tags and Dockerfile context.
- **`src`**: The source directory to be included in the build context.
- **`dockerFile`**: Path to the Dockerfile (though not used directly in the provided code).
- **`f`**: A function that takes an `io.Reader` and processes it. This allows custom handling of the build output.

**Returns**:
- An error if the build operation fails.

### `ImagesPush` Function

```go
func (d *DockerOppai) ImagesPush(image string, opt image.PushOptions, f func(rd io.Reader) error) error {
	res, err := d.client.ImagePush(d.ctx, image, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	// Hook function
	return f(res)
}
```

**Parameters**:
- **`image`**: The image name to push to the Docker registry.
- **`opt`**: `image.PushOptions` struct that contains options for pushing the image.
- **`f`**: A function that takes an `io.Reader` and processes it. This allows custom handling of the push output.

**Returns**:
- An error if the push operation fails.

### `ImagesPull` Function

```go
func (d *DockerOppai) ImagesPull(ref string, opt image.PullOptions, f func(rd io.Reader) error) error {
	res, err := d.client.ImagePull(d.ctx, ref, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	// Hook function
	return f(res)
}
```

**Parameters**:
- **`ref`**: Reference to the image to pull.
- **`opt`**: `image.PullOptions` struct that contains options for pulling the image.
- **`f`**: A function that takes an `io.Reader` and processes it. This allows custom handling of the pull output.

**Returns**:
- An error if the pull operation fails.

### `ImagesImport` Function

```go
func (d *DockerOppai) ImagesImport(ref string, src image.ImportSource, opt image.ImportOptions, f func(rd io.Reader) error) error {
	res, err := d.client.ImageImport(d.ctx, src, ref, opt)
	if err != nil {
		return err
	}
	defer res.Close()

	return f(res)
}
```

**Parameters**:
- **`ref`**: Reference for the imported image.
- **`src`**: `image.ImportSource` representing the source of the import, such as a tarball.
- **`opt`**: `image.ImportOptions` struct for import options.
- **`f`**: A function that takes an `io.Reader` and processes it. This allows custom handling of the import output.

**Returns**:
- An error if the import operation fails.

In summary, these functions wrap Docker client operations and provide hooks for custom processing of the Docker API responses. The `context.Context` and `*client.Client` are used to interact with Docker, while the other parameters help control the specific details of each Docker operation.