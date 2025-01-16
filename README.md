cmd: go mod tidy
cmd: go run .\src\cmd\main.go
cmd: sqlc generate

Note:
Because we have new dynamic relationships in the database and each relationship with nested object response requires creating a Table View, so we cannot predict the number in advance, so SQLC will no longer be suitable for this project.
