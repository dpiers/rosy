angular.module( 'rosy.user', [
  'ui.state'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'user', {
    url: '/user',
    views: {
      'main': {
        controller: 'UserCtrl',
        templateUrl: 'user/user.tpl.html'
      }
    },
    data: { pageTitle: 'Dashboard' }
  });
})

.controller( 'UserCtrl', ['$scope', '$http', function UserCtrl($scope, $http) {
   $http.get('/user').
     success(function(data) {
       $scope.user = data.user;
     });
}])

;
