build-linux:
	GOOS=linux CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -ldflags "-linkmode external -extldflags -static" -o smolgu-lib lib.go

build:
	go build -o smolgu-lib lib.go

push:
	rsync -avzph --progress ./ webmaster@51.158.171.85:/home/webmaster/apps/smolgu-lib/
start:
	PORT=3001 ./smolgu-lib web -c=./conf/app.ini

run: build start

clean:
	rm -rf log/users.log data/cache.db data/sessions data/bold.db data/ids data/vk.log data/db/pages.index data/db/users.kv data/db/db.sqlite data/db/kv.bolt data/db/users.sqlite data/bolt.db
