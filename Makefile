benchmark:
	hey -m "POST" -n 1000 -d '{"correlationId": "be463535-0661-413e-b592-a24759b72fc0", "amount": 10.22 }' -H "Content-Type: application/json" http://localhost:3001/payments

proto:
	protoc --go_out=./messages/model ./messages/proto/payments.proto
