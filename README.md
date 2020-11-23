# HTTP File Server using GO

## Download

[Release](https://github.com/chentanyi/fileserver/releases)

## Docker

[Docker Hub](https://hub.docker.com/r/chentanyi/fileserver)

## Build

#### with GO 1.12+

```bash
git clone https://github.com/chentanyi/fileserver
bash prebuild.sh
go build
```

#### with Docker only

```bash
git clone https://github.com/chentanyi/fileserver
docker build -f Dockerfile.prebuild .
```

#### with Docker and GO 1.12+

```bash
git clone https://github.com/chentanyi/fileserver
bash prebuild.sh
docker build .
```