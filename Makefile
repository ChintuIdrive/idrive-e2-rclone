IMAGE_NAME = rclone-test

build:
	GOOS=linux go build -o app
	docker build -t ${IMAGE_NAME} .

run:
	docker run ${IMAGE_NAME}

clean:
	docker rmi $(IMAGE_NAME)
	docker system prune -f
