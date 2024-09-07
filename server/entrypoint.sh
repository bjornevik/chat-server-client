#!/bin/sh

echo "Container is ready. Run the server with docker exec -it server /app/server"

# Keep the container running indefinitely
tail -f /dev/null
