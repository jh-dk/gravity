/*
Copyright 2018-2020 Gravitational, Inc.

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

package service

import (
	"bytes"
	"context"
	"path/filepath"
	"time"

	"github.com/gravitational/gravity/lib/app"
	apptest "github.com/gravitational/gravity/lib/app/service/test"
	"github.com/gravitational/gravity/lib/blob/fs"
	"github.com/gravitational/gravity/lib/compare"
	"github.com/gravitational/gravity/lib/defaults"
	"github.com/gravitational/gravity/lib/helm"
	"github.com/gravitational/gravity/lib/loc"
	"github.com/gravitational/gravity/lib/pack"
	"github.com/gravitational/gravity/lib/pack/localpack"
	"github.com/gravitational/gravity/lib/schema"
	"github.com/gravitational/gravity/lib/storage"
	"github.com/gravitational/gravity/lib/storage/keyval"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	. "gopkg.in/check.v1"
)

type PullerSuite struct {
	srcPack pack.PackageService
	dstPack pack.PackageService
	srcApp  app.Applications
	dstApp  app.Applications
}

var _ = Suite(&PullerSuite{})

func (s *PullerSuite) SetUpTest(c *C) {
	_, s.srcPack, s.srcApp = setupServices(c)
	_, s.dstPack, s.dstApp = setupServices(c)
	err := s.srcPack.UpsertRepository("example.com", time.Time{})
	c.Assert(err, IsNil)
	err = s.dstPack.UpsertRepository("example.com", time.Time{})
	c.Assert(err, IsNil)
}

func (s *PullerSuite) TestPullPackage(c *C) {
	loc := loc.MustParseLocator("example.com/package:0.0.1")
	logger := log.WithField("test", "PullPackage")

	_, err := s.srcPack.CreatePackage(loc, bytes.NewBuffer([]byte("data")))
	c.Assert(err, IsNil)

	puller := app.Puller{
		FieldLogger: logger,
		SrcPack:     s.srcPack,
		DstPack:     s.dstPack,
	}
	err = puller.PullPackage(context.TODO(), loc)
	c.Assert(err, IsNil)

	env, err := s.dstPack.ReadPackageEnvelope(loc)
	c.Assert(err, IsNil)
	c.Assert(env.Locator, Equals, loc)

	puller = app.Puller{
		FieldLogger: logger,
		SrcPack:     s.srcPack,
		DstPack:     s.dstPack,
	}
	err = puller.PullPackage(context.TODO(), loc)
	c.Assert(trace.IsAlreadyExists(err), Equals, true)
}

func (s *PullerSuite) TestPullApp(c *C) {
	s.pullApp(c, 0)
}

func (s *PullerSuite) TestPullAppInParallel(c *C) {
	s.pullApp(c, 2)
}

func (s *PullerSuite) pullApp(c *C, parallel int) {
	clusterAppLoc := loc.MustParseLocator("gravitational.io/app:0.0.2")
	existingLoc := apptest.NewDependency("example.com/existing:0.0.1")
	clusterApp := apptest.DefaultClusterApplication(clusterAppLoc)
	clusterApp.Manifest.Dependencies.Packages = append(clusterApp.Manifest.Dependencies.Packages, []schema.Dependency{
		apptest.NewDependency("example.com/new:0.0.1"),
		apptest.NewDependency("example.com/new:0.0.2"),
		existingLoc,
	}...)
	apptest.CreateApplication(apptest.AppRequest{
		App:      clusterApp,
		Packages: s.srcPack,
		Apps:     s.srcApp,
	}, c)
	// `existing` package is also available in the destination service
	apptest.CreatePackage(apptest.PackageRequest{
		Package:  apptest.Package{Loc: existingLoc.Locator},
		Packages: s.dstPack,
	}, c)

	puller := app.Puller{
		SrcPack:  s.srcPack,
		DstPack:  s.dstPack,
		SrcApp:   s.srcApp,
		DstApp:   s.dstApp,
		Upsert:   true,
		Parallel: parallel,
	}
	err := puller.PullApp(context.TODO(), clusterAppLoc)
	c.Assert(err, IsNil)

	verifyPackages(s.dstPack, []loc.Locator{
		loc.MustParseLocator("gravitational.io/app:0.0.2"),
		apptest.RuntimeApplicationLoc,
		apptest.RuntimePackageLoc,
		loc.MustParseLocator("example.com/existing:0.0.1"),
		loc.MustParseLocator("example.com/new:0.0.1"),
		loc.MustParseLocator("example.com/new:0.0.2"),
	}, c)

	local, err := s.dstApp.GetApp(clusterAppLoc)
	c.Assert(err, IsNil)
	c.Assert(local.Package, Equals, clusterAppLoc)

	puller = app.Puller{
		SrcPack:  s.srcPack,
		DstPack:  s.dstPack,
		SrcApp:   s.srcApp,
		DstApp:   s.dstApp,
		Parallel: parallel,
	}
	err = puller.PullApp(context.TODO(), clusterAppLoc)
	c.Assert(trace.IsAlreadyExists(err), Equals, true)
}

func setupServices(c *C) (storage.Backend, pack.PackageService, *Applications) {
	dir := c.MkDir()

	backend, err := keyval.NewBolt(keyval.BoltConfig{
		Path: filepath.Join(dir, "bolt.db"),
	})
	c.Assert(err, IsNil)

	objects, err := fs.New(dir)
	c.Assert(err, IsNil)

	packService, err := localpack.New(localpack.Config{
		Backend:     backend,
		UnpackedDir: filepath.Join(dir, defaults.UnpackedDir),
		Objects:     objects,
	})
	c.Assert(err, IsNil)

	charts, err := helm.NewRepository(helm.Config{
		Packages: packService,
		Backend:  backend,
	})
	c.Assert(err, IsNil)

	appService, err := New(Config{
		Backend:  backend,
		StateDir: filepath.Join(dir, defaults.ImportDir),
		Packages: packService,
		Charts:   charts,
	})
	c.Assert(err, IsNil)

	return backend, packService, appService
}

func verifyPackages(packages pack.PackageService, expected []loc.Locator, c *C) {
	repositories, err := packages.GetRepositories()
	c.Assert(err, IsNil)

	var result []loc.Locator
	for _, repository := range repositories {
		packages, err := packages.GetPackages(repository)
		c.Assert(err, IsNil)
		result = append(result, locators(packages)...)
	}

	c.Assert(packagesByName(result), compare.SortedSliceEquals, packagesByName(expected))
}

func locators(envelopes []pack.PackageEnvelope) []loc.Locator {
	out := make([]loc.Locator, 0, len(envelopes))
	for _, env := range envelopes {
		out = append(out, env.Locator)
	}
	return out
}

type packagesByName []loc.Locator

func (r packagesByName) Len() int           { return len(r) }
func (r packagesByName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r packagesByName) Less(i, j int) bool { return r[i].String() < r[j].String() }
