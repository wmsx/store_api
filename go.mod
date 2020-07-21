module github.com/wmsx/store_api

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.5.0
	github.com/minio/minio-go/v7 v7.0.1
	github.com/wmsx/pkg v0.0.0-20200721154220-d35217b41458
	github.com/wmsx/store_svc v0.0.0-20200721143106-d801144c1fdf
	github.com/wmsx/xconf v0.0.0-20200721142237-75926266fd1a
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
