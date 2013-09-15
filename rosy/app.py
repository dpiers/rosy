from flask import Flask, send_from_directory, render_template, request, redirect, session, jsonify
from flask.ext.sqlalchemy import SQLAlchemy  # pylint: disable=E0611
import requests


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:////tmp/test.db'
app.secret_key = 'THUPER THECRET'
db = SQLAlchemy(app)
EVAL_URL = 'http://tryrosy.com:9000'


class Assignment(db.Model):  # ???
    id = db.Column(db.Integer, primary_key=True)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship('User', backref=db.backref('assignments', lazy='dynamic'))
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


def get_user_from_session(session):
    email = session.get('email')
    if email is None:
        return None
    return User.query.filter_by(email=email).one()

def eval_code(code):
    r = requests.post(EVAL_URL + '/python', data={'code': code})
    print r.text
    return r.text

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
    return jsonify({'user': {'email': u.email, 'id': u.id}})

@app.route('/register', methods=['GET', 'POST'])
def register():
    if request.method == 'POST':
        # TODO validation
        # XXX TODO HASH PASSWORDS
        u = User(request.form.get('email'), request.form.get('password'))
        db.session.add(u)
        db.session.commit()

        session['email'] = u.email
        return redirect('/')
    return render_template('register.html')

@app.route('/assignments')
def assignments():
    u = get_user_from_session(session)
    if u is None:
        assignments = []
    else:
        assignments = [a.to_json() for a in u.assignments.all()]
    return jsonify({'assignments': assignments})

@app.route('/assignment/<id>')
def assignment(id):
    a = Assignment.query.filter_by(id=id).one()
    return jsonify(a.to_json())

@app.route('/assignment/<id>/submit', methods=['POST'])
def submit_assignment(id):
    assignment = Assignment.query.filter_by(id=id).one()
    assignment.attempts += 1
    output = eval_code(request.form.get('code'))
    if output == assignment.output:
        assignment.complete = True
        correct = True
    else:
        correct = False
    db.session.add(assignment)
    db.session.commit()
    return jsonify({'correct': correct, 'output': output})

if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
