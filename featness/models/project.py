from mongoengine import Document, StringField, ReferenceField
from featness.models.team import Team


class Project(Document):
    name = StringField(required=True)
    slug = StringField(required=True)
    team = ReferenceField(Team)

    def __str__(self):
        return self.name
