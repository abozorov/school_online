.PHONY: services

run_api:
	go run cmd/api_gateway/main.go

run_user:
	go run cmd/user/main.go

run_classroom:
	go run cmd/classroom/main.go


.PHONY:
build_user_proto:
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/userpb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/userpb --go-grpc_opt=paths=source_relative \
	user/v1/user.proto

build_raiting_proto:
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/raitingpb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/raitingpb --go-grpc_opt=paths=source_relative \
	raiting/v1/raiting.proto

build_classroom_proto:
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/classroompb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/classroompb --go-grpc_opt=paths=source_relative \
	classroom/v1/classroom.proto


.PHONY:
migrate_classroom:
	migrate -path cmd/classroom/migrations -database "postgres://postgres:postgres@localhost:5433/school_onlinedb?sslmode=disable&x-migrations-table=schema_migrations_classroom" up


migrate_user:
	migrate -path cmd/user/migrations -database "postgres://postgres:postgres@localhost:5433/school_onlinedb?sslmode=disable&x-migrations-table=schema_migrations_user" up


migrate_raiting:
	migrate -path cmd/raiting/migrations -database "postgres://postgres:postgres@localhost:5433/school_onlinedb?sslmode=disable&x-migrations-table=schema_migrations_raiting" up