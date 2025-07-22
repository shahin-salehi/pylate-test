async function search() {
    const query = document.getElementById('searchInput').value;
    const category = document.getElementById('categorySelect').value;
    const resultsContainer = document.getElementById('results');

    // Clear previous results
    resultsContainer.innerHTML = `<div class="text-gray-500 italic">Searching for "${query}"...</div>`;

    try {
        console.log('category', category)
        const response = await fetch('http://localhost:8080/query', {
            method: 'GET',
            headers: {
                'query': query,
                'category': category
            }
        });

        if (!response.ok) {
            throw new Error(`Server responded with ${response.status}`);
        }

        const html = await response.text();

        // Replace container contents with response HTML
        resultsContainer.innerHTML = html;

    } catch (err) {
        console.error('Search failed:', err);
        resultsContainer.innerHTML = `<div class="text-red-600 font-medium">Error: ${err.message}</div>`;
    }
}
