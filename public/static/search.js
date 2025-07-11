// Copyright (c) 2024 Julian Müller (ChaoticByte)
(() => {
    let searchBox = document.getElementById("search-box");
    let searchResults = document.getElementById("search-results");
    let toc = document.getElementById("toc");

    /**
     * @param {string} results
     */
    function showSearchResults(results) {
        searchResults.innerHTML = "";
        if (results.length > 0) {
            results.forEach(r => {
                let resultElem = document.createElement("div");
                let resultAnchor = document.createElement("a");
                resultAnchor.href = r[0]; // we should be at /, so this is right
                resultAnchor.innerText = r[1];
                resultElem.appendChild(resultAnchor);
                searchResults.appendChild(resultElem);
            });
        }
        toc.classList.add("hidden");
        searchResults.classList.remove("hidden");
    }
    function hideSearchResults() {
        toc.classList.remove("hidden");
        searchResults.classList.add("hidden");
    }

    async function handleSearchInput() {
        // get search query
        const query = searchBox.value;
        if (query == "") {
            hideSearchResults();
            return
        }
        // make request
        const response = await fetch("/search/" + query);
        if (!response.ok) {
            throw new Error(`Search API returned status code ${response.status}`);
        }
        let results_raw = await response.text();
        let results = [];
        if (results_raw.length > 0) {
            results_raw.split('\n').forEach(r => {
                results.push(r.split("|", 2));
            });
        }
        showSearchResults(results);
    }

    searchBox.addEventListener("input", handleSearchInput);
})();
