import nltk
from nltk.corpus import stopwords
import string


nltk.download('punkt')
nltk.download('stopwords')

def mostCount(words):
	d = {}
	tokens = nltk.word_tokenize(words)
	for w in tokens:
		if w.lower() not in stopwords.words('english') and str(w) not in string.punctuation and w != "https" and len(w) > 2:
			if w in d:
				d[w] += 1
			else:
				d[w] = 1
	a = list(d.items())
	a.sort( key = lambda x: int(x[1]))
	final6 = a[len(a) - 30:]
	return final6


f = open("python/wordcloud.txt", "r")
a = f.read()
toPrint = mostCount(a)
buildString = ""
for word, count in toPrint:
	buildString += word + " " + str(count) + "\n"

print(buildString)