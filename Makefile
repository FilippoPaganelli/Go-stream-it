build:
	GOOS=linux GOARCH=arm64 go build -o go-stream-it
	chmod +x go-stream-it

deploy:
	sudo systemctl restart go-stream-it.service

run:
	./go-stream-it

logs:
	journalctl -u go-stream-it.service -f
