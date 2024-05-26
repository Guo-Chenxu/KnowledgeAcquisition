import datetime
import scrapy
import logging

logger = logging.getLogger("ThebookSpider")
logger.setLevel(logging.DEBUG)


class ThebookSpider(scrapy.Spider):
    name = "thebook"
    allowed_domains = ["doc.rust-lang.org"]
    start_urls = ["http://doc.rust-lang.org/book/"]
    local_site = "rust-book_response.html"
    id = 0

    def start_request(self):
        urls = [] + self.start_urls

        for i, url in enumerate(urls):
            yield scrapy.Request(url=url, callback=self.parse, cb_kwargs={})

    def parse(self, response):
        chapters = response.xpath(
            '//ol[@class="chapter"]//li[@class="chapter-item expanded " or @class="chapter-item expanded affix "]/a'
        )
        hrefs = chapters.xpath("@href").getall()
        texts = chapters.xpath("text()").getall()
        # Number of hrefs should be equal to number of texts
        assert len(hrefs) == len(texts)

        # Yield scrapy request for each chapter
        for href, chapter in zip(hrefs, texts):
            url = response.urljoin(href)
            yield scrapy.Request(
                url=url,
                callback=self.parse_chapter,
                cb_kwargs={"chapter": chapter},
            )

    def parse_chapter(self, response, chapter="Unknown"):
        # Scrawl chapter content
        content = response.xpath("//main/*")
        keywords = content.xpath("text()").getall()
        content = content.getall()

        self.id = self.id + 1
        yield {
            "id": str(self.id),
            "title": content[0],
            "content": "".join(p for p in content),
            "keywords": "".join(p for p in keywords),
            "url": response.url,
            "data": datetime.date.today().strftime("%Y-%m-%d"),
        }
