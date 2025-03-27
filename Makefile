.PHONY: install run dev build start
 
install-backend:
	cd backend && go mod tidy && go install

run-backend:
	cd backend && go run cmd/main.go

install-frontend:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-start:
	cd frontend && npm run start