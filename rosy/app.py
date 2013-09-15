from flask import Flask, send_from_directory, render_template, request, \
    redirect, session, jsonify, make_response
from flask.ext.sqlalchemy import SQLAlchemy  # pylint: disable=E0611
import requests
import json


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:////tmp/test.db'
app.secret_key = 'THUPER THECRET'
db = SQLAlchemy(app)
EVAL_URL = 'http://tryrosy.com:9000'


class Assignment(db.Model):  # ???
    id = db.Column(db.Integer, primary_key=True)
    student_id = db.Column(db.Integer, db.ForeignKey('student.id'))
    student = db.relationship('Student',
                           backref=db.backref('assignments', lazy='dynamic'))
    teacher_id = db.Column(db.Integer, db.ForeignKey('teacher.id'))
    teacher = db.relationship('Teacher',
                           backref=db.backref('assignments', lazy='dynamic'))
    title = db.Column(db.String(128))
    description = db.Column(db.Text)
    language = db.Column(db.String(32))
    code = db.Column(db.Text)
    output = db.Column(db.Text)
    attempts = db.Column(db.Integer)
    complete = db.Column(db.Boolean)

    def to_json(self):
        return {
            'id': self.id,
            'teacher': self.teacher.id,
            'title': self.title,
            'description': self.description,
            'language': self.language,
            'code': self.code,
            'attempts': self.attempts,
            'complete': self.complete
            }


class User(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    email = db.Column(db.String(120), unique=True)
    # XXX TODO NO THIS ISN'T HASHED CHANGE IT ASAP
    password = db.Column(db.String(128))

    def __init__(self, email, password):
        self.email = email
        self.password = password

    def __repr__(self):
        return '<User: %r>' % self.email

    def to_json(self):
        user_type = 'student'
        teacher = Teacher.query.filter_by(user=self).one()
        if teacher:
            user_type = 'teacher'
        return {
            'id': self.id,
            'email': self.email,
            'type': user_type
            }


class Student(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship("User", backref=db.backref("student", uselist=False))
    teacher_id = db.Column(db.Integer, db.ForeignKey('teacher.id'))
    teacher = db.relationship("Teacher", backref=db.backref("students", lazy='dynamic'))


class Teacher(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship("User", backref=db.backref("teacher", uselist=False))


def get_user_from_session(sesh):
    email = sesh.get('email')
    if email is None:
        return None
    return User.query.filter_by(email=email).one()


def eval_code(code, language):
    print 'submitting'
    url = EVAL_URL + '/' + language
    r = requests.post(url, data={'input': code})
    return r.text.strip()


@app.route('/')
def index():
    print session.get('email')
    return send_from_directory('static/build', 'index.html')


@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        # do login
        email = request.form.get('email')
        u = User.query.filter_by(email=email).first()
        if u and request.form.get('password') == u.password:
            session['email'] = email
            return redirect('/')
    return render_template('login.html')


@app.route('/logout')
def logout():
    session.pop('email', None)
    return redirect('/')


@app.route('/user')
def user():
    u = get_user_from_session(session)
    if u is None:
        return jsonify({'user': None})
    return jsonify({'user': u.to_json()})


@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        # TODO validation
        # XXX TODO HASH PASSWORDS
        u = User(request.form.get('email'), request.form.get('password'))
        db.session.add(u)
        if request.form.get('type') == 'teacher':
            teacher = Teacher()
            teacher.user = u
            db.session.add(teacher)
        else:
            student = Student()
            student.user = u
            db.session.add(student)
        db.session.commit()

        session['email'] = u.email
        return redirect('/')
    return render_template('register.html')


@app.route('/eval/<language>', methods=['POST'])
def eval_route(language):
    data = json.loads(request.data)
    output = eval_code(data.get('code'), language)
    return jsonify({'output': output})


@app.route('/assignments')
def list_assignments():
    u = get_user_from_session(session)
    if u is None:
        assignments = []
    else:
        if u.teacher:
            teacher = u.teacher
        else:
            teacher = u.student.teacher
        assignments = [a.to_json() for a in teacher.assignments.all()]
    return jsonify({'assignments': assignments})


@app.route('/assignment/<aid>')
def assignment_detail(aid):
    a = Assignment.query.filter_by(id=aid).one()
    return jsonify(a.to_json())


@app.route('/assignment/<aid>/submit', methods=['POST'])
def submit_assignment(aid):
    assignment = Assignment.query.filter_by(id=aid).one()
    assignment.attempts += 1
    data = json.loads(request.data)
    output = eval_code(data.get('code'), assignment.language)
    if output == assignment.output:
        assignment.complete = True
        correct = True
    else:
        correct = False
    db.session.add(assignment)
    db.session.commit()
    return jsonify({'correct': correct, 'output': output, 'attempts': assignment.attempts})


@app.route('/assignments/new', methods=['POST'])
def new_assignment():
    user = get_user_from_session(session)
    teacher = Teacher.query.filter_by(user=user).one()
    if not teacher:
        return make_response('', 404)
    data = json.loads(request.data)
    for student in teacher.students.all():
        a = Assignment()
        a.title = data.get('title', '')
        a.language = data.get('language', '').lower()
        a.description = data.get('description', '')
        a.output = data.get('output', '')
        a.code = data.get('code', '')
        a.student = student
        a.teacher = teacher
        a.attempts = 0
        a.complete = False
        db.session.add(a)
    db.session.commit()

    return jsonify(a.to_json())


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
