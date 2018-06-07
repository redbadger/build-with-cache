# Build a docker image from a multi-stage Dockerfile with layer caching

A cli command written in Go that uses a Docker registry to store layer caches in order to speed up build times. Useful in CI pipelines.

The tool parses the `Dockerfile` for the stage targets and attempts to pull respective images from the specified registry. Any images it finds are used as layer caches for the docker build. Updated images for each stage are pushed back to the registry ready for the next build.

### Steps:

1.  Parse `Dockerfile` looking for stages.
1.  Attempt to pull an image for each stage.
    1.  Image name is that supplied by `--tag` appended with stage name or number: e.g. `my-registry/my-image-builder` or `my-registry/my-image-0`
1.  For every image found, pass a `--cache-from` directive to the build.
1.  Tag the images created for each stage.
1.  Push these images back to the registry.

### Use another registry as the layer cache:

The `--cache` flag allows the specification of a separate registry to use as the layer cache. For example, `--cache localhost:5000` would mean that the images in the above example would be named like this: `localhost:5000/my-image-builder` and `localhost:5000/my-image-0`.

## To install locally:

```bash
go get github.com/redbadger/build-with-cache
```

## Usage:

The usage of `build-with-cache` is identical to [`docker build`](https://docs.docker.com/engine/reference/commandline/build/), with the following exceptions:

1.  the `--tag` flag must be specified if you want caching enabled.
1.  the `--cache` flag can be used to specify a separate registry for storing the cache images.
1.  all flags (except `--cache`) are passed unaltered to `docker build`.

## Examples:

1.  Build with local context:

    ```bash
    build-with-cache . \
      --tag=my-registry/my-image
    ```

1.  Build with tarred context on `stdin`, use a local Docker registry as the cache, and pass a `--build-arg` flag to the build:

    ```bash
    tar c .| build-with-cache - \
      --cache=localhost:5000 \
      --file=Dockerfile \
      --tag=my-registry.example.com/my-image
      --build-arg http_proxy=http://localhost:3128
    ```

## To build

1.  Sync dependencies (they have been vendored, see [`/vendor/vendor.json`](vendor/vendor.json), but not checked into github), so use [`govendor`](https://github.com/kardianos/govendor) to sync:

    ```bash
    govendor sync
    ```

1.  Build everything in your repository only: `govendor install +local` or test your repository only: `govendor test +local`
