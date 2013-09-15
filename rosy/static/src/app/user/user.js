angular.module( 'rosy.user', [
  'ui.router'
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
  $scope.user.then(function(data) {
    var user = data.data;

    $http.get('/assignments').
      success(function(data) {
        $scope.assignments = data.assignments;
    });
    if (user.type === 'teacher') {
      $http.get('/students').
        success(function(data) {
          $scope.students = data.students;
      });
    }

  });

  $scope.numAssignments = function(student) {
    return student.assignments.length;
  };

  $scope.numCompleted = function(student) {
    var completed = 0;
    student.assignments.forEach(function(assn) {
      if (assn.complete) {
        completed++;
      }
    });
    return completed;
  };
}])

;
