module slem

go 1.17

replace slem/cmd => ./cmd

require (
	github.com/aws/aws-sdk-go v1.44.107
	github.com/mattn/go-sqlite3 v1.14.15
	github.com/spf13/cobra v1.5.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
