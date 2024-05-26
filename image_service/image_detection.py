import json
import io
from keras.applications.resnet import (
    ResNet50,
    preprocess_input,
    decode_predictions,
)
from keras.preprocessing import image
from keras.utils import img_to_array
from PIL import Image
import numpy as np
import logging
import tensorflow as tf

config = tf.compat.v1.ConfigProto(
    gpu_options=tf.compat.v1.GPUOptions(allow_growth=True))
sess = tf.compat.v1.Session(config=config)

log = logging.getLogger("ImageToKeywords")

# Image object detection model
model = ResNet50(weights="imagenet")


def image_to_keywords(file: str) -> tuple[str, int]:
    if file.filename == "":
        return ("No selected file", 400)
    log.info(file.filename)

    # Open image
    try:
        img = (
            Image.open(io.BytesIO(file.read()))
            .convert("RGB")
            .resize((224, 224))
        )
    except Exception as e:
        msg = "Failed to load image: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image to array
    try:
        x = img_to_array(img)
    except Exception as e:
        msg = "Failed to convert image to array: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image preprocess
    try:
        x = np.expand_dims(x, axis=0)
        x = preprocess_input(x)
    except Exception as e:
        msg = "Failed to preprocess image: " + str(e)
        log.error(msg)
        return (msg, 500)

    # Image detection
    try:
        preds = model.predict(x)
        predictions = decode_predictions(preds, top=5)[0]
        # Return top three predictions or all if predictions are less than three
        if len(predictions) >= 3:
            keywords = [pred[1] for pred in predictions[:3]]
        else:
            keywords = [pred[1] for pred in predictions]
        keywords = " ".join(kw for kw in keywords)
        log.info(keywords)
        return (json.dumps({"keyword": keywords}), 200)
    except Exception as e:
        msg = "Failed to process image: " + str(e)
        log.error(msg)
        return (msg, 500)
