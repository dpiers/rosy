angular.module( 'rosy', [
  'templates-app',
  'templates-common',
  'rosy.home',
  'rosy.about',
  'rosy.try',
  'rosy.user',
  'rosy.assignment',
  'ui.state',
  'ui.route',
  'ui.ace'
])

.config( function myAppConfig ( $stateProvider, $urlRouterProvider ) {
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
      if (toState.name === 'home') {
        $http.get('/user').
          success(function(data) {
            if (data.user) {
              event.preventDefault();
              $state.transitionTo('user', { location: true, inherit: true, relative: $state.$current });
            }
          });
      }
  });
})

;

