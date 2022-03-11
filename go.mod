module github.com/prometheus-community/pushprox

go 1.13

require (
	github.com/Showmax/go-fqdn v1.0.0
	github.com/cenkalti/backoff/v4 v4.1.2
	github.com/go-kit/log v0.2.0
	github.com/google/uuid v1.3.0
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/prometheus/common v0.32.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	k8s.io/sample-controller v0.23.4
)
