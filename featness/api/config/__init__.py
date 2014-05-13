#!/usr/bin/python
# -*- coding: utf-8 -*-

from derpconf.config import Config  # NOQA


Config.define('MONGODB_HOST', '127.0.0.1', 'MongoDB Host', 'MongoDB')
Config.define('MONGODB_PORT', 3333, 'MongoDB port', 'MongoDB')
Config.define('MONGODB_DATABASE', 'featness', 'MongoDB database name', 'MongoDB')
Config.define('MONGODB_USER', None, 'MongoDB Authenticating User', 'MongoDB')
Config.define('MONGODB_PASS', None, 'MongoDB Authenticating Password', 'MongoDB')
