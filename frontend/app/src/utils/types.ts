export interface SearchResult {
  id: string;
  title: string;
  content: string;
  url: string;
  date: string;
  score: number;
}

export interface DocumentAbstract {
  entities: { [entity: string]: number };
  hot_words: { [word: string]: number };
}
