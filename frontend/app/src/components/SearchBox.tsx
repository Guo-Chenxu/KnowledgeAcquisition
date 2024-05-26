import "./SearchBox.css";
import { useState } from "react";
import React from "react";

interface SearchBoxProps {
  onSearch: (query: string, page: number) => void;
  onImageSearch: (image: File | null) => void;
  isSearching: boolean;
  setIsSearching: (isSearching: boolean) => void;
  query: string;
  setQuery: (query: string) => void;
}

const SearchBox: React.FC<SearchBoxProps> = ({
  onSearch,
  onImageSearch,
  isSearching,
  setIsSearching,
  query,
  setQuery,
}) => {
  const [file, setFile] = useState<File | null>(null);

  const handleTextSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSearching(true);
    onSearch(query, 1); // 1 as the initial page number
  };

  const handleImageSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSearching(true);
    if (file) {
      onImageSearch(file);
    } else {
      console.error("Error searching by image: file not given.");
    }
  };

  return (
    <div className="SearchBox">
      <form onSubmit={handleTextSearchSubmit}>
        {/* Search by text */}
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search..."
        />
        <button type="submit" disabled={isSearching}>
          Search
        </button>
      </form>

      <form onSubmit={handleImageSearchSubmit}>
        {/* Search by image */}
        <input
          type="file"
          accept="image/*"
          onChange={(e) => setFile(e.target.files?.item(0) || null)}
          id="fileUpload"
        />
        <button type="submit" disabled={isSearching}>
          Search<br></br>by Image
        </button>
      </form>
    </div>
  );
};

export default SearchBox;
