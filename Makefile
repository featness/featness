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
	@mongod --dbpath /tmp/featness/mongodata --logpath /tmp/featness/mongolog --port 3333 --quiet &

clear_mongo drop drop_db:
	@rm -rf /tmp/featness && mkdir -p /tmp/featness/mongodata

kill_mongo_test:
	@-ps aux | egrep -i 'mongod.+3334' | egrep -v egrep | awk '{ print $$2 }' | xargs kill -9

mongo_test: kill_mongo_test
	@rm -rf /tmp/featness/mongotestdata && mkdir -p /tmp/featness/mongotestdata
	@mongod --dbpath /tmp/featness/mongotestdata --logpath /tmp/featness/mongotestlog --port 3334 --quiet &

run_dashboard run-dashboard dashboard dash: mongo
	@cd featness/dashboard && ./manage.py runserver --settings=featness.dashboard.featness_dashboard.settings_local

update_dump:
	@rm -rf ./mongodump && mkdir -p ./mongodump && mongodump --host localhost --port 3333 --out ./mongodump

restore_dump:
	@mongorestore --host localhost --port 3333 ./mongodata
