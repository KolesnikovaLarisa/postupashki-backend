# Настройки
GO_BUILDER = go
GO_TESTER = go test
BUILD_BIN = server
TEST_PACKAGE = ./...
PYTEST_COMMAND = python -m pytest tests.py

# Правило для сборки
build:
	${GO_BUILDER} build -o ${BUILD_BIN} .

# Правило для запуска тестов
test:
	${PYTEST_COMMAND}

# Правило для запуска сервера
run:
	${GO_BUILDER} run . 

# Правило для очистки
clean:
	rm -f ${BUILD_BIN}

# Правило для запуска всех задач
all: 
	build test run