## Build

```
docker build -t samprom:v1 .
```

## Run

```
docker run --name prometheus -d -p 9090:9090 samprom:v1
```
