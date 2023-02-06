package proto

//go:generate 	protoc --go_out=internal/server/pb --go-grpc_out=internal/server/pb --openapiv2_out . api/EventService.proto
