# Gitmon

Gitmon is a command-line tool for monitoring Git repositories and automatically
committing and pushing changes. It helps you stay updated on your repositories
without manual intervention.

## Features

- **Automatic Monitoring**: Continuously monitors specified Git repositories
  for changes.
- **Commit and Push**: Automatically commits changes and pushes them to
  the repository.
- **Graceful Shutdown**: Handles interruptions gracefully and ensures no data loss.

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/clevrf0x/gitmon.git
   cd gitmon
   ```

2. **Build the executable**:

   ```bash
   make build
   ```

## Usage

```bash
./bin/gitmon [repo1] [repo2] ... [repoN]
```

- Replace `[repo1]`, `[repo2]`, etc., with the paths to the Git repositories
  you want to monitor.

## Commands

- **`make all`**: Runs tests and builds the executable.
- **`make test`**: Executes all unit tests.
- **`make build`**: Compiles the application into the `bin/gitmon` executable.
- **`make run`**: Runs the compiled `gitmon` executable with specified arguments.
- **`make clean`**: Removes the `bin/` directory and its contents.

## Example

Monitor changes in two Git repositories:

```bash
./bin/gitmon /path/to/repo1 /path/to/repo2
```

## Dependencies

- `github.com/clevrf0x/gitmon/internal/fslock`: File system locking mechanism
  for ensuring single instance execution.
- Standard Go packages (`os`, `log`, `flag`, `context`, `sync`,
  `os/exec`, `os/signal`, `path/filepath`, `syscall`).

### Contributing

- Fork the repository, make your changes, and submit a pull request.
- Report issues or suggest improvements using the issue tracker.

### License

This project is licensed under the GNU GPLv3 License -
see the [LICENSE](LICENSE) file for details.

---
