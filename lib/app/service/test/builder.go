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

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gravitational/gravity/lib/app"
	"github.com/gravitational/gravity/lib/archive"
	"github.com/gravitational/gravity/lib/loc"
	"github.com/gravitational/gravity/lib/pack"
	"github.com/gravitational/gravity/lib/schema"

	check "gopkg.in/check.v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PackageRequest describes an intent to create a test package
type PackageRequest struct {
	// Packages specifies the package service where the package is to be created
	Packages pack.PackageService
	// Package describes the package to create
	Package Package
}

// AppRequest describes an intent to create a test application package
type AppRequest struct {
	// Packages specifies the package service where packages will be created
	Packages pack.PackageService
	// Apps specifies the application service where the package is to be created.
	Apps app.Applications
	// App defines the application package to create
	App App
}

// Dependencies groups package/application dependencies
type Dependencies struct {
	Packages []Package
	Apps     []App
}

// Package describes a test package
type Package struct {
	// Loc identifies the package to create
	Loc loc.Locator
	// Labels optionally specifies package labels
	Labels map[string]string
	// Items optionally specifies the contents of the package.
	Items []*archive.Item
}

func (r *App) WithDependencies(deps Dependencies) *App {
	r.Dependencies = deps
	return r
}

func (r *App) WithSchemaDependencies(deps schema.Dependencies) *App {
	r.Manifest.Dependencies = deps
	return r
}

func (r *App) WithSchemaPackageDependencies(deps ...loc.Locator) *App {
	packages := make([]Package, 0, len(deps))
	schemaPackages := make([]schema.Dependency, 0, len(deps))
	for _, pkg := range deps {
		packages = append(packages, Package{Loc: pkg})
		schemaPackages = append(schemaPackages, schema.Dependency{Locator: pkg})
	}
	r.Dependencies.Packages = packages
	r.Manifest.Dependencies.Packages = schemaPackages
	return r
}

func (r *App) WithAppDependencies(deps ...App) *App {
	schemaApps := make([]schema.Dependency, 0, len(deps))
	for _, app := range deps {
		schemaApps = append(schemaApps, schema.Dependency{Locator: app.Manifest.Locator()})
	}
	r.Dependencies.Apps = deps
	r.Manifest.Dependencies.Apps = schemaApps
	return r
}

func (r *App) WithPackageDependencies(deps ...Package) *App {
	schemaPackages := make([]schema.Dependency, 0, len(deps))
	for _, pkg := range deps {
		schemaPackages = append(schemaPackages, schema.Dependency{Locator: pkg.Loc})
	}
	r.Dependencies.Packages = deps
	r.Manifest.Dependencies.Packages = schemaPackages
	return r
}

func (r *App) WithItems(items ...*archive.Item) *App {
	r.Items = items
	return r
}

func (r App) Locator() loc.Locator {
	return r.Manifest.Locator()
}

func (r App) Build() App { return r }

// App describes a test application package
type App struct {
	// Manifest describes the application to create
	Manifest schema.Manifest
	// Base describes the base (runtime) application
	Base *App
	// Labels optionally specifies application package labels
	Labels map[string]string
	// Items optionally specifies the contents of the application package.
	Items []*archive.Item
	// Dependencies optionally specify the application's dependencies.
	// These override dependencies from the manifest with actual data.
	Dependencies Dependencies
}

// CreateApplication creates a new test application as described by the given request
func CreateApplication(req AppRequest, c *check.C) (app *app.Application) {
	pkgDeps := make(map[loc.Locator]Package)
	appDeps := make(map[loc.Locator]App)
	if req.App.Base != nil {
		collectBaseDependencies(*req.App.Base, pkgDeps, appDeps, c)
	}
	// collectManifestDependencies(req.App.Manifest, pkgDeps, appDeps)
	collectDependencies(req.App, pkgDeps, appDeps)
	// // override with dependencies from the configuration
	// for _, d := range req.App.Dependencies.Packages {
	// 	pkgDeps[d.Loc] = d
	// }
	// for _, d := range req.App.Dependencies.Apps {
	// 	appDeps[d.Manifest.Locator()] = d
	// }
	for _, pkg := range pkgDeps {
		CreatePackage(PackageRequest{
			Package:  pkg,
			Packages: req.Packages,
		}, c)
	}
	for _, app := range appDeps {
		fmt.Println("Create application ", app.Manifest.Locator())
		data := CreatePackageData(ApplicationLayout(app, c), c)
		_, err := req.Apps.CreateApp(app.Manifest.Locator(), &data, app.Labels)
		c.Assert(err, check.IsNil)
	}

	fmt.Println("Create application ", req.App.Manifest.Locator())
	data := CreatePackageData(ApplicationLayout(req.App, c), c)
	app, err := req.Apps.CreateApp(req.App.Manifest.Locator(), &data, req.App.Labels)
	c.Assert(err, check.IsNil)

	return app
}

// CreatePackage creates a new test package as described by the given request
func CreatePackage(req PackageRequest, c *check.C) *pack.PackageEnvelope {
	fmt.Println("Create package ", req.Package.Loc)
	items := req.Package.Items
	if len(items) == 0 {
		// Create a package with a test payload
		items = append(items, archive.ItemFromString("data", req.Package.Loc.String()))
	}
	input := CreatePackageData(items, c)
	c.Assert(req.Packages.UpsertRepository(req.Package.Loc.Repository, time.Time{}), check.IsNil)

	pkg, err := req.Packages.CreatePackage(req.Package.Loc, &input, pack.WithLabels(req.Package.Labels))
	c.Assert(err, check.IsNil)
	c.Assert(pkg, check.NotNil)

	return pkg
}

// CreateDummyApplication creates an application with a valid manifest, but fake content.
// It returns the application created in the last service specified with services
func CreateDummyApplication(locator loc.Locator, c *check.C, services ...app.Applications) (result *app.Application) {
	for _, s := range services {
		result = CreateApplication(AppRequest{
			App:  DefaultClusterApplication(locator).Build(),
			Apps: s,
		}, c)
	}
	return result
}

var (
	// RuntimeApplicationLoc specifies the default runtime application locator
	RuntimeApplicationLoc = loc.MustParseLocator("gravitational.io/kubernetes:0.0.1")
	// RuntimePackageLoc specifies the default runtime package locator
	RuntimePackageLoc = loc.MustParseLocator("gravitational.io/planet:0.0.1")
)

// NewDependency is a convenience helper to create a manifest Dependency from a package locator
func NewDependency(pkgLoc string) schema.Dependency {
	return schema.Dependency{
		Locator: loc.MustParseLocator(pkgLoc),
	}
}

// DefaultRuntimeApplication returns a default test runtime application manifest
func DefaultRuntimeApplication() *App {
	return RuntimeApplication(RuntimeApplicationLoc, RuntimePackageLoc)
}

// RuntimeApplication returns a test runtime application manifest
// given the application locator and the locator for the runtime (planet) package
func RuntimeApplication(appLoc, runtimePackageLoc loc.Locator) *App {
	return &App{
		Manifest: schema.Manifest{
			Header: schema.Header{
				TypeMeta: metav1.TypeMeta{
					Kind:       schema.KindRuntime,
					APIVersion: schema.APIVersionV2Cluster,
				},
				Metadata: schema.Metadata{
					Repository:      appLoc.Repository,
					Name:            appLoc.Name,
					ResourceVersion: appLoc.Version,
				},
			},
			SystemOptions: &schema.SystemOptions{
				Runtime: &schema.Runtime{
					Locator: loc.Runtime.WithLiteralVersion(appLoc.Version),
				},
				Dependencies: schema.SystemDependencies{
					Runtime: &schema.Dependency{
						Locator: runtimePackageLoc,
					},
				},
			},
		},
	}
}

// SystemApplication creates a new test system application manifest
func SystemApplication(appLoc loc.Locator) *App {
	return &App{
		Manifest: schema.Manifest{
			Header: schema.Header{
				TypeMeta: metav1.TypeMeta{
					Kind:       schema.KindSystemApplication,
					APIVersion: schema.APIVersionV2Cluster,
				},
				Metadata: schema.Metadata{
					Repository:      appLoc.Repository,
					Name:            appLoc.Name,
					ResourceVersion: appLoc.Version,
				},
			},
			Hooks: &schema.Hooks{
				Install: &schema.Hook{
					Type: schema.HookInstall,
					Job: `apiVersion: batch/v1
kind: Job
metadata:
name: app-install
spec:
template:
  spec:
    containers:
      - name: hook
	image: quay.io/gravitational/debian-tall:buster
	command: ["/install"]`,
				},
			},
		},
	}
}

// DefaultClusterApplication creates a new cluster application with defaults
func DefaultClusterApplication(appLoc loc.Locator) *App {
	return ClusterApplication(appLoc, DefaultRuntimeApplication().Build())
}

// ClusterApplication creates a new cluster application with the given locator
// and the runtime application
func ClusterApplication(appLoc loc.Locator, base App) *App {
	return &App{
		Manifest: schema.Manifest{
			Header: schema.Header{
				TypeMeta: metav1.TypeMeta{
					Kind:       schema.KindCluster,
					APIVersion: schema.APIVersionV2Cluster,
				},
				Metadata: schema.Metadata{
					Repository:      appLoc.Repository,
					Name:            appLoc.Name,
					ResourceVersion: appLoc.Version,
				},
			},
			Installer: &schema.Installer{
				Flavors: schema.Flavors{
					Items: []schema.Flavor{
						{
							Name: "one",
							Nodes: []schema.FlavorNode{
								{
									Profile: "node",
									Count:   1,
								},
							},
						},
					},
				},
			},
			NodeProfiles: schema.NodeProfiles{
				{
					Name: "node",
				},
				{
					Name: "kmaster",
					Labels: map[string]string{
						"node-role.kubernetes.io/master": "true",
					},
				},
				{
					Name: "knode",
					Labels: map[string]string{
						"node-role.kubernetes.io/node": "true",
					},
				},
			},
			Hooks: &schema.Hooks{
				NodeAdding: &schema.Hook{
					Type: schema.HookNodeAdding,
					Job: `apiVersion: batch/v1
kind: Job
metadata:
name: pre-join
spec:
template:
  spec:
    containers:
      - name: hook
	image: quay.io/gravitational/debian-tall:buster
	command: ["/bin/echo", "Pre-join hook"]`,
				},
				NodeAdded: &schema.Hook{
					Type: schema.HookNodeAdded,
					Job: `apiVersion: batch/v1
kind: Job
metadata:
name: post-join
spec:
template:
  spec:
    containers:
      - name: hook
	image: quay.io/gravitational/debian-tall:buster
	command: ["/bin/echo", "Post-join hook"]`,
				},
				NetworkInstall: &schema.Hook{
					Type: schema.HookNetworkInstall,
					Job: `apiVersion: batch/v1
kind: Job
metadata:
name: post-join
spec:
template:
  spec:
    containers:
    - name: hook
      image: quay.io/gravitational/debian-tall:buster
      command: ["/bin/echo", "Install overlay network hook"]`,
				},
			},
			SystemOptions: &schema.SystemOptions{
				Runtime: &schema.Runtime{
					Locator: base.Manifest.Locator(),
				},
			},
		},
		Base: &base,
	}
}

// KubernetesResources returns test kubernetes resources
func KubernetesResources() *archive.Item {
	const resourceBytes = `
apiVersion: v1
kind: Pod
metadata:
  name: webserver
  labels:
    app: sample-application
    role: webserver
spec:
  containers:
  - name: webserver
    image: alpine:edge
    ports:
      - containerPort: 80
  nodeSelector:
    role: webserver
---
apiVersion: v1
kind: Pod
metadata:
  name: platform
  labels:
    app: sample-application
    role: server
spec:
  containers:
  - name: platform
    image: busybox:1
    ports:
      - containerPort: 50001
  nodeSelector:
    role: server`
	return archive.ItemFromString("resources/resources.yaml", resourceBytes)
}

// CreateApplicationFromData ??
func CreateApplicationFromData(apps app.Applications, locator loc.Locator, files []*archive.Item, c *check.C) *app.Application {
	data := CreatePackageData(files, c)
	return CreateApplicationFromBinaryData(apps, locator, data, c)
}

// CreateApplicationFromBinaryData ??
func CreateApplicationFromBinaryData(apps app.Applications, locator loc.Locator, data bytes.Buffer, c *check.C) *app.Application {
	var labels map[string]string
	app, err := apps.CreateApp(locator, &data, labels)
	c.Assert(err, check.IsNil)
	c.Assert(app, check.NotNil)

	return app
}

// CreatePackageData generates and returns a new tarball with the specified contents
func CreatePackageData(items []*archive.Item, c *check.C) bytes.Buffer {
	var buf bytes.Buffer
	archive := archive.NewTarAppender(&buf)
	defer archive.Close()

	c.Assert(archive.Add(items...), check.IsNil)

	return buf
}

// ApplicationLayout builds a list of items to place into an application package
func ApplicationLayout(app App, c *check.C) []*archive.Item {
	manifestBytes, err := json.Marshal(app.Manifest)
	c.Assert(err, check.IsNil)
	return append([]*archive.Item{
		archive.DirItem("resources"),
		archive.ItemFromString("resources/app.yaml", string(manifestBytes)),
	}, app.Items...)
}

func collectBaseDependencies(base App, pkgDeps map[loc.Locator]Package, appDeps map[loc.Locator]App, c *check.C) {
	appDeps[base.Manifest.Locator()] = base
	// Add runtime package to dependencies
	runtimePackage, err := base.Manifest.DefaultRuntimePackage()
	c.Assert(err, check.IsNil)
	pkgDeps[*runtimePackage] = Package{
		Loc: *runtimePackage,
	}
	collectDependencies(base, pkgDeps, appDeps)
}

func collectDependencies(app App, pkgDeps map[loc.Locator]Package, appDeps map[loc.Locator]App) {
	collectManifestDependencies(app.Manifest, pkgDeps, appDeps)
	overrideDependencies(app.Dependencies, pkgDeps, appDeps)
	// // apply dependency overrides if configured
	// for _, app := range app.Dependencies.Apps {
	// 	appDeps[app.Manifest.Locator()] = app
	// 	collectDependencies(app, pkgDeps, appDeps)
	// }
}

func collectManifestDependencies(m schema.Manifest, pkgDeps map[loc.Locator]Package, appDeps map[loc.Locator]App) {
	for _, pkg := range m.Dependencies.Packages {
		pkgDeps[pkg.Locator] = Package{
			Loc: pkg.Locator,
		}
	}
	for _, app := range m.Dependencies.Apps {
		appDeps[app.Locator] = SystemApplication(app.Locator).Build()
	}
}

func overrideDependencies(deps Dependencies, pkgDeps map[loc.Locator]Package, appDeps map[loc.Locator]App) {
	for _, pkg := range deps.Packages {
		pkgDeps[pkg.Loc] = pkg
	}
	for _, app := range deps.Apps {
		appDeps[app.Manifest.Locator()] = app
		collectDependencies(app, pkgDeps, appDeps)
	}
}
