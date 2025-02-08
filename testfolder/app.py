from flask import Flask, render_template, request
import random
import string

app = Flask(__name__)

@app.route("/", methods=["GET", "POST"])
def index():
    password = ""
    if request.method == "POST":
        length = int(request.form["length"])
        include_numbers = "numbers" in request.form
        include_symbols = "symbols" in request.form
        characters = string.ascii_letters
        if include_numbers:
            characters += string.digits
        if include_symbols:
            characters += string.punctuation
        password = ''.join(random.choice(characters) for i in range(length))
    return render_template("index.html", password=password)

if __name__ == "__main__":
    app.run(debug=True)