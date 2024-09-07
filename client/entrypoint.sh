#!/bin/sh

echo "Container is ready. To run the clients run docker exec -it client<NUMBER> /app/client"

# Keep the container running indefinitely
tail -f /dev/null
