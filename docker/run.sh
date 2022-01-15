export $(grep -v '^#' .env | xargs)
docker-compose up
