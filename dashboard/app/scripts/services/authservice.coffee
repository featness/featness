'use strict'

class AuthService
  constructor: ->

  isAuthenticated: ->
    return false

angular.module('dashboardApp')
  .service 'AuthService', ->
    return new AuthService()
