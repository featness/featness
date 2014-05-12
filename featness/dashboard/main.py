#import datetime
import os

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

# Flask views
@app.route('/')
def index():
    return '<a href="/admin/">Click me to get to Admin!</a>'

class PasswordEnabledView(ModelView):
    def on_model_change(self, form, model, is_created):
        model.password = app.bcrypt.generate_password_hash(form.password.data)
        super(PasswordEnabledView, self).on_model_change(form, model, is_created)


def start_server():
    host = os.environ.get('HOST', '0.0.0.0')
    port = int(os.environ.get('PORT', 8000))
    debug = os.environ.get('DEBUG', None) is not None

    app.config['SECRET_KEY'] = '9124518258'
    app.config['MONGODB_SETTINGS'] = {'DB': 'featness', 'PORT': 3333}

    db = MongoEngine()
    db.init_app(app)

    login_manager.init_app(app)

    app.bcrypt = Bcrypt(app)

    adm = admin.Admin(app, 'Featness')

    adm.add_view(PasswordEnabledView(User))
    adm.add_view(ModelView(Team))
    adm.add_view(ModelView(Project))
    adm.add_view(ModelView(Feature))

    app.run(debug=debug, host=host, port=port)


if __name__ == '__main__':
    start_server()
