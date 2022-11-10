up:
	docker-compose up -d --remove-orphans

restart:
	docker-compose restart

build-app:
	docker-compose build acl

# get service IP-address
# ip address | grep inet

wrk_identity:
	wrk \
		-t10 \
		-c2000 \
		-d30s \
		-H 'Host: sweet-home.com' \
		--timeout 2s \
		-s ./load_test/wrk_identity.lua \
		http://localhost:8080/identity

wrk_alert:
	wrk \
		-t10 \
		-c2000 \
		-d30s \
		-H 'Host: sweet-home.com' \
		--timeout 2s \
		-s ./load_test/wrk_alert.lua \
		http://localhost:8080/alert

# add data to the database (POST)
# curl -X POST -F 'token=admin_secret_token' http://127.0.0.1:8080/identity
# curl http://127.0.0.1:8080/alert

# performance test
# wrk -t1 -c1 -d5m -s ./load_test/wrk_identity.lua http://localhost:8080/identity
# wrk -t1 -c1 -d5m -s ./load_test/wrk_alert.lua http://localhost:8080/alert