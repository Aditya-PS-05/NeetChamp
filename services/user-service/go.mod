module github.com/Aditya-PS-05/NeetChamp/user-service

go 1.22.2

require (
	github.com/Aditya-PS-05/NeetChamp/auth-service v0.0.0-00010101000000-000000000000
	github.com/Aditya-PS-05/NeetChamp/shared-libs v0.0.0-00010101000000-000000000000
	gorm.io/gorm v1.25.12
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5 // indirect
	gorm.io/driver/postgres v1.5.11
)

replace (
	github.com/Aditya-PS-05/NeetChamp/auth-service => ../../services/auth-service
	github.com/Aditya-PS-05/NeetChamp/shared-libs => ../../shared-libs
)
