package pkg

import (
	"context"
	"encoding/json"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/specgen"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func CreateWithSpec(ctx context.Context, s *specgen.SpecGenerator) (entities.ContainerCreateResponse, error) {
	var ccr entities.ContainerCreateResponse
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return ccr, err
	}
	specgenString, err := json.Marshal(s)
	if err != nil {
		return ccr, err
	}
	stringReader := strings.NewReader(string(specgenString))
	response, err := conn.DoRequest(ctx, stringReader, http.MethodPost, "/containers/create", nil, nil)
	if err != nil {
		return ccr, err
	}
	return ccr, response.Process(&ccr)
}

func CreateVolume(name, path string) (string, error) {
	_, err := exec.Command("podman", "volume", "create", "--ignore", "-o", "device="+path, "-o", "o=bind", name).Output()
	return name, err
}

// List obtains a list of containers in local storage.  All parameters to this method are optional.
// The filters are used to determine which containers are listed. The last parameter indicates to only return
// the most recent number of containers.  The pod and size booleans indicate that pod information and rootfs
// size information should also be included.  Finally, the sync bool synchronizes the OCI runtime and
// container state.
func List(ctx context.Context, name string) ([]entities.ListContainer, error) { // nolint:typecheck
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	var containers []entities.ListContainer
	params := url.Values{}
	params.Set("limit", "1")
	filter, err := json.Marshal(map[string][]string{"name": {name}})
	if err != nil {
		return nil, err
	}
	params.Set("filters", string(filter))
	response, err := conn.DoRequest(ctx, nil, http.MethodGet, "/containers/json", params, nil)
	if err != nil {
		return containers, err
	}
	return containers, response.Process(&containers)
}

// Start starts a non-running container.The nameOrID can be a container name
// or a partial/full ID. The optional parameter for detach keys are to override the default
// detach key sequence.
func Start(ctx context.Context, nameOrID string) error {
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return err
	}
	response, err := conn.DoRequest(ctx, nil, http.MethodPost, "/containers/%s/start", url.Values{}, nil, nameOrID)
	if err != nil {
		return err
	}
	return response.Process(nil)
}
