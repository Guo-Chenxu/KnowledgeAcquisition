# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://docs.scrapy.org/en/latest/topics/item-pipeline.html


# useful for handling different item types with a single interface
import json
from itemadapter import ItemAdapter


class RustbookPipeline:
    def __init__(self):
        self.book = None

    def open_spider(self, spider):
        try:
            self.book = open('rust_book.json', 'w', encoding='utf-8')
        except Exception as e:
            print(e)

    def process_item(self, item, spider):
        dict_item = dict(item)
        json_str = json.dumps(dict_item, ensure_ascii=False) + "\n"
        self.book.write(json_str)
        return item

    def close_spider(self, spider):
        self.book.close()

