
---
authors: Dmitri Shelenin (dshelenin@mulesoft.com)
state: draft
---

# RFD 1 - Kubernetes Runtime Bundles

## What

This document concerns itself with the problem of extracting the Kubernetes runtime
out of the Planet container.

## Why

Historically, Gravity cluster bundles included a specific Kubernetes distribution (henceforth k8s)
bundled inside the container runtime environment called Planet (henceforth Planet).
Thus k8s version upgrades amounted to building a new cluster bundle
(possibly with the same application version) that included the new version of
Planet with the set of new k8s components.
Moreover, this necessarily required that Gravity developers release a new version of Gravity
since a specific Gravity release is explicitly bound to a k8s release (via the Planet dependency).

This is inconvenient as it unnecessarily complicates the vulnerability patching by requiring the
whole release cycle (Gravity-side) and generation of the new cluster bundle (user side).
It requires the upgrade of the whole cluster while, in effect, only the version
of k8s needs to be changed.

## Details

A Gravity cluster bundle contains a snapshot of the world for a specific k8s
application. The bundle includes the snapshot of Planet which, in turn, includes
the snapshot of the specific k8s version.

...

## User Stories

### As a Gravity developer, I want to release a new version of Gravity.

It should be possible to release new versions of Gravity as before. The bundle 
should not have a hard binding to the specific version of Kubernetes - instead it should define
the _baseline_ and the scope of patching or upgrades.

For example, whereas previously, Gravity 7.0.32 depended on k8s 1.17.9, the new release should
define v1.21 as the new baseline, enable seamless upgrades between all v1.21 patches without
releasing a new Gravity version (and consequently a new version of Planet).
It should also define the extended upgrade scope of k8s versions that it can safely support in order
to provide friendly diagnostics when unsupported versions of k8s are either rejected or users
get a warning.

### As a Gravity user, I want to upgrade my application to the new version.

It should be possible for users to build their cluster bundles and upgrade existing installations
as before.

### As a Gravity developer, I want to release a new version of k8s - perhaps as a patch release.

It should be possible to build and publish k8s bundles which could be used for k8s patching and upgrades
on existing installations.

### As a Gravity user, I want to patch my existing installation with a new version of k8s.

It should be possible for users to patch or upgrade k8s versions on an existing installation
using a k8s bundle.

### (Future) As a Gravity user, I want to patch my existing installation with a new version of k8s from Internet

It should be possible for users to patch or upgrade k8s versions on an existing installation with k8s bundles
available from the Internet.


## Kubernetes Bundle

1. As a tarball with a snapshot of all components (e.g. as a Gravity package with a registry dump of all images).
1. As a manifest which describes the component sources that can be pulled from (either directly from the Internet
or from a private registry which is previously populated with the required images)
