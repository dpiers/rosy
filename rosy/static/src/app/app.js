angular.module( 'rosy', [
  'templates-app',
  'templates-common',
  'rosy.home',
  'rosy.about',
  'rosy.try',
  'rosy.user',
  'rosy.assignment',
  'rosy.newAssignment',
  'ui.state',
  'ui.route',
  'ui.ace'
])

.config( function myAppConfig ( $stateProvider, $urlRouterProvider, $httpProvider ) {
  $urlRouterProvider.otherwise( '/home' );
})

.run( function run () {
})

.controller( 'AppCtrl', function AppCtrl ( $scope, $location, $http, $state ) {
  $scope.$on('$stateChangeSuccess', function(event, toState, toParams, fromState, fromParams){
    if ( angular.isDefined( toState.data.pageTitle ) ) {
      $scope.pageTitle = toState.data.pageTitle + ' | rosy' ;
    }
  });
  $scope.$on('$stateChangeStart',
    function(event, toState, toParams, fromState, fromParams){
      $http.get('/user').
        success(function(data) {
          $scope.user = data.user;
          if (data.user) {
            if (toState.name === 'home') {
              event.preventDefault();
              $state.transitionTo('user', { location: true, inherit: true, relative: $state.$current });
            }
          }
        });
  });
})

;

