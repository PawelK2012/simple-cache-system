### Simple cache server

[![Go](https://github.com/PawelK2012/simple-cache-system/actions/workflows/go.yml/badge.svg)](https://github.com/PawelK2012/simple-cache-system/actions/workflows/go.yml)

Small library for caching user information to increase databse throughput and availability. 

### Run it
 Simple exec test by typing  `make test` into terminal. This will exec test that is simulating hitting API 1000 times.

The goal is to retrive only 10% of traffic from DB and other 90% from cache. From test logs you can see this is a case and program hits DB only 100 times. 

To start server simply type `make run` and make GET request to `http://localhost:3000/user?id=9`

Exec code coverage with `make coverage` command