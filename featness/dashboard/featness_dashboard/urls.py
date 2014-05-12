from django.conf.urls import patterns, include, url

from django.contrib import admin
admin.autodiscover()

from mongoadmin import site

urlpatterns = patterns('',
    url(r'^admin/', include(site.urls)),
)
