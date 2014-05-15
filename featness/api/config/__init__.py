#!/usr/bin/python
# -*- coding: utf-8 -*-

from derpconf.config import Config  # NOQA


Config.define('MONGODB_HOST', '127.0.0.1', 'MongoDB Host', 'MongoDB')
Config.define('MONGODB_PORT', 3333, 'MongoDB port', 'MongoDB')
Config.define('MONGODB_DATABASE', 'featness', 'MongoDB database name', 'MongoDB')
Config.define('MONGODB_USER', None, 'MongoDB Authenticating User', 'MongoDB')
Config.define('MONGODB_PASS', None, 'MongoDB Authenticating Password', 'MongoDB')

Config.define('REDIS_HOST', '127.0.0.1', 'Redis Host', 'Redis')
Config.define('REDIS_PORT', 4444, 'Redis port', 'Redis')
Config.define('REDIS_DB', 0, 'Redis database index', 'Redis')
Config.define('REDIS_PASS', None, 'Redis Authenticating Password', 'Redis')
