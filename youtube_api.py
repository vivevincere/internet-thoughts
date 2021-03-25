import requests
import os
import heapq
import googleapiclient.discovery




api_service_name = "youtube"
api_version = "v3"
DEVELOPER_KEY = "AIzaSyCtOaqMVQasejuYBG6S2uIp-mGzlZ2YYqM"


#given a searchTerm, returns a list of videoIDs
def searchForVideos(searchTerm, language,numberOfVideos):
	youtube = googleapiclient.discovery.build(
	    api_service_name, api_version, developerKey = DEVELOPER_KEY)
	thelist = []
	nextToken = None


	while numberOfVideos > 0:
		curNum = 50
		if numberOfVideos < curNum:
			curNum = numberOfVideos

		request = youtube.search().list(
		        part="snippet",
		        maxResults = curNum,
		        q= searchTerm,
		        regionCode="US",
		        type = "video",
		        pageToken = nextToken,
		        relevanceLanguage= language
		    )

		response = request.execute()
		nextToken = response['nextPageToken']
		for parentVideo in response['items']:
			thelist.append(parentVideo['id']['videoId'])
		numberOfVideos -= curNum
	return thelist

#given a videoID, retrieves a list of comment tuples [likeCount,comment]
def getCommentThread(videoID, numberOfComments): 
#currently retrieving comments by relevance i.e. the most popular comments, can be changed to the most recent comments
	youtube = googleapiclient.discovery.build(
	    api_service_name, api_version, developerKey = DEVELOPER_KEY)

	thelist = []
	nextToken = None

	while numberOfComments > 0:
		curNum = 100
		if numberOfComments < curNum:
			curNum = numberOfComments
		
		request = youtube.commentThreads().list(
	            part="snippet,id",
	        videoId= videoID,
	        pageToken = nextToken,
	        maxResults = curNum,
	        order = "relevance"
	)

		response = request.execute()

		
		nextToken = response['nextPageToken']
		for parentComment in response['items']:
			unit = []
			unit.append(parentComment['snippet']['topLevelComment']['snippet']['likeCount'])
			unit.append(parentComment['snippet']['topLevelComment']['snippet']['textDisplay'])
			
			thelist.append(unit)
		numberOfComments -= curNum
	return thelist


#Gets a list of comment tuples [likeCount,comment]
#searchTerm is the keyword, videoCount is the number of videos to get comments from, language is the 2 character representation of desired language e.g. "en"
def getCommentsFromVideos(searchTerm, videoCount, commentsPerVideo, language):

	videoList = searchForVideos(searchTerm, language,videoCount)

	commentList = []

	for videoID in videoList:
		comments = getCommentThread(videoID, commentsPerVideo)
		commentList += comments
	return commentList


#Gets the topCommentCount most liked comments from a list of comment tuples [likeCount, comment] 
def getMostLiked(comments, topCommentCount): 

	n = 0
	retList = []

	for x in comments:
		if n < topCommentCount:
			n += 1
			heapq.heappush(retList,x)
		else:
			likes = x[0]
			comment = x[1]
			
			if likes > retList[0][0]:
				heapq.heappop(retList)
				heapq.heappush(retList,x)
	return retList



#comments = getCommentThread("uQYLGiuQqpA",1000)
#print(getMostLiked(comments, 20))
# print(getCommentsFromVideos("BTS", 5, 100, "en"))