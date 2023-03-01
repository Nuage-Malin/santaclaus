##
## Santaclaus - Indexer
##

GO		=	go
PROTOC	=	protoc

NAME	=	santaclaus

SRCDIR	=	src

SRC		=	main.go \
			SantaclausServer.go

SRC		:= $(addprefix $(SRCDIR)/, $(SRC))

GOFLAGS =	--trimpath --mod=vendor

all: $(NAME)

$(NAME):
	$(GO) mod vendor
	$(GO) build $(GOFLAGS) -o $(NAME) $(SRC)

fclean:
	rm -f  $(NAME)

re: fclean all

# PROTOBUF - GRPC

PROTODIR	=	./third_parties/protobuf-interfaces

PROTOSRCDIR	=	$(PROTODIR)/src

PROTOOUTDIR	=	$(PROTODIR)/generated

PROTOSRC	=	common/File.proto \
				Maestro_Santaclaus/Maestro_Santaclaus.proto

PROTOSRC	:= $(addprefix $(PROTOSRCDIR)/, $(PROTOSRC))

PATH := $(PATH):$(shell go env GOPATH)/bin
export $(PATH)

gRPC:	
	$(PROTOC) \
	--go_out=$(PROTOOUTDIR) \
	--go-grpc_out=$(PROTOOUTDIR) \
	-I $(PROTOSRCDIR) \
	$(PROTOSRC)

UNIT_TESTS	=	unit_tests

unit_tests:
	$(GO) test -c -o $(UNIT_TESTS) ./$(SRCDIR)