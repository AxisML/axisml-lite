package standalone

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	cerrdefs "github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

type fakeImageEngine struct {
	inspectErr   error
	pullErr      error
	pullResponse string
	inspectCalls int
	pullCalls    int
}

func (f *fakeImageEngine) ImageInspect(
	context.Context,
	string,
	...client.ImageInspectOption,
) (image.InspectResponse, error) {
	f.inspectCalls++
	return image.InspectResponse{}, f.inspectErr
}

func (f *fakeImageEngine) ImagePull(
	context.Context,
	string,
	image.PullOptions,
) (io.ReadCloser, error) {
	f.pullCalls++
	if f.pullErr != nil {
		return nil, f.pullErr
	}
	return io.NopCloser(strings.NewReader(f.pullResponse)), nil
}

func TestEnsureImage(t *testing.T) {
	t.Run("Always pulls without inspecting", func(t *testing.T) {
		cli := &fakeImageEngine{pullResponse: `{"status":"Image is up to date"}` + "\n"}
		require.NoError(t, ensureImage(context.Background(), cli, "busybox:latest", corev1.PullAlways))
		assert.Zero(t, cli.inspectCalls)
		assert.Equal(t, 1, cli.pullCalls)
	})

	t.Run("IfNotPresent uses a local image", func(t *testing.T) {
		cli := &fakeImageEngine{}
		require.NoError(t, ensureImage(context.Background(), cli, "busybox:1.37", corev1.PullIfNotPresent))
		assert.Equal(t, 1, cli.inspectCalls)
		assert.Zero(t, cli.pullCalls)
	})

	t.Run("IfNotPresent pulls a missing image", func(t *testing.T) {
		cli := &fakeImageEngine{
			inspectErr:   cerrdefs.ErrNotFound,
			pullResponse: `{"status":"Downloaded newer image"}` + "\n",
		}
		require.NoError(t, ensureImage(context.Background(), cli, "busybox:1.37", corev1.PullIfNotPresent))
		assert.Equal(t, 1, cli.inspectCalls)
		assert.Equal(t, 1, cli.pullCalls)
	})

	t.Run("IfNotPresent propagates inspect failures", func(t *testing.T) {
		cli := &fakeImageEngine{inspectErr: errors.New("daemon unavailable")}
		err := ensureImage(context.Background(), cli, "busybox:1.37", corev1.PullIfNotPresent)
		require.ErrorContains(t, err, "inspect image")
		assert.Zero(t, cli.pullCalls)
	})

	t.Run("Never uses a local image", func(t *testing.T) {
		cli := &fakeImageEngine{}
		require.NoError(t, ensureImage(context.Background(), cli, "busybox:1.37", corev1.PullNever))
		assert.Equal(t, 1, cli.inspectCalls)
		assert.Zero(t, cli.pullCalls)
	})

	t.Run("Never rejects a missing image without pulling", func(t *testing.T) {
		cli := &fakeImageEngine{inspectErr: cerrdefs.ErrNotFound}
		err := ensureImage(context.Background(), cli, "busybox:1.37", corev1.PullNever)
		require.ErrorContains(t, err, "not present locally")
		assert.Zero(t, cli.pullCalls)
	})
}

func TestPullImageReadsDockerStreamErrors(t *testing.T) {
	cli := &fakeImageEngine{
		pullResponse: `{"errorDetail":{"message":"manifest unknown"},"error":"manifest unknown"}` + "\n",
	}
	err := pullImage(context.Background(), cli, "example.invalid/missing:latest")
	require.Error(t, err)
	assert.ErrorContains(t, err, "manifest unknown")
}

func TestPullImagePropagatesRequestErrors(t *testing.T) {
	cli := &fakeImageEngine{pullErr: errors.New("registry unavailable")}
	err := pullImage(context.Background(), cli, "busybox:latest")
	require.ErrorContains(t, err, "registry unavailable")
}
