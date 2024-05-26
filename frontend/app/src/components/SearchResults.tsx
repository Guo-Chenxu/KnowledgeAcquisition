import { renderToString } from "react-dom/server";
import "./SearchResults.css";
import React, { useEffect, useState } from "react";
import { DocumentAbstract, SearchResult } from "../utils/types";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark as codeTheme } from "react-syntax-highlighter/dist/esm/styles/prism";
import { backend_url } from "../utils/config";
import {
  Table,
  TableHead,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Paper,
} from "@mui/material";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import { setSyntheticTrailingComments } from "typescript";

const darkTheme = createTheme({
  palette: {
    mode: "dark",
  },
});

interface SearchResultProps {
  result: SearchResult;
  onFeedback: (resultId: string, score: number) => void;
}

const SearchResultItem: React.FC<SearchResultProps> = ({
  result,
  onFeedback,
}) => {
  const [content, setContent] = useState("");
  const [abstract, setAbstract] = useState(<div>Loading abstract...</div>);
  const [contentFetched, setContentFetched] = useState<boolean>(false);
  const [abstractFetched, setAbstractFetched] = useState<boolean>(false);
  const [isFolded, setIsFolded] = useState(true); // new state
  const toggleFold = () => {
    setIsFolded(!isFolded);
  }; // function to toggle the fold
  const [entities, setEntities] = useState<{ [entity: string]: number }>({});
  const [hotWords, setHotWords] = useState<{ [entity: string]: number }>({});
  // Regex info extract
  const [regex, setRegex] = useState<string>("");
  const [regexResult, setRegexResult] = useState<string | null>(null);

  useEffect(() => {
    // Abstract fetch
    if (!abstractFetched) {
      fetch(`${backend_url}/extract_info?id=${result.id}`)
        .then((response) => response.json())
        .then((data) => {
          console.debug("data: ", data);
          let parsedData: DocumentAbstract = {
            entities: data.entities,
            hot_words: data.hot_words,
          };
          // Set states
          setEntities(parsedData.entities);
          setHotWords(parsedData.hot_words);
          setAbstractFetched(true);
        })
        .catch((error) => {
          console.error("Error loading abstract: ", error);
          setAbstract(<div>{"Failed: " + error}</div>);
        });
    }

    // Content fetch
    if (!isFolded && !contentFetched) {
      fetch(`${backend_url}/document?id=${result.id}`)
        .then((response) => response.json())
        .then((data) => {
          let doc = renderCodeBlocks(data.content);
          setContent(doc.body.innerHTML);
          setContentFetched(true);
        });
    }

    // Update abstract
    setAbstract(buildAbstract());
  }, [isFolded, content, contentFetched, result, entities, hotWords]);

  const handleRegexSubmit = () => {
    const matched = content.match(new RegExp(regex));
    console.log(matched);
    setRegexResult(matched ? matched[0] : "No matches found");
  };

  const generateItemFeedbackHandler = (
    items: { [s: string]: number },
    setItems: React.Dispatch<React.SetStateAction<{ [s: string]: number }>>,
    api: string
  ) => {
    return (resultId: string, item: string, score: number) => {
      fetch(api, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ resultId, item, score }),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log(data);
        })
        .catch((error) => {
          console.error(error);
        });

      let newItems = { ...items };
      let index = Object.entries(newItems).findIndex(
        (entry) => entry[0] == item
      );
      if (score === 1) {
        // Swap the item with the one before it to move it top
        newItems[item] += 1;
      } else if (score === 0) {
        if (index == newItems.length - 1 || newItems[item] === 1) {
          delete newItems[item];
        } else {
          newItems[item] -= 1;
        }
      }
      setItems(newItems);
    };
  };

  const handleEntityFeedback = generateItemFeedbackHandler(
    entities,
    setEntities,
    `${backend_url}/entity_feedback`
  );

  const handleHotWordFeedback = generateItemFeedbackHandler(
    hotWords,
    setHotWords,
    `${backend_url}/hotword_feedback`
  );

  const buildAbstract = () => {
    const buildTable = (
      title: string,
      name: string,
      entries: { [s: string]: number },
      feedbackHandler: (resultId: string, key: string, score: number) => void
    ) => {
      return (
        <ThemeProvider theme={darkTheme}>
          <TableContainer component={Paper}>
            <h2 style={{ padding: "10px" }}>{title}</h2>
            <Table sx={{ tableLayout: "fixed" }}>
              <TableHead>
                <TableRow>
                  <TableCell className="maxWidthCell">{name}</TableCell>
                  <TableCell className="maxWidthCell">Frequency</TableCell>
                  <TableCell className="maxWidthCell">Feedback</TableCell>
                </TableRow>
              </TableHead>

              <TableBody>
                {Object.entries(entries).map(([key, value]) => {
                  console.debug(key, value);
                  return (
                    <TableRow key={key}>
                      <TableCell className="maxWidthCell">
                        <div className="maxWidthCell">{key}</div>
                      </TableCell>
                      <TableCell className="maxWidthCell">
                        <div className="maxWidthCell">{value}</div>
                      </TableCell>
                      <TableCell>
                        <button
                          className="feedback"
                          onClick={() => feedbackHandler(result.id, key, 1)}
                        >
                          👍
                        </button>
                        <button
                          className="feedback"
                          onClick={() => feedbackHandler(result.id, key, 0)}
                        >
                          👎
                        </button>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </TableContainer>
        </ThemeProvider>
      );
    };
    const entitiesTable = buildTable(
      "Entities",
      "Entity",
      entities,
      handleEntityFeedback
    );
    const hotWordsTable = buildTable(
      "HotWords",
      "Hot word",
      hotWords,
      handleHotWordFeedback
    );

    return (
      <div>
        {entitiesTable}
        {hotWordsTable}
      </div>
    );
  };

  const renderCodeBlocks = (content: string) => {
    let parser = new DOMParser();
    let doc = parser.parseFromString(content, "text/html");
    let codeBlocks = doc.querySelectorAll("pre");

    codeBlocks.forEach((block) => {
      let codeString = block.textContent || "";
      let highlightedCodeString = renderToString(
        <SyntaxHighlighter language="cpp" style={codeTheme}>
          {codeString}
        </SyntaxHighlighter>
      );
      block.innerHTML = highlightedCodeString;
    });

    return doc;
  };

  const handleFeedback = (resultId: string, score: number) => {
    fetch(`${backend_url}/feedback`, {
      method: "POST",
      headers: {
        "Content-TYpe": "application/json",
      },
      body: JSON.stringify({ resultId, score }),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
        onFeedback(resultId, score);
      })
      .catch((error) => console.error("Error: ", error));
  };

  return (
    <div>
      <h2 onClick={toggleFold}>
        <div dangerouslySetInnerHTML={{ __html: result.title }} />
      </h2>
      {!isFolded && <div dangerouslySetInnerHTML={{ __html: content }} />}
      {abstract}
      <p>{result.url}</p>
      <p>{result.date}</p>
      <p>Relevance: {result.score}</p>

      {/* Regex info extraction */}
      <div>
        <input
          type="text"
          value={regex}
          onChange={(e) => setRegex(e.target.value)}
        />
        <button className="regex" onClick={handleRegexSubmit}>
          Submit Regex
        </button>
        {regexResult && (
          <p>
            Regex result:{" "}
            <div dangerouslySetInnerHTML={{ __html: regexResult }} />
          </p>
        )}
      </div>
      <button className="feedback" onClick={() => handleFeedback(result.id, 1)}>
        👍
      </button>
      <button className="feedback" onClick={() => handleFeedback(result.id, 0)}>
        👎
      </button>
    </div>
  );
};

interface SearchResultsProps {
  results: SearchResult[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ results }) => {
  const [adjustResults, setAdjustResults] = useState(results);

  useEffect(() => {
    setAdjustResults(results);
  }, [results]);

  const handleFeedback = (resultId: string, score: number) => {
    let newResults = [...adjustResults];

    let index = newResults.findIndex((result) => result.id === resultId);
    if (score === 1) {
      if (index > 0) {
        // Swap the item with the one before it to move it up
        [newResults[index], newResults[index - 1]] = [
          newResults[index - 1],
          newResults[index],
        ];
      }
    } else if (score === 0 && index < newResults.length - 1) {
      // Swap the item with the one after it to move it down
      [newResults[index], newResults[index + 1]] = [
        newResults[index + 1],
        newResults[index],
      ];
    }

    setAdjustResults(newResults);
  };

  console.debug(`Generate SearchResults for`);
  console.debug(results);

  return (
    <div>
      {adjustResults.map((result, index) => (
        <SearchResultItem
          key={`${result.id}-${index}`}
          result={result}
          onFeedback={handleFeedback}
        />
      ))}
    </div>
  );
};

export default SearchResults;
