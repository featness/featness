#!/usr/bin/python
# -*- coding: utf-8 -*-

from cow.server import Server
from tornado.httpclient import AsyncHTTPClient
from mongoengine import connect

from featness import __version__
from featness.api.handlers import BaseHandler
from featness.api.config import Config


def main():
    AsyncHTTPClient.configure("tornado.curl_httpclient.CurlAsyncHTTPClient")
    FeatnessApiServer.run()


class VersionHandler(BaseHandler):
    def get(self):
        self.write(__version__)


class FeatnessApiServer(Server):
    def __init__(self, db=None, debug=None, *args, **kw):
        super(FeatnessApiServer, self).__init__(*args, **kw)

        self.force_debug = debug
        self.db = db

    def initialize_app(self, *args, **kw):
        super(FeatnessApiServer, self).initialize_app(*args, **kw)

        if self.force_debug is not None:
            self.debug = self.force_debug

    def get_config(self):
        return Config

    def get_handlers(self):
        handlers = [
            ('/version/?', VersionHandler),
        ]

        return tuple(handlers)

    def after_start(self, io_loop):
        connect(
            self.config.MONGODB_DATABASE,
            host=self.config.MONGODB_HOST,
            port=self.config.MONGODB_PORT,
            username=self.config.MONGODB_USER,
            password=self.config.MONGODB_PASS
        )


if __name__ == '__main__':
    main()
