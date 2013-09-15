angular.module( 'rosy.lab', [
  'ui.router',
  'firebase'
])

.config(function config($stateProvider) {
  $stateProvider.state( 'lab', {
    url: '/lab',
    views: {
      'main': {
        controller: 'LabCtrl',
        templateUrl: 'lab/lab.tpl.html'
      }
    },
    data: { pageTitle: 'Lab' }
  });
})

.controller( 'LabCtrl', ['$scope', '$http', 'angularFire', function LabCtrl($scope, $http, angularFire) {
  var ref = new Firebase('https://rosy.firebaseio.com/lab');

  $scope.readonly = true;

  $scope.user.then(function(data) {
    var user = data.data;

    if (user.type === 'teacher') {
      $scope.control = user.id;
      $scope.readonly = false;
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

  angularFire(ref, $scope, "code");
}])

;
