from flask import Flask, send_from_directory, render_template, request, redirect, session, jsonify
from flask.ext.sqlalchemy import SQLAlchemy  # pylint: disable=E0611


app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:////tmp/test.db'
app.secret_key = 'THUPER THECRET'
db = SQLAlchemy(app)


class Assignment(db.Model):  # ???
    id = db.Column(db.Integer, primary_key=True)
    description = db.Column(db.Text)
    user_id = db.Column(db.Integer, db.ForeignKey('user.id'))
    user = db.relationship('User', backref=db.backref('assignments', lazy='dynamic'))


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
        assignments = [{'id': a.id, 'description': a.description} for a in u.assignments.all()]
    return jsonify({'assignments': assignments})

if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
