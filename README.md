# SimulationCraft REST API

This is a simple REST API for SimulationCraft.

## Configuration

Set the `PORT` environment variable to specify the listening port.

## Usage

Use `PUT/{name}` to create a new simulation. The body must be a valid simc profile.
Use `GET/{name}` to query an running or completed simulation.
Use `HEAD/{name}` to check whether a simulation is completed without downloading the actual results.

Running/pending simulations will result in a `202 Accepted` status code.
`GET` and `HEAD` will use the `Accept` header in order to switch between `text/html` (default), `application/json` and `text/plain` output.