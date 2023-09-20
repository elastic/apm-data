module github.com/elastic/apm-data/tools

go 1.20

replace github.com/planetscale/vtprotobuf => github.com/elastic/vtprotobuf v0.0.0-20230920122325-4b486ed04f6e

require github.com/planetscale/vtprotobuf v0.5.0

require google.golang.org/protobuf v1.26.0 // indirect
