(() => {
    let searchBox = document.getElementById("search-box");
    let searchResults = document.getElementById("search-results");
    let toc = document.getElementById("toc");

    /**
     * @param {string} results
     */
    function updateSearchResults(results) {
        if (results.length > 0) {
            searchResults.innerHTML = "";
            for (let i = 0; i < results.length; i++) {
                let resultElem = document.createElement("div");
                let resultAnchor = document.createElement("a");
                resultAnchor.href = results[i]; // we should be at /, so this is right
                resultAnchor.innerText = results[i];
                resultElem.appendChild(resultAnchor);
                searchResults.appendChild(resultElem);
            }
            toc.classList.add("hidden");
            searchResults.classList.remove("hidden");
        } else {
            toc.classList.remove("hidden");
            searchResults.classList.add("hidden");
        }
    }

    async function handleSearchInput() {
        // get search query
        const query = searchBox.value;
        if (query == "") {
            updateSearchResults([]);
            return
        }
        // make request
        const response = await fetch("/search/" + query);
        if (!response.ok) {
            throw new Error("Couldn't search, status code ", response.status);
        }
        let result = await response.text();
        updateSearchResults(result.split('\n'));
    }

    searchBox.addEventListener("input", handleSearchInput);
})();
