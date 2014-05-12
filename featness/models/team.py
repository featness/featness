from mongoengine import Document, StringField, ListField, ReferenceField
from featness.models.user import User


class Team(Document):
    name = StringField(required=True)
    slug = StringField(required=True)
    members = ListField(ReferenceField(User))
    owner = ReferenceField(User)
    
    def __str__(self):
        return "Team %s" % self.name
