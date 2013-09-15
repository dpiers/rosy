angular.module( 'rosy.lecture', [
  'ui.router',
  'firebase'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'lecture', {
    url: '/lecture',
    views: {
      'main': {
        controller: 'LectureCtrl',
        templateUrl: 'lecture/lecture.tpl.html'
      }
    },
    data: { pageTitle: 'Lecture' }
  });
})

.controller( 'LectureCtrl', ['$scope', '$http', 'angularFire', function LectureCtrl($scope, $http, angularFire) {

  var refCode = new Firebase('https://rosy.firebaseio.com/lecture/code');
  var refOutput = new Firebase('https://rosy.firebaseio.com/lecture/output');
  var refControl = new Firebase('https://rosy.firebaseio.com/lecture/control');

  $scope.readonly = true;
  $scope.connectedStudents = [];

  $scope.isTeacher = function(type) {
    return (type === 'teacher');
  };

  $scope.students = [];

  $scope.$watch('control', function() {
    if (!angular.isDefined('$scope.user.data.id')) {
      return;
    }
    if (!$scope.userData) {
      return;
    }
    var uid = $scope.control;
    if ($scope.userData.type === 'student') {
      if (uid === $scope.userData.id) {
        $scope.readonly = false;
      } else {
        $scope.readonly = true;
      }
    }
  });

  $scope.user.then(function(data) {
    var user = data.data;
    $scope.userData = data.data;

    if (user.type === 'teacher') {
      $scope.control = user.id;
      $scope.readonly = false;

      $http.get('/students').
        success(function(data) {
          $scope.students = [user].concat(data.students);
          console.log($scope.students);
          $scope.control = $scope.students[0].id; // actually the teacher
      });
    } else {
      $scope.connectedStudents.push(user);
    }
  });

  $scope.languages = [
    {name: 'Ruby'},
    {name: 'Python'},
    {name: 'JavaScript'},
    {name: 'Haskell'},
    {name: 'C++'},
    {name: 'Go'}
  ];
  $scope.language = $scope.languages[1];

  $scope.code = 'print \'hello\'';
  $scope.output = 'Waiting for code...';

  $scope.editorOptions = {
    theme: 'tomorrow',
    mode: $scope.language.name.toLowerCase()
  };

  $scope.runCode = function(code, language) {
    var data = JSON.stringify({code: code});
    $http.post('/eval/' + language.name.toLowerCase(), data).
      success(function(data) {
        console.log(data);
        $scope.output = data.output;
      }).
      error(function(data) {
        console.log(data);
        $scope.output = 'error running code';
      });
  };

  angularFire(refCode, $scope, 'code');
  angularFire(refOutput, $scope, 'output');
  angularFire(refControl, $scope, 'control');
}])

;
