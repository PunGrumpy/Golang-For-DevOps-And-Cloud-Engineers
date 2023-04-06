# test-server

The test server is used in multiple lectures. You can use the following command to start the test-server:
```bash
./start-test-server.sh
```

## Docker

The test server can also be started using docker. The following command will start the test server on port 8080:
```bash
docker build -t test-server .
docker run -p 8080:8080 test-server
```

## Docker Compose

The test server can also be started using docker-compose. The following command will start the test server on port 8080:
```bash
docker-compose up -d
# or
docker compose up -d
```

# Notes
If you're using zsh, make sure to use quotes around the URL when testing.
