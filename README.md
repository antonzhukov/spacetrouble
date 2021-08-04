# spacetrouble
This is a coding challenge implementation based on the [requirements](REQUIREMENTS.md)

The solution implements 2 endpoints: add booking, get all bookings.

## Install and run

    make build
    docker-compose up --build -d

## Test

    make test

## API
The server exposes /bookings endpoint:
- POST method to create a booking
- GET method to get all bookings

### Examples
Successful add booking request

    curl -X POST localhost:8000/bookings  -d '{"first_name":"Anton", "last_name":"Zhukov", "birthday":"2021-01-01T00:00:00+00:00", "launchpad_id":"5e9e4502f509092b78566f87", "destination_id": "abc", "launch_date": "2031-01-01T00:00:00+00:00"}' -i

Add booking requests which fails due to flight being cancelled

    curl -X POST localhost:8000/bookings  -d '{"first_name":"Anton", "last_name":"Zhukov", "birthday":"2021-01-01T00:00:00+00:00", "launchpad_id":"5e9e4502f509092b78566f87", "destination_id": "abc", "launch_date": "2022-01-01T00:00:00+00:00"}' -i

Get all added bookings

    curl localhost:8000/bookings -i