
# zODC - Masterflow
This is a project

## Installation

### First step - Important preparation

Using [Golang 1.23.4](*https://go.dev/dl/) and follow instruction below :

Install lazy library :
```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

.Env File
```
Put the .env file in the root directory
```

### Final step - Setup require libs and packages

Prepare the necessary libraries
```bash
task setup
```

## Run server

In development :
```bash
task run-dev
```

# Note

### Before Commit

Remove unnecessary libs in the project
```bash
task clean
```
