# santaclaus
Indexer

## Docker 
### Build

```shell
docker compose --profile launch --env-file env/santaclaus.env build
```

### Run

```shell
docker compose --profile launch --env-file env/santaclaus.env up
```

## Manually
### Build

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go mod download
go mod verify
make gRPC
make
```

### Run

```shell
# the database has to be executed separately
./santaclaus
```

## Test
### Build and run nit tests

```shell
## launches the database to be able to connect with mongo compass (even after the tests have run)
docker compose --env-file env/unit_tests.env up --build 
./scripts/unit_tests.sh
```
Alternatively, `unit_tests.sh` script can launch docker :
```shell
./scripts/unit_tests.sh --docker
```

<!-- ## Learn -->

<!-- ### Documentation -->

<!-- ### Contribute -->
