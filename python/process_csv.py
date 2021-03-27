import pandas as pd

all_categories = ["trust", "pessimism", "love", "surprise", "fear", "joy", "disgust", "sadness", "anticipation", "optimism", "anger"]
categories = ["surprise", "fear", "joy", "disgust", "sadness", "anger"]
# ordered by least to most frequently appearing - to maintain size when stripping duplicates

def process_csv(filepath: str) -> None:
    data = pd.read_csv(filepath)
    drop_col = [col for col in all_categories if col not in categories]
    data.drop(drop_col, axis=1, inplace=True)
    for cat in categories:
        d = {1: cat, 0: 'None'}
        data[cat] = data[cat].map(d)
    data["label"] = data[categories].agg(','.join, axis=1)
    categories.append("ID")
    data.drop(categories, axis=1, inplace=True)
    data["label"] = [row.replace("None,","") for row in data["label"]]
    data["label"] = [row.replace("None","") for row in data["label"]]
    data.replace("", float("NaN"), inplace=True)
    data.dropna(axis=0, inplace=True, subset=["label"])
    data.to_csv("./data/all.csv", index=False)



if __name__ == "__main__":
    process_csv("./data/dataset.csv")
    