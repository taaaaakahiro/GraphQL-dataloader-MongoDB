dck=docker-compose

conf:
	${dck} config
up:
	${dck} up --build -d
down:
	${dck} down -v