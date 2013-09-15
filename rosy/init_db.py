from app import *


def create_db():
    db.drop_all()
    db.create_all()

def fixture():
    teacheru = User('teacher@tryrosy.com', 'test')
    db.session.add(teacheru)
    teacher = Teacher()
    teacher.user = teacheru
    db.session.add(teacher)

    studentu = User('student@tryrosy.com', 'test')
    db.session.add(studentu)
    student = Student()
    student.user = studentu
    student.teacher = teacher
    db.session.add(student)

    a = Assignment()
    a.teacher = teacher
    a.title = 'First Assignment'
    a.description = 'Assign the variable a to the number 3.\n' \
            'Assign the variable b to the number 4.'
    a.language = 'python'
    a.code = 'a = 0\nb = 1\n\nprint a, b'
    a.output = '3 4'
    a.attempts = 0
    a.complete = False
    db.session.add(a)

    db.session.commit()

if __name__ == '__main__':
    create_db()
    fixture()
