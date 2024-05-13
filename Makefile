createmigration:
	migrate create --ext=mysql -dir=internal/infra/database/migrations -seq

.PHONY: migrate