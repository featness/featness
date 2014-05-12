from mongoengine import Document, StringField, ReferenceField
from featness.models.project import Project


class Feature(Document):
    name = StringField(required=True)
    slug = StringField(required=True)
    project = ReferenceField(Project)

    def __str__(self):
        return "Feature %s" % self.name
