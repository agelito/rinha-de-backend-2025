benchmark:
	hey -m "POST" -n 1000 -d '{"correlationId": "be463535-0661-413e-b592-a24759b72fc0", "amount": 10.22 }' -H "Content-Type: application/json" http://localhost:3001/payments

proto:
	protoc --go_out=./messages/model ./messages/proto/payments.proto
	protoc --go_out=./messages/model ./messages/proto/servers.proto

stop-containers:
	docker compose -f payment-processors/docker-compose.yml down
	docker compose down

start-payment-processors:
	docker compose -f payment-processors/docker-compose.yml up -d

start-backend:
	docker compose up

run-k6-tests:
	k6 run tests/rinha.js

test: stop-containers start-payment-processors start-backend run-k6-tests
