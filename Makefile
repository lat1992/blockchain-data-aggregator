NAME			=	blockchain-data-aggregator

BUILD_NAME		=	$(NAME)

BUILD_DIR		=	./build

CMD_DIR			=	./cmd/aggregator

GO				=	go

GO_BUILD		=	$(GO) build

GO_CLEAN		=	$(GO) clean

GO_TEST			=	$(GO) test

RM				=	rm -f

MKDIR			=	mkdir -p

DOCKER			=	docker

all				:	build

build			:
					$(MKDIR) $(BUILD_DIR)
					$(GO_BUILD) -o $(BUILD_DIR)/$(BUILD_NAME) -v $(CMD_DIR)/main.go

test			:
					$(GO_TEST) -v ./...

vet 			:
					$(GO) vet ./...

clean			:
					$(GO_CLEAN)
					$(RM) $(BUILD_DIR)/$(BUILD_NAME)

docker-compose	:
					$(DOCKER) compose up -d

docker-build	:
	                $(DOCKER) build -t $(NAME) .

docker-run      :
	                $(DOCKER) run -d --name $(NAME) -p 8080:8080 $(NAME)

install			:	docker-build docker-compose

rebuild			:	clean build

re				:	rebuild

.PHONY			:	all make build test vet clean docker-compose docker-build docker-run install rebuild re
