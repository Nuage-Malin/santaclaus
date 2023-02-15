# santaclaus
Indexer

## Docker 
### Build

```shell
docker compose --env-file santaclaus.env build
```

### Run

```shell
docker compose --env-file santaclaus.env up
```

## Manually
### Build

```shell
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
docker compose --env-file unit_tests.env up --build
```

<!-- ## Learn -->

<!-- ### Documentation -->

<!-- ### Contribute -->