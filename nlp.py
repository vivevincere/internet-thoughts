from google.cloud import language_v1
from enum import Enum

SENTIMENT = Enum("SENTIMENT", "positive negative mixed neutral")
SENTIMENT_THRESHOLD = { SENTIMENT.positive: 0.25, SENTIMENT.negative: -0.25} # need to tune
SENTIMENT_SCORE = "sentiment_score" # score range: normalized -1 to 1
SENTIMENT_MAGNITUDE = "sentiment_magnitude" # magnitude range: 0 to inf
NEUTRAL_THRESHOLD = 1.0 # if magnitude above threshold, considered "mixed" instead of "neutral"

def analyze_sentiment(text_content: str) -> dict():
    """
    Analyzing Sentiment in a String

    Args:
      text_content The text content to analyze

    Return:
      map of SENTIMENT_SCORE and SENTIMENT_MAGNITUDE
    """

    client = language_v1.LanguageServiceClient()

    # Available types: PLAIN_TEXT, HTML
    type_ = language_v1.Document.Type.PLAIN_TEXT

    # Optional. If not specified, the language is automatically detected.
    # For list of supported languages:
    # https://cloud.google.com/natural-language/docs/languages
    document = {"content": text_content, "type_": type_}

    # Available values: NONE, UTF8, UTF16, UTF32
    encoding_type = language_v1.EncodingType.UTF8

    response = client.analyze_sentiment(request = {'document': document, 'encoding_type': encoding_type})
    # Get overall sentiment of the input document
    print(u"Document sentiment score: {}".format(response.document_sentiment.score))
    print(
        u"Document sentiment magnitude: {}".format(
            response.document_sentiment.magnitude
        )
    )
    # Get sentiment for all sentences in the document
    for sentence in response.sentences:
        print(u"Sentence text: {}".format(sentence.text.content))
        print(u"Sentence sentiment score: {}".format(sentence.sentiment.score))
        print(u"Sentence sentiment magnitude: {}".format(sentence.sentiment.magnitude))

    # Get the language of the text, which will be the same as
    # the language specified in the request or, if not specified,
    # the automatically-detected language.
    print(u"Language of the text: {}".format(response.language))
    
    # return { SENTIMENT_SCORE : response.document_sentiment.score, SENTIMENT_MAGNITUDE : response.document_sentiment.magnitude }
    return (response.document_sentiment.score, response.document_sentiment.magnitude)


def get_sentiment_score(res: tuple()) -> float:
    return res[0]
    # return res.get(SENTIMENT_SCORE, None)

def get_sentiment_magnitude(res: dict()) -> float:
    return res[1]
    # return res.get(SENTIMENT_MAGNITUDE, None)

def get_sentiment(res: tuple()) -> SENTIMENT:
    score = get_sentiment_score(res)
    magnitude = get_sentiment_magnitude(res)
    if score >= SENTIMENT_THRESHOLD[SENTIMENT.positive]:
        return SENTIMENT.positive
    elif score <= SENTIMENT_THRESHOLD[SENTIMENT.negative]:
        return SENTIMENT.negative
    elif magnitude >= NEUTRAL_THRESHOLD:
        return SENTIMENT.mixed
    else:
        return SENTIMENT.neutral
