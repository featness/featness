'use strict'

angular.module('dashboardApp')
  .directive('noteams', ->
    restrict: 'E'
    replace: true
    template: '''
      <div class="jumbotron">
        <h1>No teams found...</h1>
        <p class="lead">
          It looks like there are no teams created in <a href="/" class="logo"><strong>feat</strong>ness</a>.
        </p>
        <p class="toolbar">
          <a href="/teams/new" class="btn btn-lg btn-success new-team"><span class="icon"><span class="first glyphicon glyphicon-user"></span><span class="glyphicon glyphicon-user"></span></span> create a new team</a>
        </p>
      </div>
    '''
    link: (scope, element, attrs) ->
  )
