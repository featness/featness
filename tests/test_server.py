#!/usr/bin/python
# -*- coding: utf-8 -*-


from preggy import expect
from mock import patch

from featness import __version__
import featness.api.server
from tests.base import ApiTestCase


class ApiServerTestCase(ApiTestCase):
    def test_healthcheck(self):
        response = self.fetch('/healthcheck')
        expect(response.code).to_equal(200)
        expect(response.body).to_be_like('WORKING')

    def test_get_version(self):
        response = self.fetch('/version')
        expect(response.code).to_equal(200)
        expect(response.body).to_be_like(__version__)

    def test_server_handlers(self):
        srv = featness.api.server.FeatnessApiServer()
        handlers = srv.get_handlers()

        expect(handlers).not_to_be_null()
        expect(handlers).to_length(1)

    def test_server_plugins(self):
        srv = featness.api.server.FeatnessApiServer()
        plugins = srv.get_plugins()

        expect(plugins).to_length(0)

    @patch('featness.api.server.FeatnessApiServer')
    def test_server_main_function(self, server_mock):
        featness.api.server.main()
        expect(server_mock.run.called).to_be_true()

    def test_has_connected_to_redis(self):
        app = self.get_app()
        expect(app.redis).not_to_be_null()


