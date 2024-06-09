<template>
  <v-container>
    <!--  class="text-center" -->
    <v-row>
      <v-col cols="12">
        <v-img
          :src="require('../assets/logo.svg')"
          class="my-3"
          contain
          height="200"
        />
      </v-col>

      <v-col cols="12">
        <h1 class="display-2 font-weight-bold text-center">
          Knowledge Acquisition
        </h1>
      </v-col>

      <v-col cols="12" class="d-flex align-center">
        <h2 class="headline font-weight-bold mb-1 mr-1">关键词搜索</h2>
        <v-text-field
          v-model="searchText"
          label="输入关键字"
          outlined
          dense
          class="flex-grow-1"
        ></v-text-field>
        <v-btn @click="searchByKeyword" color="primary" class="mb-1 ml-1"
          ><span class="white--text font-weight-black">搜索</span></v-btn
        >
      </v-col>

      <v-col cols="12" class="d-flex align-center">
        <h2 class="headline font-weight-bold mb-1">图片检索</h2>
        <v-file-input
          v-model="imageFile"
          label="点击上传图片"
          outlined
          dense
          class="flex-grow-1"
          @change="onFileChange"
        ></v-file-input>
        <v-btn @click="searchByImage" color="primary" class="mb-1 ml-1"
          ><span class="white--text font-weight-black">搜索</span></v-btn
        >
      </v-col>

      <!-- Search Results Section -->
      <v-col cols="12">
        <v-row>
          <v-card
            width="100%"
            class="mb-1"
            v-for="(result, index) in searchResults"
            :key="index"
          >
            <v-card-title
              @click="toggleDetail(result.Doc.id)"
              v-html="result.Doc.title"
            ></v-card-title>
            <v-card-subtitle>{{ result.Doc.date }}</v-card-subtitle>
            <v-card-text
              ><div>介绍：{{ result.Doc.content }}</div>
              <div>相关度：{{ result.Score }}</div>
              链接：
              <a :href="result.Doc.url" target="_blank">
                {{ result.Doc.url }}
              </a>
              <div
                v-if="
                  detailMap[result.Doc.id] && detailMap[result.Doc.id].visible
                "
              >
                <div v-html="detailMap[result.Doc.id].content"></div>
                <div>关键字: {{ getShortKeywords(result.Doc.id) }}</div>
                <div>
                  语言:
                  {{ detailMap[result.Doc.id].Lang === 0 ? "英文" : "中文" }}
                </div>
                <!-- Entity Table -->
                <v-simple-table dense>
                  <template v-slot:default>
                    <thead>
                      <tr>
                        <th class="text-left">实体</th>
                        <th class="text-left">频率</th>
                        <th class="text-left">反馈</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(value, key) in detailMap[result.Doc.id]
                          .entities"
                        :key="key"
                      >
                        <td>{{ key }}</td>
                        <td>{{ value }}</td>
                        <td>
                          <v-rating
                            dense
                            hover
                            small
                            v-model="
                              detailMap[result.Doc.id].entities[key].score
                            "
                            @input="
                              handleEntityFeedback(
                                result.Doc.id,
                                key,
                                detailMap[result.Doc.id].entities[key].score
                              )
                            "
                          ></v-rating>
                        </td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
                <!-- Hot Words Table -->
                <v-simple-table dense>
                  <template v-slot:default>
                    <thead>
                      <tr>
                        <th class="text-left">热词</th>
                        <th class="text-left">频率</th>
                        <th class="text-left">反馈</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr
                        v-for="(value, key) in detailMap[result.Doc.id]
                          .hot_words"
                        :key="key"
                      >
                        <td>{{ key }}</td>
                        <td>{{ value }}</td>
                        <td>
                          <v-rating
                            dense
                            hover
                            small
                            v-model="
                              detailMap[result.Doc.id].hot_words[key].score
                            "
                            @input="
                              handleHotwordFeedback(
                                result.Doc.id,
                                key,
                                detailMap[result.Doc.id].hot_words[key].score
                              )
                            "
                          ></v-rating>
                        </td>
                      </tr>
                    </tbody>
                  </template>
                </v-simple-table>
              </div>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-rating
                dense
                hover
                v-model="result.Score"
                @input="handleOverallFeedback(result.Doc.id, result.Score)"
              ></v-rating>
            </v-card-actions>
          </v-card>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import axios from "axios";

export default {
  name: "MainPage",

  data: () => ({
    searchText: "", // 用于存储关键字搜索文本
    imageFile: null, // 用于存储要上传的图片文件
    searchResults: [],
    detailMap: {},
  }),

  methods: {
    getShortKeywords(id) {
      const keywords = this.detailMap[id] ? this.detailMap[id].keywords : "";
      return keywords.length > 50 ? keywords.slice(0, 50) + "..." : keywords;
    },
    searchByKeyword() {
      const params = { q: this.searchText, page: 1, limit: 10 };
      axios
        .get("/api/v1/search", { params })
        .then((response) => {
          console.log(response.data);
          this.searchResults = response.data;
        })
        .catch((error) => {
          console.error("Error during keyword search:", error);
        });
    },

    onFileChange(file) {
      this.imageFile = file;
    },

    searchByImage() {
      if (!this.imageFile) {
        alert("Please upload an image.");
        return;
      }
      const formData = new FormData();
      formData.append("image", this.imageFile);
      axios
        .post("api/v1/search_by_image", formData, {
          headers: { "Content-Type": "multipart/form-data" },
        })
        .then((response) => {
          this.searchResults = response.data.results;
          this.searchText = response.data.keywords;
        })
        .catch((error) => {
          console.error("Error during image search:", error);
        });
    },

    toggleDetail(id) {
      if (this.detailMap[id] && this.detailMap[id].visible) {
        this.$set(this.detailMap[id], "visible", false);
      } else if (this.detailMap[id] && !this.detailMap[id].visible) {
        this.$set(this.detailMap[id], "visible", true);
      } else {
        axios
          .all([
            axios.get(`api/v1/document`, { params: { id } }),
            axios.get(`api/v1/extract_info`, { params: { id } }),
          ])
          .then(
            axios.spread((DocRes, infoRes) => {
              const entitiesWithScore = Object.entries(
                infoRes.data.entities
              ).reduce((acc, [key, value]) => {
                acc[key] = { value, score: 0 }; // 初始化每个实体的评分为0
                return acc;
              }, {});

              const hotWordsWithScore = Object.entries(
                infoRes.data.hot_words
              ).reduce((acc, [key, value]) => {
                acc[key] = { value, score: 0 }; // 初始化每个热词的评分为0
                return acc;
              }, {});
              this.$set(this.detailMap, id, {
                visible: true,
                content: DocRes.data.content,
                keywords: DocRes.data.keywords,
                Lang: DocRes.data.Lang,
                entities: infoRes.data.entities,
                hot_words: infoRes.data.hot_words,
              });
            })
          )
          .catch((error) => {
            console.error("Error fetching Document details:", error);
          });
      }
    },
    handleEntityFeedback(resultId, item, score) {
      const payload = {
        item,
        resultId,
        score,
      };
      axios
        .post("api/v1/entity_feedback", payload)
        .then((response) => {
          console.log("Entity Feedback sent successfully", response);
        })
        .catch((error) => {
          console.error("Error sending entity feedback", error);
        });
    },

    handleHotwordFeedback(resultId, item, score) {
      const payload = {
        item,
        resultId,
        score,
      };
      axios
        .post("api/v1/hotword_feedback", payload)
        .then((response) => {
          console.log("Hotword Feedback sent successfully", response);
        })
        .catch((error) => {
          console.error("Error sending hotword feedback", error);
        });
    },

    handleOverallFeedback(resultId, score) {
      const payload = {
        resultId,
        Score: score,
      };
      axios
        .post("api/v1/feedback", payload)
        .then((response) => {
          console.log("Overall Feedback sent successfully", response);
        })
        .catch((error) => {
          console.error("Error sending overall feedback", error);
        });
    },
  },
};
</script>


<style scoped>
.mb-1 {
  margin-bottom: 26px !important;
}
.mr-1 {
  margin-right: 10px !important;
}
.ml-1 {
  margin-left: 10px !important;
}
.mt-1 {
  margin-top: 10px !important;
}
</style>
