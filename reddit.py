import praw

reddit = praw.Reddit("internet_thoughts")


def search_reddit(query, limit=1000):
    comments_list = []
    count = 0
    for submission in reddit.subreddit("all").search(query, time_filter="week"):
        submission.comment_sort = "top"
        submission.comments.replace_more(limit=None)
        for top_level_comment in submission.comments:
            comments_list.append((top_level_comment.score, top_level_comment.body))
            count += 1
            if count > limit:
                break
    return comments_list

print(search_reddit("bts", 20))

