#!/usr/bin/python
# -*- coding: utf-8 -*-

from ujson import dumps
from tornado.web import RequestHandler

from featness import __version__


class BaseHandler(RequestHandler):
    def initialize(self, *args, **kw):
        super(BaseHandler, self).initialize(*args, **kw)
        self._session = None

    def log_exception(self, typ, value, tb):
        for handler in self.application.error_handlers:
            handler.handle_exception(
                typ, value, tb, extra={
                    'url': self.request.full_url(),
                    'ip': self.request.remote_ip,
                    'holmes-version': __version__
                }
            )

        super(BaseHandler, self).log_exception(typ, value, tb)

    def write_json(self, obj):
        self.set_header("Content-Type", "application/json")
        self.write(dumps(obj))
