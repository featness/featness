'use strict'

class AuthService
  constructor: ->

  isAuthenticated: ->
    return false

  authenticateWithGoogle: (callback) ->
    gapi.auth.signIn(
      callback: (authResult) =>
        gapi.client.load('plus','v1', @handleProfileLoad(authResult, callback))
    ) # Will use page level configuration

  handleProfileLoad: (authResult, callback) ->
    return =>
      request = gapi.client.plus.people.get( {'userId' : 'me'} )
      request.execute(@handleProfileLoaded(authResult, callback))

  handleProfileLoaded: (authResult, callback) ->
    return (profile) =>
      callback(authResult, profile)

angular.module('dashboardApp')
  .service 'AuthService', ->
    return new AuthService()
