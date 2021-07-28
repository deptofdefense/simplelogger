# SimpleLogger

SimpleLogger is a golang package for logging.

Use SimpleLogger when a bare-bones logger is required for logging.

To use this module add `import "github.com/deptofdefense/simplelogger"` to your package.

## Example

To use the lock in AWS S3 you could use this code:

```golang
logger, errInitLogger := log.InitLogger(log.Stdout)
logger.Log("Message", map[string]interface{}{"metadata": metadata})
```

SimpleLogger also provides a primitive to build your own logger named `SimpleLogger`.

## Testing

Testing can be done using the `Makefile` targets `make test` and `make test_coverage`.

## Development

Development for this has been geared towards MacOS users. Install dependencies to get started:

```sh
brew install circleci go golangci-lint pre-commit shellcheck
```

Install the pre-commit hooks and run them before making pull requests:

```sh
pre-commit install
pre-commit run -a
```

## License

This project constitutes a work of the United States Government and is not subject to domestic copyright protection under 17 USC § 105.  However, because the project utilizes code licensed from contributors and other third parties, it therefore is licensed under the MIT License.  See LICENSE file for more information.
