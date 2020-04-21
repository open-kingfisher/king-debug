module kingfisher/king-debug

go 1.14

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-20171023200535-7848b8beb9d3
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	k8s.io/apimachinery v0.0.0-20190820100750-21ddcbbef9e1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubernetes v1.13.1
	kingfisher/kf v0.0.0-00010101000000-000000000000
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
	kingfisher/kf => ../kf
)
