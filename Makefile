.PHONY: run

run:
	@echo "🔄 Запуск Kafka..."
	@osascript -e 'tell application "Terminal" to do script "/opt/homebrew/opt/kafka/bin/kafka-server-start /opt/homebrew/etc/kafka/server.properties"'
	sleep 5
	@echo "💬 Запуск Telegram-бота (producer)..."
	@osascript -e 'tell application "Terminal" to do script "cd $(PWD) && go run cmd/main.go"'
	sleep 2
	@echo "🧠 Запуск consumer..."
	@osascript -e 'tell application "Terminal" to do script "cd $(PWD) && go run cmd/consumer/main.go"'