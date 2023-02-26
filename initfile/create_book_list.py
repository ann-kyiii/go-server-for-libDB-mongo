import argparse
import json
import traceback

import isbnlib
import requests
import numpy as np
import pandas as pd
from tqdm import tqdm


def result_isbn_openBD(res):
    if res.text == "[null]":
        return np.nan

    try:
        res = json.loads(res.text)
        title_element = res[0]["onix"]["DescriptiveDetail"]["TitleDetail"]["TitleElement"]
        summary = res[0]["summary"]

        title = title_element["TitleText"]["content"].replace("\\", "")
        subtitle = title_element["Subtitle"]["content"].replace("\\", "") if "Subtitle" in title_element.keys() else ""
        partNumber = title_element["PartNumber"].replace("\\", "") if "PartNumber" in title_element.keys() else ""
        # subtitleやpartNumberが存在する場合は結合する
        bookName = f"{title} {subtitle}" if subtitle else title
        bookName = f"{bookName} {partNumber}" if partNumber else bookName

        cover = summary["cover"].replace("\\", "") if summary["cover"] != "" else np.nan
        return {
            "bookName": bookName,
            "author": summary["author"].replace("\\", ""),
            "publisher": summary["publisher"].replace("\\", ""),
            "pubdate": summary["pubdate"],
            "imgURL": cover,
        }
    except:
        traceback.print_exc()
        print(res.text)


def get_books_from_openBD(book_finds):
    # 存在する本のISBNを取得
    isbns = book_finds[book_finds.exist != np.nan].ISBN.dropna()

    # openBDで書誌情報・書影を取得
    tqdm.pandas(desc="Get book information using openBD")
    openBD = lambda isbn: requests.get(f"https://api.openbd.jp/v1/get?isbn={isbn}")
    isbn_responses = isbns.progress_apply(openBD)
    results = isbn_responses.apply(result_isbn_openBD)

    res_ind = results.index
    res_dict = {
        "bookName": [],
        "author": [],
        "publisher": [],
        "pubdate": [],
        "imgURL": [],
    }
    for res in results.values:
        if isinstance(res, dict):
            res_dict["bookName"] += [res["bookName"]]
            res_dict["author"] += [res["author"]]
            res_dict["publisher"] += [res["publisher"]]
            res_dict["pubdate"] += [res["pubdate"]]
            res_dict["imgURL"] += [res["imgURL"]]
        else:
            res_dict["bookName"] += [np.nan]
            res_dict["author"] += [np.nan]
            res_dict["publisher"] += [np.nan]
            res_dict["pubdate"] += [np.nan]
            res_dict["imgURL"] += [np.nan]

    # bookNameがNaNのものを除外した結果を取得
    res_openBD = pd.DataFrame(res_dict).set_index(res_ind).dropna(subset=["bookName"])
    # 取得した結果で上書き
    book_finds["imgURL"] = np.nan
    for i, item in res_openBD.iterrows():
        for key, val in item.items():
            if isinstance(val, str):
                book_finds.at[i, key] = val

    return book_finds


def fill_nan_with_default(book_finds):
    # genre
    book_finds["genre"] = book_finds["genre"].fillna("unidentified")
    # subGenre
    book_finds["subGenre"] = book_finds["subGenre"].fillna("unidentified")
    # find
    book_finds["find"] = book_finds["find"].fillna(1)
    # sum
    book_finds["sum"] = book_finds["sum"].fillna(1)
    # author
    book_finds["author"] = book_finds["author"].fillna("unidentified")
    # publisher
    book_finds["publisher"] = book_finds["publisher"].fillna("unidentified")
    # pubdate
    book_finds["pubdate"] = book_finds["pubdate"].fillna("unidentified")
    # locateAt4F
    book_finds["locateAt4F"] = book_finds["locateAt4F"].fillna(False)
    book_finds["locateAt4F"][book_finds["locateAt4F"] == "〇"] = True
    # withDisc
    book_finds["withDisc"] = book_finds["withDisc"].fillna("なし")
    # other
    book_finds["other"] = book_finds["other"].fillna("なし")
    # ISBN
    book_finds["ISBN"] = book_finds["ISBN"].fillna("unidentified")
    # exist
    book_finds["exist"] = book_finds["exist"].fillna("unidentified")
    # imgURL
    book_finds["imgURL"] = book_finds["imgURL"].fillna("unidentified")

    return book_finds


def to_isbn_13(isbn):
    isbn_ = isbnlib.canonical(isbn)
    isbn_ = isbn_ if isbnlib.is_isbn13(isbn_) else isbnlib.to_isbn13(isbn_)
    if isbn_ is None or isbn_ == "":
        if isbn != "unidentified" and isbn != "不明":
            print("ISBNを変換できませんでした:", isbn)
        return "unidentified"
    return isbn_


def main(input_file, output_file, start_id):
    # ファイルの読み込み
    book_finds = pd.read_csv(input_file, dtype="object")

    # openBDから書誌情報を取得
    book_finds = get_books_from_openBD(book_finds)

    # NaNをデフォルト値で埋める
    book_finds = fill_nan_with_default(book_finds)

    # ISBNを13桁に統一
    book_finds["ISBN"] = book_finds["ISBN"].map(to_isbn_13)

    # id列を追加
    book_finds.insert(0, "id", range(start_id, start_id + len(book_finds)))

    # MongoDB用に型名を付与
    book_finds.rename(
        columns={
            "id": "id.int64()",
            "bookName": "bookName.string()",
            "genre": "genre.string()",
            "subGenre": "subGenre.string()",
            "ISBN": "ISBN.string()",
            "find": "find.int64()",
            "sum": "sum.int64()",
            "author": "author.string()",
            "publisher": "publisher.string()",
            "pubdate": "pubdate.string()",
            "exist": "exist.string()",
            "locateAt4F": "locateAt4F.boolean()",
            "withDisc": "withDisc.string()",
            "other": "other.string()",
            "imgURL": "imgURL.string()",
        },
        inplace=True,
    )

    # Dataframeのindexを除外してファイルに保存
    book_finds.to_csv(output_file, index=False)
    print("ファイルの作成が完了しました")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--input_file", type=str, default="booklist.csv")
    parser.add_argument("--output_file", type=str, default="initdata.csv")
    parser.add_argument("--start_id", type=int, default=0)
    args = parser.parse_args()

    main(args.input_file, args.output_file, args.start_id)
