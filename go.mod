module github.com/ggoop/mdf

require (
	github.com/360EntSecGroup-Skylar/excelize v1.4.1
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20200620013148-b91950f658ec
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/kataras/iris v0.0.0-20190723163631-de9904ff6d38
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/lib/pq v1.7.1
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/robfig/cron v1.2.0
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/spf13/viper v1.7.0
	go.uber.org/dig v1.7.0
	go.uber.org/zap v1.15.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.42.0
	go.uber.org/atomic => github.com/uber-go/atomic v1.4.0
	go.uber.org/multierr => github.com/uber-go/multierr v1.1.0
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190627132806-fd42eb6b336f
	golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9

	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422

	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190711165009-e47acb2ca7f9

	golang.org/x/net => github.com/golang/net v0.0.0-20190628185345-da137c7871d7

	golang.org/x/protobuf => github.com/golang/protobuf v1.3.2

	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

	golang.org/x/sys => github.com/golang/sys v0.0.0-20190712062909-fae7ac547cb7
	golang.org/x/text => github.com/golang/text v0.3.2

	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4

	golang.org/x/tools => github.com/golang/tools v0.0.0-20190717194535-128ec6dfca09

	google.golang.org/api => github.com/googleapis/google-api-go-client v0.7.0

	google.golang.org/appengine => github.com/golang/appengine v1.6.1

	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190716160619-c506a9f90610
	google.golang.org/grpc => github.com/grpc/grpc-go v1.22.0
	gopkg.in/yaml.v2 => github.com/go-yaml/yaml v2.1.0+incompatible
)

go 1.13
