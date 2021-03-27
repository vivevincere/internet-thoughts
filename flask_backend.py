import json

from flask import Flask, request, jsonify

app = Flask(__name__)

#Sentiment count


@app.route('/', methods = ['POST'])
def sentiment_request():
	record = json.loads(request.data)
	searchTerm = record["searchTerm"]
	print(searchTerm)
	return jsonify(searchTerm)
#word cloud



#related searches


#relevant tweets

