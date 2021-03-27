import pandas as pd

all_categories = ["trust", "pessimism", "love", "surprise", "fear", "joy", "disgust", "sadness", "anticipation", "optimism", "anger"]
categories = ["pessimism", "love", "surprise", "fear", "joy", "disgust", "sadness", "optimism", "anger"]
# ordered by least to most frequently appearing - to maintain size when stripping duplicates

def process_csv(filepath: str) -> None:
    data = pd.read_csv(filepath)
    for cat in categories:
        split_df = data[data[cat] == 1]["Tweet"]
        split_df = split_df.to_frame()
        split_df["label"] = cat
        split_df.to_csv("./data/" + cat + ".csv", index=False)
    newdata = pd.DataFrame()
    for cat in categories:
        path = "./data/" + cat + ".csv"
        data = pd.read_csv(path)
        newdata = pd.concat([data, newdata])
    newdata.drop_duplicates(subset=['Tweet'], keep='first', inplace=True)
    newdata.to_csv("./data/all.csv", index=False)



if __name__ == "__main__":
    process_csv("./data/dataset.csv")
    