container:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o dumpos
	docker build -t xtracdev/dumpos:latest .

push:
	docker push xtracdev/dumpos latest

clean:
	rm -f dumpos
