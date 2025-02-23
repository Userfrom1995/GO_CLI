from flask import Flask, render_template, request, redirect, url_for, flash
import os
from werkzeug.utils import secure_filename
import sqlite3

app = Flask(__name__)
app.secret_key = os.urandom(24) #this line is important for flash messages to work
UPLOAD_FOLDER = 'uploads'
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
os.makedirs(UPLOAD_FOLDER, exist_ok=True) #ensure the uploads folder exists

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        #handle file upload here.  Will add later
        pass
    return render_template('index.html')

if __name__ == '__main__':
    app.run(debug=True)
