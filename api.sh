#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <up|down>"
    exit 1
fi

if [ "$1" == "up" ]; then
    # Step 1: Move to directory
    cd ./api-1.4

    # Step 2: Run app.js in detached mode
    node app.js &

    # Wait for a few seconds to ensure the API server is up and running
    sleep 5

    # Step 3: Send HTTP POST request
    curl -X POST -H "Content-Type: application/json" -d '{"username":"boss","orgName":"Org1"}' http://localhost:4000/users

    echo "Server started."
elif [ "$1" == "down" ]; then
    # Find and kill the node process running app.js
    pkill -f "node app.js"

    echo "Server stopped."
else
    echo "Invalid argument. Usage: $0 <up|down>"
    exit 1
fi

echo "Script execution completed."
