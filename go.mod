module github.com/gravitational/gravity

go 1.14

require (
	cloud.google.com/go v0.34.0
	github.com/DATA-DOG/go-sqlmock v1.5.0 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38 // indirect
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.0.0-20210611225437-1a2716eca9d6 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/apparentlymart/go-cidr v1.0.0 // indirect
	github.com/aws/aws-sdk-go v1.25.41
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/boltdb/bolt v1.3.1
	github.com/boombuler/barcode v0.0.0-20161226211916-fe0f26ff6d26 // indirect
	github.com/buger/goterm v0.0.0-20140416104154-af3f07dadc88
	github.com/cenkalti/backoff v1.1.0
	github.com/chai2010/gettext-go v0.0.0-20170215093142-bf70f2a70fb1 // indirect
	github.com/cloudflare/cfssl v0.0.0-20180726162950-56268a613adf
	github.com/cloudfoundry/gosigar v1.1.1-0.20180406153506-1375283248c3
	github.com/codahale/hdrhistogram v0.9.1-0.20161010025455-3a0bb77429bd
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-oidc v0.0.1 // indirect
	github.com/coreos/go-semver v0.2.0
	github.com/coreos/prometheus-operator v0.0.0-00010101000000-000000000000 // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20191101170500-ac7306503d23
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20150526203908-9cbd2a1374f4
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/emicklei/go-restful v2.11.0+incompatible // indirect
	github.com/evanphx/json-patch v3.0.0+incompatible
	github.com/fatih/color v1.9.0
	github.com/fsouza/go-dockerclient v1.6.5
	github.com/garyburd/redigo v0.0.0-20151029235527-6ece6e0a09f2 // indirect
	github.com/ghodss/yaml v1.0.1-0.20180820084758-c7ce16629ff4
	github.com/gizak/termui v2.3.0+incompatible
	github.com/go-openapi/analysis v0.19.4 // indirect
	github.com/go-openapi/runtime v0.19.3
	github.com/go-openapi/strfmt v0.19.2 // indirect
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/gokyle/hotp v0.0.0-20160218004637-c180d57d286b
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff // indirect
	github.com/golang/protobuf v1.3.5
	github.com/google/btree v1.0.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gorilla/handlers v0.0.0-20151124211609-e96366d97736 // indirect
	github.com/gravitational/configure v0.0.0-20191213111049-fce91dea0d0d
	github.com/gravitational/coordinate v0.0.0-20180225144834-2bc9a83f6fe2
	github.com/gravitational/form v0.0.0-20151109031454-c4048f792f70
	github.com/gravitational/go-vhost v0.0.0-20171024163855-94d0c42e3263
	github.com/gravitational/kingpin v2.1.10+incompatible // indirect
	github.com/gravitational/license v0.0.0-20171013193735-f3111b1818ce
	github.com/gravitational/magnet v0.2.7-0.20210609203954-0a59f53f530e
	github.com/gravitational/oxy v0.0.0-20180629203109-e4a7e35311e6 // indirect
	github.com/gravitational/rigging v0.0.0-20200803191640-2a0fba75cac5
	github.com/gravitational/roundtrip v1.0.0
	github.com/gravitational/satellite v0.0.0-00010101000000-000000000000
	github.com/gravitational/tail v1.0.1
	github.com/gravitational/teleport v3.2.17+incompatible
	github.com/gravitational/trace v1.1.14
	github.com/gravitational/ttlmap v0.0.0-20171116003245-91fd36b9004c
	github.com/gravitational/version v0.0.2-0.20170324200323-95d33ece5ce1
	github.com/gravitational/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/hashicorp/go-getter v0.0.0-20180809191950-4bda8fa99001 // indirect
	github.com/hashicorp/go-hclog v0.0.0-20180828044259-75ecd6e6d645 // indirect
	github.com/hashicorp/go-plugin v0.0.0-20180814222501-a4620f9913d1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20180822193130-ed8144cda141 // indirect
	github.com/hashicorp/hil v0.0.0-20170627220502-fa9f258a9250 // indirect
	github.com/hashicorp/terraform v0.11.7
	github.com/hashicorp/yamux v0.0.0-20180826203732-cc6d2ea263b2 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/jonboulle/clockwork v0.2.0
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/julienschmidt/httprouter v1.2.0
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/kylelemons/godebug v0.0.0-20170820004349-d65d576e9348
	github.com/lib/pq v1.2.0 // indirect
	github.com/magefile/mage v1.9.0
	github.com/mailgun/lemma v0.0.0-20160211003854-e8b0cd607f58
	github.com/mailgun/metrics v0.0.0-20150124003306-2b3c4565aafd // indirect
	github.com/mailgun/minheap v0.0.0-20131208021033-7c28d80e2ada // indirect
	github.com/mailgun/timetools v0.0.0-20150505213551-fd192d755b00
	github.com/mailgun/ttlmap v0.0.0-20150816203249-16b258d86efc // indirect
	github.com/maruel/panicparse v1.1.2-0.20180806203336-f20d4c4d746f // indirect
	github.com/mdp/rsc v0.0.0-20160131164516-90f07065088d // indirect
	github.com/miekg/dns v1.1.41
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-ps v1.0.0
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/mreiferson/go-httpclient v0.0.0-20160630210159-31f0106b4474 // indirect
	github.com/nsf/termbox-go v0.0.0-20190325093121-288510b9734e // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/selinux v1.3.0
	github.com/pborman/uuid v1.2.0
	github.com/pquerna/otp v0.0.0-20160912161815-54653902c20e // indirect
	github.com/prometheus/alertmanager v0.17.0
	github.com/prometheus/client_golang v1.4.0
	github.com/prometheus/common v0.9.1
	github.com/rubenv/sql-migrate v0.0.0-20190902133344-8926f37f0bc1 // indirect
	github.com/russellhaering/gosaml2 v0.0.0-20170515204909-8908227c114a // indirect
	github.com/russellhaering/goxmldsig v1.1.1-0.20200930045116-0bf1c1013037 // indirect
	github.com/santhosh-tekuri/jsonschema v1.2.2
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tstranex/u2f v0.0.0-20160508205855-eb799ce68da4
	github.com/ulikunitz/xz v0.5.4 // indirect
	github.com/vulcand/oxy v0.0.0-20160623194703-40720199a16c
	github.com/vulcand/predicate v1.1.0
	github.com/xtgo/set v1.0.0
	github.com/zclconf/go-cty v0.0.0-20180829180805-c2393a5d54f2 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.mongodb.org/mongo-driver v1.0.4 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1
	golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44
	golang.org/x/tools v0.0.0-20191212051200-825cb0626375 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20160818020120-3f83fa500528 // indirect
	gopkg.in/square/go-jose.v2 v2.2.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.0.0
	k8s.io/apimachinery v0.17.5-beta.0
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/helm v2.15.2+incompatible
	k8s.io/kube-aggregator v0.0.0
	k8s.io/kubernetes v1.15.5
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
)

replace (
	github.com/boltdb/bolt => github.com/gravitational/bolt v1.3.2-gravitational
	github.com/cloudflare/cfssl => github.com/gravitational/cfssl v0.0.0-20180619163912-4b8305b36ad0
	github.com/coreos/go-oidc => github.com/gravitational/go-oidc v0.0.1
	github.com/coreos/prometheus-operator => github.com/gravitational/prometheus-operator v0.35.2
	github.com/docker/docker => github.com/gravitational/moby v1.4.2-0.20191008111026-2adf434ca696
	github.com/fvbommel/sortorder => github.com/fvbommel/sortorder v1.0.1
	github.com/google/certificate-transparency-go => github.com/gravitational/certificate-transparency-go v0.0.0-20180803094710-99d8352410cb
	github.com/gravitational/satellite => github.com/a-palchikov/satellite v0.0.9-0.20210701113341-c00eafc55855
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
	github.com/julienschmidt/httprouter => github.com/julienschmidt/httprouter v1.1.0
	github.com/magefile/mage => github.com/knisbet/mage v1.9.1-0.20210609142646-749a704341ac
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	github.com/sirupsen/logrus => github.com/gravitational/logrus v1.4.3
	gopkg.in/alecthomas/kingpin.v2 => github.com/gravitational/kingpin v2.1.11-0.20180808090833-85085db9f49b+incompatible
	k8s.io/api => k8s.io/api v0.15.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.15.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.15.7
	k8s.io/apiserver => k8s.io/apiserver v0.15.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.15.7
	k8s.io/client-go => k8s.io/client-go v0.15.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.15.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.15.7
	k8s.io/code-generator => k8s.io/code-generator v0.15.7
	k8s.io/component-base => k8s.io/component-base v0.15.7
	k8s.io/cri-api => k8s.io/cri-api v0.15.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.15.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.15.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.15.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.15.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.15.7
	k8s.io/kubelet => k8s.io/kubelet v0.15.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.15.7
	k8s.io/metrics => k8s.io/metrics v0.15.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.15.7
)
