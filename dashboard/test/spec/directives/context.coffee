'use strict'

describe 'Directive: context', ->

  # load the directive's module
  beforeEach module 'dashboardApp'

  scope = {}

  beforeEach inject ($controller, $rootScope) ->
    scope = $rootScope.$new()

  it 'should make hidden element visible', inject ($compile) ->
    element = angular.element '<context></context>'
    element = $compile(element) scope
    expect(element.text()).toBe 'this is the context directive'
