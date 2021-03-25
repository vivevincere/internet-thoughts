from multiprocessing.pool import Pool
from collections import defaultdict
from nlp import get_sentiment, get_sentiment_score, analyze_sentiment, SENTIMENT_SCORE, SENTIMENT_MAGNITUDE
import youtube_api

BATCH_SIZE = 100

def analyze_sentiment_batch(text_batch):
    assert len(text_batch) == BATCH_SIZE
    # counts = dict.fromkeys([e.name for e in SENTIMENT], 0)
    counts = defaultdict(int)
    total_score = 0.0
    with Pool(processes=BATCH_SIZE) as pool:
        # returns a list of sentiment names and scores for each list item
        res = pool.map(analyze_sentiment, text_batch)
        for sentiment in res:
            if not sentiment: continue
            counts[get_sentiment(sentiment).name] += 1
            total_score += get_sentiment_score(sentiment)
    return counts, total_score

def analyze_sentiment_list1(text_list):
    batch = []
    # counts = dict.fromkeys([e.name for e in SENTIMENT], 0)
    counts = defaultdict(int)
    tot = 0
    length = 0
    for l in text_list:
        length += 1
        batch.append(l)
        if len(batch) == BATCH_SIZE:
            cur_counts, cur_tot = analyze_sentiment_batch(batch)
            for c in cur_counts:
                counts[c] += cur_counts[c]
            tot += cur_tot
            batch = []
    return { "counts" : counts, "average" : tot / length }
        



def calcNum(n):#some arbitrary, time-consuming calculation on a number
    print(f"Calcs Started on {n}")
    m = n
    for i in range(5000000):
        m += i%25
    if m > n*n:
        m /= 2
    return m

if __name__ == "__main__":
    # result = 0
    str_list = list(map(lambda x: x[1], youtube_api.getCommentsFromVideos("BTS", 10, 10, "en")))
    print(analyze_sentiment_list1(str_list))

    # with Pool(processes=BATCH_SIZE) as pool:
    #     # pool.map(f, text_batch)
    #     nums = [12,25,76,38,8,2,5]
    #     finList = []
    #     result = pool.map(calcNum, nums)
    # print(result)
    # print(sum(result))
