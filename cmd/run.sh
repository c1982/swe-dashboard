export $(grep -v '^#' .env | xargs)
go run .
