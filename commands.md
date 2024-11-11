# akits-e9

```bash
cd 3-exercise-ring
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative grpc/proto.proto
```

```bash
cd 2-exercise
alacritty -e zsh -c "go run . 4 5050; exec zsh" &
alacritty -e zsh -c "go run . 4 5051 5050; exec zsh" &
alacritty -e zsh -c "go run . 4 5052 5051; exec zsh" &
alacritty -e zsh -c "go run . 4 5053 5052; exec zsh" &

```

```bash
cd 2-exercise
if [[ $(sort log.txt | uniq -d) ]]; then
    echo "Duplicates found:"
    sort log.txt | uniq -d
else
    echo "No duplicates found."
fi
```