module github.com/wmsx/store_api

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/gorilla/sessions v1.2.0
	github.com/micro/go-micro/v2 v2.9.1
	github.com/minio/minio-go/v6 v6.0.57
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/wmsx/pkg v0.0.0-20200710124640-b827730961c0
	github.com/wmsx/store_svc v0.0.0-00010101000000-000000000000
	github.com/wmsx/xconf v0.0.0-20200710193800-f97c7e3c9e84
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	github.com/wmsx/pkg => /Users/zengqiang96/codespace/sx/pkg
	github.com/wmsx/store_svc => /Users/zengqiang96/codespace/sx/store_svc
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
