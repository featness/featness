'use strict'

class AuthService
  constructor: (@http, @window) ->
    @storage = @window.sessionStorage

  isAuthenticated: ->
    return @getToken()?

  getToken: ->
    return @storage.getItem("featness-token")

  getUser: ->
    if not @storage.getItem("featness-account")?
      return null

    return {
      account: @storage.getItem("featness-account")
      name: @storage.getItem("featness-account-name")
      picture: @storage.getItem("featness-account-image")
    }

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
      userAccount = ""
      for email in profile.emails
        if email.type == "account"
          userAccount = email.value
          break

      @getAuthenticationHeader("google", userAccount, authResult.access_token, callback)

  initializeFacebook: (handleFacebookCallback) ->
    id = 'facebook-jssdk'
    ref = document.getElementsByTagName('script')[0]

    return if document.getElementById(id)?

    js = document.createElement('script')
    js.id = id
    js.async = true
    js.src = "//connect.facebook.net/en_US/all.js"

    ref.parentNode.insertBefore(js, ref)

    @window.fbAsyncInit = =>
      #Executed when the SDK is loaded
      FB.init(
        appId: '843188275707721',
        channelUrl: 'app/channel.html',
        status: true,
        cookie: true,
        xfbml: true
      )

      FB.Event.subscribe('auth.authResponseChange', @handleFacebookAuthenticationResponseChange)

    @facebookCallback = handleFacebookCallback

  handleFacebookAuthenticationResponseChange: (response) =>
    if (response.status == 'connected')
      ###
      The user is already logged,
      is possible retrieve his personal info
      ###
      FB.api('/me', (meResponse) =>
        @getAuthenticationHeader("facebook", response.authResponse.accessToken, meResponse.username, meResponse.name, "http://graph.facebook.com/#{ meResponse.username }/picture", @facebookCallback)
      )

    else
      ###
      The user is not logged to the app, or into Facebook:
      destroy the session on the server.
      ###
      console.log('unauthenticated')

  authenticateWithFacebook: (callback) ->
    FB.login()

  getAuthenticationHeader: (provider, accessToken, userAccount, name, imageUrl, callback) ->
    @http(
      url: "http://local.featness.com:8000/authenticate/#{ provider }",
      method: "POST",
      headers: {
        'X-Auth-Data': "#{userAccount};#{name};#{imageUrl};#{accessToken}"
      }
      data: {}
    ).success((data, status, headers, config) =>
      token = headers('X-Auth-Token')
      if token?
        @storage.setItem("featness-token", token)
        @storage.setItem("featness-account", userAccount)
        @storage.setItem("featness-account-name", name)
        @storage.setItem("featness-account-image", imageUrl)
        callback(userAccount, name, token)
      else
        callback(null, null, null)
    ).error((data, status, headers, config) =>
      callback(null, null, null)
    )


angular.module('dashboardApp')
  .service 'AuthService', ($http, $window) ->
    return new AuthService($http, $window)
