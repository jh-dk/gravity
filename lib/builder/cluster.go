/*
Copyright 2020 Gravitational, Inc.

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

package builder

import (
	"context"
	"io/ioutil"
	"os"

	libapp "github.com/gravitational/gravity/lib/app"
	"github.com/gravitational/gravity/lib/app/service"
	"github.com/gravitational/gravity/lib/constants"
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/loc"
	"github.com/gravitational/gravity/lib/localenv"
	"github.com/gravitational/gravity/lib/pack"
	"github.com/gravitational/gravity/lib/schema"
	"github.com/gravitational/gravity/lib/utils"

	"github.com/gravitational/trace"

	"github.com/coreos/go-semver/semver"
)

// NewClusterBuilder returns a builder that produces cluster images.
func NewClusterBuilder(config Config) (*ClusterBuilder, error) {
	engine, err := newEngine(config)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	runtimeVersions, err := parseVersions(config.UpgradeVia)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return &ClusterBuilder{
		engine:     engine,
		upgradeVia: runtimeVersions,
	}, nil
}

type ClusterBuilder struct {
	engine *Engine
	// upgradeVia lists intermediate runtime versions to embed in the resulting installer
	upgradeVia []semver.Version
}

// ClusterRequest combines parameters for building a cluster image.
type ClusterRequest struct {
	// SourcePath specifies the path to build the cluster image out of.
	SourcePath string
	// OutputPath is the resulting cluster image output file path.
	OutputPath string
	// Overwrite is whether to overwrite existing output file.
	Overwrite bool
	// Vendor combines vendoring parameters.
	Vendor service.VendorRequest
	// BaseImage is optional base image provided on the command line.
	BaseImage string
}

// Build builds a cluster image according to the provided parameters.
func (b *ClusterBuilder) Build(ctx context.Context, req ClusterRequest) error {
	imageSource, err := GetClusterImageSource(req.SourcePath)
	if err != nil {
		return trace.Wrap(err)
	}

	manifest, err := imageSource.Manifest()
	if err != nil {
		return trace.Wrap(err)
	}

	if req.BaseImage != "" {
		locator, err := loc.MakeLocator(req.BaseImage)
		if err != nil {
			return trace.Wrap(err)
		}
		manifest.SetBase(*locator)
	}

	outputPath, err := checkOutputPath(manifest, req.OutputPath, req.Overwrite)
	if err != nil {
		return trace.Wrap(err)
	}

	locator := imageLocator(manifest, req.Vendor)
	b.engine.NextStep("Building cluster image %v %v from %v", locator.Name,
		locator.Version, imageSource.Type())

	b.engine.NextStep("Selecting base image version")
	runtimeVersion, err := b.engine.SelectRuntime(manifest)
	if err != nil {
		return trace.Wrap(err)
	}
	err = b.engine.checkVersion(runtimeVersion)
	if err != nil {
		return trace.Wrap(err)
	}
	err = b.engine.SyncPackageCache(ctx, manifest, runtimeVersion, b.upgradeVia...)
	if err != nil {
		return trace.Wrap(err)
	}

	vendorDir, err := ioutil.TempDir("", "vendor")
	if err != nil {
		return trace.Wrap(err)
	}
	defer os.RemoveAll(vendorDir)

	b.engine.NextStep("Discovering and embedding Docker images")
	stream, err := b.engine.Vendor(ctx, VendorRequest{
		SourceDir: imageSource.Dir(),
		VendorDir: vendorDir,
		Manifest:  manifest,
		Vendor:    req.Vendor,
	})
	if err != nil {
		return trace.Wrap(err)
	}
	defer stream.Close()

	b.engine.NextStep("Creating application")
	application, err := b.engine.CreateApplication(stream)
	if err != nil {
		return trace.Wrap(err)
	}

	b.engine.NextStep("Packaging cluster image")
	installer, err := b.engine.GenerateInstaller(manifest, *application)
	if err != nil {
		return trace.Wrap(err)
	}
	defer installer.Close()

	b.engine.NextStep("Saving cluster image to %v", outputPath)
	err = b.engine.WriteInstaller(installer, outputPath)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// Close closes the builder
func (b *ClusterBuilder) Close() error {
	return b.engine.Close()
}

// appForRuntime builds an application object with the specified runtime version
// as the base to be able to collect dependencies of the specified base application.
func (b *ClusterBuilder) appForRuntime(runtimeVersion semver.Version) libapp.Application {
	return libapp.Application{
		Package:  b.Locator(),
		Manifest: b.Manifest.WithBase(loc.Runtime.WithVersion(runtimeVersion)),
	}
}

// collectUpgradeDependencies computes and returns a set of package dependencies for each
// configured intermediate runtime version.
// result contains combined dependencies marked with a label per runtime version.
func (b *ClusterBuilder) collectUpgradeDependencies() (result *libapp.Dependencies, err error) {
	apps, err := b.Env.AppServiceLocal(localenv.AppConfig{})
	if err != nil {
		return nil, trace.Wrap(err)
	}
	result = &libapp.Dependencies{}
	for _, runtimeVersion := range b.UpgradeVia {
		app, err := apps.GetApp(loc.Runtime.WithVersion(runtimeVersion))
		if err != nil {
			return nil, trace.Wrap(err)
		}
		req := libapp.GetDependenciesRequest{
			App:  *app,
			Apps: apps,
			Pack: b.Env.Packages,
		}
		dependencies, err := libapp.GetDependencies(req)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		dependencies.Apps = append(dependencies.Apps, *app)
		addUpgradeVersionLabel(dependencies, runtimeVersion.String())
		result.Packages = append(result.Packages, filterUpgradePackageDependencies(dependencies.Packages)...)
		result.Apps = append(result.Apps, dependencies.Apps...)
	}
	return result, nil
}

// imageLocator returns locator of the image that's being built.
func imageLocator(manifest *schema.Manifest, vendor service.VendorRequest) loc.Locator {
	name := manifest.Metadata.Name
	if vendor.PackageName != "" {
		name = vendor.PackageName
	}
	version := manifest.Metadata.ResourceVersion
	if vendor.PackageVersion != "" {
		version = vendor.PackageVersion
	}
	return loc.Locator{
		Repository: defaults.SystemAccountOrg,
		Name:       name,
		Version:    version,
	}
}

// filterUpgradePackageDependencies returns the list of package dependencies
// to include as additional dependencies when building an installer.
// packages lists all package dependencies for a specific intermediate version.
// The resulting list will only include the packages the upgrade will need
// for each intermediate hop which includes the gravity binary, teleport and planet container packages.
// All other packages are not necessary for an intermediate upgrade hop and will be omitted.
func filterUpgradePackageDependencies(packages []pack.PackageEnvelope) (result []pack.PackageEnvelope) {
	result = packages[:0]
	for _, pkg := range packages {
		if pkg.Locator.Repository != defaults.SystemAccountOrg {
			continue
		}
		switch pkg.Locator.Name {
		case constants.TeleportPackage,
			constants.GravityPackage,
			constants.PlanetPackage:
		default:
			continue
		}
		result = append(result, pkg)
	}
	return result
}

func addUpgradeVersionLabel(dependencies *libapp.Dependencies, version string) {
	for i := range dependencies.Packages {
		dependencies.Packages[i].RuntimeLabels = utils.CombineLabels(
			dependencies.Packages[i].RuntimeLabels,
			pack.RuntimeUpgradeLabels(version),
		)
	}
	for i := range dependencies.Apps {
		dependencies.Apps[i].PackageEnvelope.RuntimeLabels = utils.CombineLabels(
			dependencies.Apps[i].PackageEnvelope.RuntimeLabels,
			pack.RuntimeUpgradeLabels(version),
		)
	}
}

func parseVersions(versions []string) (result []semver.Version, err error) {
	result = make([]semver.Version, 0, len(versions))
	for _, version := range versions {
		runtimeVersion, err := semver.NewVersion(version)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		result = append(result, *runtimeVersion)
	}
	return result, nil
}
