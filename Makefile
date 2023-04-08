build:
	go build -o ./bin/linux64

run: build
	./bin/linux64

runleader:	build
	./bin/linux64 --listen 6666

runfollower: build
	./bin/linux64  --listen 6667 --leaderaddr 6666