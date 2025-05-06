async function fetchMessages() {
    try {
        const response = await fetch("/live");
        const data = await response.json();

        const list = document.getElementById("messageList");
        list.innerHTML = "";

        data.forEach(msg => {
            const li = document.createElement("li");
            li.textContent = msg;
            list.appendChild(li);
        });
    } catch (error) {
        console.error("Error fetching messages:", error);
    }
}

// Fetch messages immediately and then every 3 seconds
fetchMessages();
setInterval(fetchMessages, 3000);

