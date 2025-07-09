package step

import (
	"context"
	"fmt"
	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"studyConductor/pkg"
)

const (
	StateConfigured = "configured"
	StateExited     = "exited"
)

type ContainerStep struct {
	*AbstractStep
	podmanContext   context.Context
	containerClient *bindings.Connection
}

func boolPtr(b bool) *bool {
	return &b
}

func matchName(haystack []string, needle string) bool {
	for _, hay := range haystack {
		if strings.HasSuffix(hay, needle) {
			return true
		}
	}
	return false
}

func (s *ContainerStep) hasLocalImage(imageName string) (bool, error) {
	report, err := images.List(s.podmanContext, nil)
	if err != nil {
		return false, err
	}
	for _, image := range report {
		if matchName(image.Names, imageName) {
			return true, nil
		}
	}
	return false, nil
}

func (s *ContainerStep) Activate(ctx context.Context) error {
	s.workingMx.Store(true)
	defer s.workingMx.Store(false)
	if _, err := s.getClient(); err != nil {
		return err
	}
	imageName := s.config["IMAGE_NAME"].(string)
	isLocal, err := s.hasLocalImage(imageName)
	if err != nil {
		return err
	}
	if !isLocal {
		if _, err = images.Pull(s.podmanContext, imageName, &images.PullOptions{Quiet: boolPtr(true)}); err != nil {
			return err
		}
	}

	sg := specgen.NewSpecGenerator(imageName, false)
	sg.User = "1000:1000"
	sg.Terminal = boolPtr(false)
	sg.Name = s.config.GetString("CONTAINER_NAME")
	sg.Env = s.config.GetStringMap("ENV")
	if s.config.GetStringSlice("PORT_MAPPING") != nil {
		portMappings, err := createPortMappings(s.config.GetStringSlice("PORT_MAPPING"))
		if err != nil {
			return err
		}
		sg.PortMappings = portMappings
	}
	// Create Port Mapping
	if s.config.GetString("DATA_PATH") != "" {
		var err error
		sg.Volumes, err = s.createMounts(s.config.GetString("DATA_PATH"))
		if err != nil {
			return err
		}
		sg.Mounts = []specs.Mount{
			{
				Source:      "/etc/passwd",
				Destination: "/etc/passwd",
				Type:        "bind",
				Options:     []string{"ro"},
			},
		}
	}

	var contId string
	oldCont, err := s.getContainer()
	if err != nil {
		return err
	}
	if oldCont == nil {
		// Container does not exist yet, so create one it
		cont, err := pkg.CreateWithSpec(s.podmanContext, sg)
		if err != nil {
			fmt.Println(cont, err)
			return err
		}
		contId = cont.ID
	} else {
		contId = oldCont.ID
	}

	return pkg.Start(s.podmanContext, contId)
}

func (s *ContainerStep) createMounts(dataFolder string) ([]*specgen.NamedVolume, error) {
	path := filepath.Join(Conf.Study.Storage.GetString("path"), "_data")
	if err := os.MkdirAll(path, 0777); err != nil {
		return nil, err
	}
	volumeName, err := pkg.CreateVolume("study_storage", path)
	if err != nil {
		panic(any(err))
	}
	vol := make([]*specgen.NamedVolume, 1)
	vol[0] = &specgen.NamedVolume{
		Name: volumeName,
		Dest: dataFolder,
	}
	return vol, nil
}

func createPortMappings(ports []string) ([]types.PortMapping, error) {
	var pm []types.PortMapping
	for _, port := range ports {
		if port == "" {
			continue
		}
		mapping := strings.Split(port, ":")
		if len(mapping) != 2 {
			return nil, fmt.Errorf("invalid port map provided")
		}
		hostPort, err := strconv.Atoi(mapping[0])
		if err != nil {
			return nil, err
		}
		containerPort, err := strconv.Atoi(mapping[1])
		if err != nil {
			return nil, err
		}
		pm = append(pm, types.PortMapping{
			HostPort:      uint16(hostPort),
			ContainerPort: uint16(containerPort),
		})
	}
	return pm, nil
}

func (s *ContainerStep) Deactivate(ctx context.Context) error {
	s.workingMx.Store(true)
	defer s.workingMx.Store(false)
	if _, err := s.getClient(); err != nil {
		return err
	}
	ctr, err := s.getContainer()
	if err != nil {
		return err
	}
	if err = containers.Stop(s.podmanContext, ctr.ID, nil); err != nil {
		return err
	}
	_, err = containers.Remove(s.podmanContext, ctr.ID, nil)
	return err
}

func (s *ContainerStep) Status(ctx context.Context) (*Status, error) {
	if s.workingMx.Load() {
		s.state = StatePending
		return &Status{
			Label: SymbolWorking,
			Color: color.RGBA{R: uint8(255), G: uint8(255)},
		}, nil
	}
	if _, err := s.getClient(); err != nil {
		return nil, err
	}
	ctr, err := s.getContainer()
	if err != nil {
		return nil, err
	}
	if ctr == nil {
		s.state = StateOff
		return &Status{
			Label: SymbolCross,
			Color: color.RGBA{R: uint8(255)},
		}, nil
	}
	if ctr.State == StateConfigured || ctr.State == StateExited {
		s.state = StateOff
		return &Status{
			Label: SymbolPause,
			Color: color.RGBA{R: uint8(255), G: uint8(255)},
		}, nil
	}
	s.state = StateRunning
	return &Status{
		Label: SymbolCheckmark,
		Color: color.RGBA{G: uint8(255)},
	}, nil
}

func (s *ContainerStep) getClient() (*bindings.Connection, error) {
	var err error
	if s.containerClient == nil {
		socket := "unix:" + os.Getenv("XDG_RUNTIME_DIR") + "/podman/podman.sock"
		s.podmanContext, err = bindings.NewConnection(context.Background(), socket)
		if err != nil {
			return nil, err
		}
		s.containerClient, err = bindings.GetClient(s.podmanContext)
		if err != nil {
			return nil, err
		}
	}
	return s.containerClient, nil
}

func (s *ContainerStep) getContainer() (*entities.ListContainer, error) {
	ctrs, err := pkg.List(s.podmanContext, s.config.GetString("CONTAINER_NAME"))
	if err != nil && ctrs == nil {
		return nil, err
	}
	for _, ctr := range ctrs {
		for _, name := range ctr.Names {
			if name == s.config.GetString("CONTAINER_NAME") {
				return &ctr, nil
			}
		}
	}
	return nil, nil
}
