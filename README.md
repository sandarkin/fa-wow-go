# Word of Wisdom TCP Server

# Task

Design and implement “Word of Wisdom” tcp server.

- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

# How to run

## With docker-compose

```sh
# up service
docker-compose up -d
# view logs
docker logs fa-wow-go-client-1 -f
# and
docker logs fa-wow-go-server-1 -f
```

## With Makefile

```sh
make run_server
make run_client
```

# Why sha256 for PoW algorithm?

It uses the widely accepted hash function sha256, which simplifies the implementation of the task because it is included in the standard libraries. 
But at the same time the calculation difficulty is sufficient to protect against DDoS attacks.