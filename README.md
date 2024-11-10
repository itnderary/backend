# Itnderary back-end

`Itnderary backend` is a lightweight Go application providing backend APIs for the Itnderary application.

## Developing

This project uses `make` to build and run the application.

Supported commands:

- `make build` - builds the binary for your OS. You can then run it using `itnderary`.
- `make package` - builds a local Docker image.
- `make run` - runs the Docker image locally, making the service accessible via `http://localhost:3000`.
- `make publish` - publishes the image in the configured Docker repository

## Contributing

Please refer to our LICENSE and CONTRIBUTING.md for our guidelines.
