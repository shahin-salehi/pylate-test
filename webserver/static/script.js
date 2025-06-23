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

document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById("uploadForm");
    const fileInput = document.getElementById("fileInput");

    uploadForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(uploadForm);
        formData.append("my-field", "some extra data for you");

        const numFiles = fileInput.files.length;
        console.log("number of files:", numFiles);
        console.log("first file:", fileInput.files[0]);
        console.log("formdata:", formData);

        try {
            const response = await fetch('/upload-pdf', {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                const message = await response.text();
                alert("✅ Upload successful");
                location.reload(); // Optionally reload to reflect updates
            } else {
                const errorText = await response.text();
                alert("❌ Upload failed: " + errorText);
            }
        } catch (error) {
            alert("❌ Network or server error: " + error.message);
        }
    });
});
