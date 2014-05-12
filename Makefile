setup: setup-python

setup-python:
	@pip install -U -e .\[tests\]

test: mongo_test _go_test

kill_redis:
	@-redis-cli -p 4444 shutdown

redis: kill_redis
	@redis-server ./redis.conf; sleep 1
	@redis-cli -p 4444 info > /dev/null

kill_mongo:
	@-ps aux | egrep -i 'mongod.+3333' | egrep -v egrep | awk '{ print $$2 }' | xargs kill -9

mongo: kill_mongo
	@rm -rf /tmp/featness/mongodata && mkdir -p /tmp/featness/mongodata
	@mongod --dbpath /tmp/featness/mongodata --logpath /tmp/featness/mongolog --port 3333 --quiet &

kill_mongo_test:
	@-ps aux | egrep -i 'mongod.+3334' | egrep -v egrep | awk '{ print $$2 }' | xargs kill -9

mongo_test: kill_mongo_test
	@rm -rf /tmp/featness/mongotestdata && mkdir -p /tmp/featness/mongotestdata
	@mongod --dbpath /tmp/featness/mongotestdata --logpath /tmp/featness/mongotestlog --port 3334 --quiet &
