document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("form");

    form.addEventListener("submit", async (event) => {
        event.preventDefault();

        const username = document.getElementById("email").value;
        const password = document.getElementById("password").value;

        const payload = {
            email: username,
            password: password 
        };

        try {
            const response = await fetch("/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(payload),
                redirect: "follow"
            });


            if (!response.ok) {
                // to understand golang responses
                const txt = await response.text()
                console.log("not ok", txt)
            }

            // this is a bit ugly
            if (response.redirected){
                // root
                window.location.href = "/"
            }

        } catch (error) {
            console.error("Login error:", error);
            alert("An unexpected error occurred. Please try again.");
        }
    });
});
