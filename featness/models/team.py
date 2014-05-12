from mongoengine import *
from mongoengine.django.auth import User


class Team(Document):
    name = StringField(required=True)
    slug = StringField(required=True)
    members = ListField(ReferenceField(User))
    owner = ReferenceField(User)
