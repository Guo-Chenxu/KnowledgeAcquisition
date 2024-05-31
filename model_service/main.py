from string import punctuation
from flask import Flask, request
import json
import logging
import entity_detection
import image_detection
from collections import Counter
from nltk.corpus import stopwords
from nltk.tokenize import word_tokenize
import re

app = Flask(__name__)
log = app.logger
log.setLevel(logging.DEBUG)


@app.route("/image_to_keywords", methods=["POST"])
def image_to_keywords():
    # Check if the post request has the file part
    if "file" not in request.files:
        return "No file part", 400
    file = request.files["file"]

    # if user does not select file, browser submits an empty part without
    # filename
    if file.filename == "":
        return "No selected file", 400
    log.info(file.filename)
    result, code = image_detection.image_to_keywords(file)
    if code != 200:
        log.error(result)
        return result, code
    return result, code


@app.route("/extract_info", methods=["POST"])
def extract_info():
    data = request.get_json()
    text = data.get("text")
    language = data.get("language")
    log.info("Data language: " + language)

    if not text or not language:
        return "Invalid request: no text or no language", 400

    if language not in ["en", "cn"]:
        return "Unsupported language: " + language, 400

    # Entity detection
    entities = entity_detection.entity_detect(text, language)

    # Extract hot words
    stop_words = (
        set(stopwords.words("english")) if language == "en" else set()
    )  # Add Chinese stopwords if needed
    word_tokens = word_tokenize(text)
    words = [
        w for w in word_tokens if not w in stop_words and not w in punctuation
    ]

    hot_words = dict(Counter(words).most_common(5))

    entities = [e["text"] for e in entities]
    entities = dict(Counter(entities).most_common(5))
    entities = {k: v for k, v in entities.items() if v > 1}

    jsonResponse = json.dumps({"entities": entities, "hot_words": hot_words})
    log.debug(jsonResponse)
    return jsonResponse
    
@app.route("/extract_info_regex", methods=["POST"])
def extract_info():
    data = request.get_json()
    text = data.get("text")
    pattern = data.get("pattern")
    language = data.get("language")
    word_class = data.get("word_class")

    if not text or not language:
        return "Invalid request: no text or no language", 400

    if language not in ["en", "cn"]:
        return "Unsupported language: " + language, 400

    # Entity detection
    entities = entity_detection.entity_detect(text, language)

    # Extract words with regex
    words = []
    entities = [{"text": item["text"], "label": item["label"]} for item in entities]
    for entity in entities:
        if word_class == entity.get("label") and bool(re.fullmatch(pattern=pattern,string=entity.get('text'))):
            words.append(entity.get('text'))

    jsonResponse = json.dumps({"words": words})
    log.debug(jsonResponse)
    return jsonResponse


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=39840, debug=True)
