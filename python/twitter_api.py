import requests
import os
import json

bearerToken = ""




#by date

#by location

#by keyword

#number of tweets

#returns number of likes

#returns number of retweets

#
def create_url(searchTerm, location, start_time, end_time, numberOfTweets, next_token):
    query = searchTerm
    # Tweet fields are adjustable.
    # Options include:
    # attachments, author_id, context_annotations,
    # conversation_id, created_at, entities, geo, id,
    # in_reply_to_user_id, lang, non_public_metrics, organic_metrics,
    # possibly_sensitive, promoted_metrics, public_metrics, referenced_tweets,
    country_fields = ""

    start_field = ""
    end_field = ""
    # source, text, and withheld
    tweet_fields = "tweet.fields=geo,created_at,public_metrics&place.fields=country,name"

    if end_time:
    	end_field = "&end_time=" + end_time

    	tweet_fields = tweet_fields + end_field
    if start_time:
    	start_field = "&start_time=" + start_time

    	tweet_fields = tweet_fields + start_field

    if location != None:
    	country_fields = " -place:\"" + location + '\"'

    if next_token:
    	tweet_fields += "&next_token=" + next_token

    tweet_fields += "&max_results=" + str(numberOfTweets)
    #query += country_fields






    url = "https://api.twitter.com/2/tweets/search/recent?query={}&{}".format(
        query, tweet_fields
    )

    return url


def create_headers(bearer_token):
    headers = {"Authorization": "Bearer {}".format(bearer_token)}
    return headers






def connect_to_endpoint(url, headers):
    response = requests.request("GET", url, headers=headers)
    print(response.status_code)
    if response.status_code != 200:
        raise Exception(response.status_code, response.text)
    return response.json()


def searchTweets(numberOfTweets, searchTerm, location, start_time, end_time):
	bearer_token = bearerToken
	thelist = []
	next_token  = None
	while numberOfTweets > 0:
		curNum = min(100,numberOfTweets)
		url = create_url("BTS",location,start_time,end_time, curNum,next_token)
		headers = create_headers(bearer_token)
		response = connect_to_endpoint(url, headers)

		#response = json.dumps(json_response, indent=4, sort_keys=True)
		next_token = response["meta"]["next_token"]
		thelist +=  response["data"]
		numberOfTweets -= curNum
	return thelist

print(searchTweets(20,"BTS", None, None, "2021-03-24T09:55:06-00:00"))