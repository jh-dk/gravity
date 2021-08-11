/*
Copyright 2018 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

import (
	"context"
	"strings"

	apps "github.com/gravitational/gravity/lib/app"
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/docker"
	"github.com/gravitational/gravity/lib/loc"
	"github.com/gravitational/gravity/lib/pack"
	"github.com/gravitational/gravity/lib/state"
	"github.com/gravitational/gravity/lib/systemservice"
	"github.com/gravitational/gravity/lib/utils"
	"github.com/gravitational/gravity/lib/vacuum/prune"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

// New creates a new registry cleaner
func New(config Config) (*Cleanup, error) {
	if err := config.checkAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	return &Cleanup{
		config: config,
	}, nil
}

func (r *Config) checkAndSetDefaults() error {
	if r.App == nil {
		return trace.BadParameter("application package is required")
	}
	if r.ImageService == nil {
		return trace.BadParameter("docker image service is required")
	}
	if r.Packages == nil {
		return trace.BadParameter("cluster package service is required")
	}
	if r.Apps == nil {
		return trace.BadParameter("cluster application service is required")
	}
	if r.FieldLogger == nil {
		r.FieldLogger = log.WithField(trace.Component, "gc:registry")
	}
	return nil
}

// Config describes configuration for the registry cleaner
type Config struct {
	// Config specifies the common pruner configuration
	prune.Config
	// App specifies the cluster application
	App *loc.Locator
	// Packages specifies the cluster package service
	Packages pack.PackageService
	// Apps specifies the cluster application service
	Apps apps.Applications
	// ImageService specifies the docker image service
	ImageService docker.ImageService
}

// Prune removes unused docker images.
// The registry state is reset by deleting the state from the filesystem
// and re-running the docker image export for the cluster application.
func (r *Cleanup) Prune(ctx context.Context) (err error) {
	r.config.PrintStepf("Stop registry service")
	if !r.config.DryRun {
		err = r.registryStop(ctx)
		defer func() {
			if err == nil {
				return
			}
			if errStart := r.registryStart(ctx); errStart != nil {
				r.config.Warn(errStart)
			}
		}()
		if err != nil {
			return trace.Wrap(err)
		}
	}

	stateDir, err := state.GetStateDir()
	if err != nil {
		return trace.Wrap(err)
	}

	dir := state.RegistryDir(stateDir)
	r.config.PrintStepf("Delete registry state directory %v", dir)
	if !r.config.DryRun {
		err = utils.RemoveContents(dir)
		if err != nil {
			return trace.Wrap(trace.ConvertSystemError(err),
				"failed to remove old registry state from %v.", dir)
		}
	}

	r.config.PrintStepf("Start registry service")
	if !r.config.DryRun {
		err = r.registryStart(ctx)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	r.config.PrintStepf("Sync application state with registry")
	if r.config.DryRun {
		return nil
	}
	syncer := apps.Syncer{
		PackService:  r.config.Packages,
		AppService:   r.config.Apps,
		ImageService: r.config.ImageService,
	}
	err = syncer.SyncApp(ctx, *r.config.App)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

func (r *Cleanup) registryStart(ctx context.Context) error {
	out, err := r.serviceStart(ctx)
	if err != nil {
		return trace.Wrap(err, "failed to start the registry service: %s.", out)
	}

	err = r.waitForService(ctx, systemservice.ServiceStatusActive)
	if err != nil {
		return trace.Wrap(err, "failed to wait for the registry service to start")
	}
	return nil
}

func (r *Cleanup) registryStop(ctx context.Context) error {
	out, err := r.serviceStop(ctx)
	if err != nil {
		return trace.Wrap(err, "failed to stop the registry service: %s.", out)
	}

	err = r.waitForService(ctx, systemservice.ServiceStatusInactive)
	if err != nil {
		return trace.Wrap(err, "failed to wait for the registry service to stop")
	}

	return nil
}

func (r *Cleanup) waitForService(ctx context.Context, status string) error {
	localCtx, cancel := defaults.WithTimeout(ctx)
	defer cancel()
	b := utils.NewUnlimitedExponentialBackOff()
	err := utils.RetryWithInterval(localCtx, b, func() error {
		out, err := r.serviceStatus(localCtx)
		actualStatus := strings.TrimSpace(string(out))
		if strings.HasPrefix(actualStatus, status) {
			return nil
		}
		return trace.Retry(err, "unexpected service status: %s", actualStatus)
	})
	return trace.Wrap(err)
}

func (r *Cleanup) serviceStop(ctx context.Context) (output []byte, err error) {
	return serviceCtl(ctx, r.config.FieldLogger, "stop")
}

func (r *Cleanup) serviceStart(ctx context.Context) (output []byte, err error) {
	return serviceCtl(ctx, r.config.FieldLogger, "start")
}

func (r *Cleanup) serviceStatus(ctx context.Context) (output []byte, err error) {
	return serviceCtl(ctx, r.config.FieldLogger, "is-active")
}

// Cleanup implements garbage collection for docker registry
type Cleanup struct {
	// config specifies the configuration for the cleanup
	config Config
}

func serviceCtl(ctx context.Context, log log.FieldLogger, args ...string) (output []byte, err error) {
	args = append([]string{"/bin/systemctl"}, append(args, "registry.service")...)
	output, err = utils.RunCommand(ctx, log, utils.PlanetCommandArgs(args...)...)
	return output, trace.Wrap(err)
}
