build:
	go build -o userapi main.go

run: build
	./userapi

watch:
	reflex -s -r '\.go$$' make run