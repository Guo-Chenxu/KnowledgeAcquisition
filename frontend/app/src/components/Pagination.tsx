import React from "react";
import "./Pagination.css";

interface PaginationProps {
  currentPage: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  onPageChange,
}) => {
  return (
    <div className="Pagination">
      <button
        className="Pagination-button"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
      >
        Previous
      </button>
      <span className="Pagination-currentPage">{currentPage}</span>
      <button
        className="Pagination-button"
        onClick={() => onPageChange(currentPage + 1)}
      >
        Next
      </button>
    </div>
  );
};

export default Pagination;
