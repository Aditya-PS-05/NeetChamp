module github.com/Aditya-PS-05/NeetChamp/grpc-gateway

go 1.23.0

toolchain go1.23.7

replace github.com/Aditya-PS-05/NeetChamp/shared-libs => ../../shared-libs

require (
	github.com/Aditya-PS-05/NeetChamp/shared-libs v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3
	google.golang.org/grpc v1.71.0
)

require (
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250303144028-a0af3efb3deb // indirect
	google.golang.org/protobuf v1.36.5 // indirect
)
