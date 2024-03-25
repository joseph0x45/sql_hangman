setup-db:
	@docker stop hangman 
	@docker rm hangman
	@docker run --name=hangman -e POSTGRES_PASSWORD=pwd -e POSTGRES_DB=game -itd -p 5432:5432 postgres:latest

migrate:
	@docker docker cp ./schema.sql hangman:/tmp/schema.sql
	@docker exec -it hangman psql -U postgres -d game -f /tmp/schema.sql
