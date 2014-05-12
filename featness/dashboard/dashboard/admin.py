from mongoadmin import site, DocumentAdmin

from featness.models import Team

class TeamAdmin(DocumentAdmin):
    pass

site.register(Team, TeamAdmin)
