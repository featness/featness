'use strict'

class AuthService
  constructor: (@http) ->
    @storage = window.sessionStorage

  isAuthenticated: ->
    return @getToken()?
  
  getToken: ->
    return @storage.getItem("featness-token")

  getAccount: ->
    return @storage.getItem("featness-account")

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
      @getAuthenticationHeader(authResult, profile, callback)

  getAuthenticationHeader: (authResult, profile, callback) ->
    userAccount = ""
    for email in profile.emails
      if email.type == "account"
        userAccount = email.value
        break

    @http(
      url: "http://local.featness.com:8000/authenticate/google",
      method: "POST",
      headers: {
        'X-Auth-Data': "#{profile.emails[0].value};#{authResult.access_token}"
      }
      data: {}
    ).success((data, status, headers, config) =>
      token = headers('X-Auth-Token')
      @storage.setItem("featness-token", token)
      @storage.setItem("featness-account", userAccount)
      callback(userAccount, token)
    ).error((data, status, headers, config) =>
      callback(null, null)
    )

angular.module('dashboardApp')
  .service 'AuthService', ($http) ->
    return new AuthService($http)
