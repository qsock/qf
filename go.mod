module github.com/qsock/qf

go 1.13

replace (
	go.etcd.io/etcd/api/v3 => go.etcd.io/etcd/api/v3 v3.0.0-20201103155942-6e800b9b0161
	go.etcd.io/etcd/pkg/v3 => go.etcd.io/etcd/pkg/v3 v3.0.0-20201103155942-6e800b9b0161
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.27.2
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b
	github.com/davecgh/go-spew v1.1.1
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.2
	github.com/kr/beanstalk v0.0.0-20180818045031-cae1762e4858
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.9.1
	github.com/qsock/qvs v0.0.0-20201120090400-82c40240c68e
	github.com/stretchr/testify v1.6.1
	go.etcd.io/etcd/api/v3 v3.5.0-pre
	go.etcd.io/etcd/client/v3 v3.0.0-20201118182908-c11ddc65cea1
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.29.1
)
