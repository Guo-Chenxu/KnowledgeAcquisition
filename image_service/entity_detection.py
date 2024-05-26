import spacy
import jieba.posseg as pseg

# Initialize English and Chinese NLP pipelines
nlp_en = spacy.load("en_core_web_sm")


def entity_detect(text: str, language: str) -> str:
    if language == "en":
        entities = en_entity_detect(text)
        return entities
    elif language == "cn":
        entities = cn_entity_detect(text)
        return entities


def en_entity_detect(text: str) -> str:
    entities = []
    doc = nlp_en(text)
    # Extract entities
    for entity in doc.ents:
        entities.append(
            {
                "text": entity.text,
                "label": entity.label_,
                # "start": entity.start_char,
                # "end": entity.end_char,
            }
        )

    return entities


def cn_entity_detect(text: str) -> str:
    entities = []
    words = pseg.cut(text)
    # Extract entities
    for word, flag in words:
        entities.append({"text": word, "label": flag})

    return entities
