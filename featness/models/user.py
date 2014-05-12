from mongoengine import Document, StringField, BooleanField


class User(Document):
    name = StringField(required=True)
    email = StringField(required=True)
    password = StringField(required=True)
    active = BooleanField(required=True, default=True)

    def is_authenticated(self):
        return True

    def is_active(self):
        return self.active

    def is_anonymous(self):
        return False

    def get_id(self):
        return str(self.id)

    def __str__(self):
        return self.name
