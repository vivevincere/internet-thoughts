import os 
nlp_key_path = os.path.join(os.path.abspath("."), "gkey.json")
print(nlp_key_path) # use this to get the absolute path on your local machine



# todo: need to set env on gcp - can't set within python script?
'''
import subprocess
command = f"export GOOGLE_APPLICATION_CREDENTIALS={nlp_key_path}".split()
print(command)
process = subprocess.Popen(command, stdout=subprocess.PIPE)
output, error = process.communicate()
print(output)
'''

# Imports the Google Cloud client library
from google.cloud import language_v1


# Instantiates a client
client = language_v1.LanguageServiceClient()

# The text to analyze
text = u"Hello, world!"
document = language_v1.Document(content=text, type_=language_v1.Document.Type.PLAIN_TEXT)

# Detects the sentiment of the text
sentiment = client.analyze_sentiment(request={'document': document}).document_sentiment

print("Text: {}".format(text))
print("Sentiment: {}, {}".format(sentiment.score, sentiment.magnitude))

