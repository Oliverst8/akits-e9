# akits-e9

```bash
cd 2-exercise
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/proto.proto
```

```bash
cd 2-exercise
kitty -- fish -c "go run . 4 5050; exec fish" &
kitty -- fish -c "go run . 4 5051 5050; exec fish" &
kitty -- fish -c "go run . 4 5052 5051; exec fish" &
kitty -- fish -c "go run . 4 5053 5052; exec fish" &

```