#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os

from cow.testing import CowTestCase
from tornado.httpclient import AsyncHTTPClient

from featness.api.config import Config
from featness.api.server import FeatnessApiServer


class ApiTestCase(CowTestCase):
    def setUp(self):
        super(ApiTestCase, self).setUp()

    def tearDown(self):
        super(ApiTestCase, self).tearDown()

    def get_config(self):
        return dict(
            MONGODB_PORT=3334,
            REDIS_PORT=4445
        )

    def get_server(self):
        cfg = Config(**self.get_config())
        debug = os.environ.get('DEBUG_TESTS', 'False').lower() == 'true'

        self.server = FeatnessApiServer(config=cfg, debug=debug)

        return self.server

    def get_app(self):
        app = super(ApiTestCase, self).get_app()
        app.http_client = AsyncHTTPClient(self.io_loop)
        return app
