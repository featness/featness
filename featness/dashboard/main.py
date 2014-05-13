#import datetime
import os
from uuid import uuid4
import hashlib

from flask import Flask
from flask.ext import admin
from flask.ext.mongoengine import MongoEngine
from flask.ext.admin.contrib.mongoengine import ModelView
import flask.ext.login as login
from flask.ext.bcrypt import Bcrypt

from featness.models import User, Team, Project, Feature


app = Flask("featness")


login_manager = login.LoginManager()


@login_manager.user_loader
def load_user(userid):
    return User.objects.get(userid)


@app.route('/')
def index():
    return '<a href="/admin/">Click me to get to Admin!</a>'


class PasswordEnabledView(ModelView):
    def on_model_change(self, form, model, is_created):
        model.password = app.bcrypt.generate_password_hash(form.password.data)
        super(PasswordEnabledView, self).on_model_change(form, model, is_created)


class FeatureCreateView(ModelView):
    form_columns = ('name', 'slug', 'project')

    def get_new_key(self):
        found = True

        new_key = None
        while found:
            m = hashlib.md5()
            m.update(str(uuid4()).encode('utf-8'))
            new_key = m.hexdigest()[:8]
            found = Feature.objects(feature_key=new_key)

        return new_key

    def _on_model_change(self, form, model, is_created):
        if is_created:
            model.feature_key = self.get_new_key()

        super(FeatureCreateView, self)._on_model_change(form, model, is_created)

def start_server():
    host = os.environ.get('HOST', '0.0.0.0')
    port = int(os.environ.get('PORT', 8000))
    debug = os.environ.get('DEBUG', None) is not None

    app.config['SECRET_KEY'] = '9124518258'
    app.config['MONGODB_SETTINGS'] = {'DB': 'featness', 'PORT': 3333}

    app.db = MongoEngine()
    app.db.init_app(app)

    login_manager.init_app(app)

    app.bcrypt = Bcrypt(app)

    adm = admin.Admin(app, 'Featness')

    adm.add_view(PasswordEnabledView(User))
    adm.add_view(ModelView(Team))
    adm.add_view(ModelView(Project))
    adm.add_view(FeatureCreateView(Feature))

    app.run(debug=debug, host=host, port=port)


if __name__ == '__main__':
    start_server()
