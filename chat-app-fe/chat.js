const chatForm = document.querySelector(".chat-form");
const messageInput = document.getElementById("message-input");
const receiverId = document.getElementById("receiver-input");
const token = localStorage.getItem("token");
const user = JSON.parse(localStorage.getItem("user"));

const ws = new WebSocket(`ws://localhost:8080/chat/${token}`);

ws.addEventListener("open", () => {
  console.log("Connected to WebSocket");
});

ws.addEventListener("close", () => {
  console.log("Disconnected from WebSocket");
});

ws.addEventListener("message", (event) => {
  console.log(`Message received: ${event.data}`);
  const message = JSON.parse(event.data);

  addMessage(message);
});

chatForm.addEventListener("submit", (event) => {
  event.preventDefault();

  const content = messageInput.value;
  messageInput.value = "";

  const message = {
    type: "message",
    data: {
      content: content,
      sender_id: user.id,
      sender_username: user.username,
      recipient_id: +receiverId.value,
      type: 0,
    },
  };

  addMessage(message, "You");

  console.log(`Message sent: ${JSON.stringify(message)}`);
  ws.send(JSON.stringify(message));
});

// Get the chat box and form elements
const chatBox = document.querySelector(".chat-box");

// Function to add a message to the chat box
function addMessage(message, sender = "") {
  const messageElement = document.createElement("div");

  console.log(message.data.content);
  messageElement.className = "chat-message";
  messageElement.innerHTML = `<span>${
    !sender ? message.data.sender_username : sender
  }:</span> ${message.data.content}`;

  chatBox.appendChild(messageElement);
  chatBox.scrollTop = chatBox.scrollHeight;
}
