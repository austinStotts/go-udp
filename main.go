<!DOCTYPE html>
<html>
<head>
    <title>Go WebSocket Example</title>
</head>
<body>
    <h1>Go WebSocket Example</h1>
    <input type="text" id="inputData" placeholder="Enter data...">
    <button onclick="sendData()">Send Data</button>
    <div id="response"></div>

    <script>
        var ws = new WebSocket("ws://localhost:8080/ws");

        ws.onopen = function(event) {
            console.log("WebSocket opened.");
        };

        ws.onmessage = function(event) {
            // Display the received message
            document.getElementById("response").innerText = event.data;
        };

        function sendData() {
            var inputData = document.getElementById("inputData").value;

            // Send data to the server
            ws.send(inputData);
        }
    </script>
</body>
</html>
