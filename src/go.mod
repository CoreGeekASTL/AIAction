module GIDS

go 1.25.0

require (
	CSPGSOMF v0.0.0-20231025071637-2532514dda8f
	CSPNTP_SDK_GO v0.0.0-20220818144051-1e83a4f813a1
	Go-chassis-extend v0.0.0-20231017151318-95bf86b05f6d
	code.huawei.com/fusionstage/auditlog v1.9.7
	code.huawei.com/fusionstage/greatwall-sdk-go v1.9.6
	gitee.com/opengauss/openGauss-connector-go-pq v1.0.7
	github.com/beego/beego/v2 v2.1.0
	github.com/google/uuid v1.6.0
	github.com/minio/minio-go/v7 v7.0.69
	github.com/redis/go-redis/v9 v9.0.5
	github.com/smartystreets/goconvey v1.8.1
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.6 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.15.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.74.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.54.0 // indirect
)

replace (
	AlarmSDK_GO => ./stubs/AlarmSDK_GO
	CSPGSOMF => ./stubs/CSPGSOMF
	CSPGoMonitorSDK => ./stubs/CSPGoMonitorSDK
	CSPNTP_SDK_GO => ./stubs/CSPNTP_SDK_GO
	Go-chassis-extend => ./stubs/Go-chassis-extend
	code.huawei.com/do-libraries/Oscar-Go => ./stubs/code.huawei.com/do-libraries/Oscar-Go
	code.huawei.com/do-libraries/goutils => ./stubs/code.huawei.com/do-libraries/goutils
	code.huawei.com/fusionstage/auditlog => ./stubs/code.huawei.com/fusionstage/auditlog
	code.huawei.com/fusionstage/greatwall-sdk-go => ./stubs/code.huawei.com/fusionstage/greatwall-sdk-go
	code.huawei.com/paaslite/dev-tool/gomockit => ./stubs/code.huawei.com/paaslite/dev-tool/gomockit
	github.com/cenkalti/backoff => github.com/cenkalti/backoff/v4 v4.1.3
	github.com/minio/minio-go/v7 => github.com/minio/minio-go/v7 v7.0.69
	github.com/stretchr/testify => github.com/stretchr/testify v1.8.2
	golang.org/x/net => golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/text => golang.org/x/text v0.3.8
	open.codehub.huawei.com/innersource/corebuf => ./stubs/open.codehub.huawei.com/innersource/corebuf
)

require (
	AlarmSDK_GO v0.0.0-00010101000000-000000000000
	CSPGoMonitorSDK v0.0.0-00010101000000-000000000000
	code.huawei.com/paaslite/dev-tool/gomockit v1.1.0
	gopkg.in/yaml.v2 v2.4.0
)
