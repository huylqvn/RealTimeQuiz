.PHONY: build dev
pull:
	git pull origin main		
stop:
	docker rm $$(docker stop $$(docker ps -a -q -f "name=quizserver")) || true
build:
	docker build -t quizserver .
run:
	docker run --log-driver=json-file --log-opt max-size=10m --log-opt max-file=5 -d --restart unless-stopped -p 7000:7000 --name quizserver quizserver
log:
	docker logs -f $$(docker ps -a -q -f "name=quizserver")
auto: pull build stop run