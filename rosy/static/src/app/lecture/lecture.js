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
  var refStudents = new Firebase('https://rosy.firebaseio.com/lecture/students');

  $scope.readonly = true;
  $scope.connectedStudents = [];

  $scope.isTeacher = function(type) {
    return (type === 'teacher');
  };

  $scope.user.then(function(data) {
    var user = data.data;

    if (user.type === 'teacher') {
      $scope.control = user.id;
      $scope.readonly = false;

      $http.get('/students').
        success(function(data) {
          console.log(data);
          $scope.students = [user].concat(data.students);
          $scope.controlling = $scope.students[0];
      });
    } else {
      $scope.connectedStudents.push(user);
    }
  });

  $scope.languages = [
    {name: 'Ruby'},
    {name: 'Python'},
    {name: 'JavaScript'}
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

  $scope.synced = {
    code: $scope.code,
    control: $scope.control
  };

  angularFire(refCode, $scope, 'code');
  angularFire(refOutput, $scope, 'output');
  //angularFire(refStudents, $scope, 'students');
}])

;
