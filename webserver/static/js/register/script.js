
async function register() {
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    console.log(username, email, password)

    try {
        const response = await fetch('http://localhost:8080/register-user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                "username": username,
                "email": email,
                "password_hash": password, //not yet hashed but im singy with types
            })
        });


        if (!response.ok) {
            throw new Error(`Server responded with ${response.status}`);
        }


        // Replace container contents with response HTML
        // write feedback here
        //resultsContainer.innerHTML = html;

    } catch (err) {
        // Replace container contents with response HTML
        console.error('Search failed:', err);
        //resultsContainer.innerHTML = `<div class="text-red-600 font-medium">Error: ${err.message}</div>`;
    }
}
document.addEventListener('DOMContentLoaded', () => {

    const registerForm = document.getElementById("register-form");


    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        register();
    });
});
