document.addEventListener('DOMContentLoaded', () => {
    const uploadForm = document.getElementById("uploadForm");
    const fileInput = document.getElementById("fileInput");
    const tagInput = document.getElementById("tagInput");


    uploadForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const tag = tagInput.value.trim();
        const formData = new FormData(uploadForm);
        formData.append("my-field", tag); // trimmed

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

// delet file
document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".delete-form").forEach(form => {
        form.addEventListener("submit", async (e) => {
            e.preventDefault(); // prevent full page reload

            const fileId = parseInt(form.getAttribute("data-id"), 10);
            const res = await fetch("/delete", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ id: fileId })
            });

            if (res.ok) {
                // Remove the parent <li> from the DOM
                form.closest("li").remove();
            } else {
                alert("Failed to delete file.");
            }
        });
    });
});
