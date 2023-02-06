##
## Santaclaus - Indexer
##

GO	=	go
PROTOC	=	protoc

NAME	=	santaclaus

SRCDIR	=	src

SRC		=	main.go

SRC			:= $(addprefix $(SRCDIR)/, $(SRC))

GOFLAGS =	--trimpath --mod=vendor

all: $(NAME)

$(NAME):
	$(GO) mod vendor
	$(GO) build $(GOFLAGS) -o $(NAME) $(SRC)

fclean:
	rm -f  $(NAME)

re: fclean all

PROTODIR	=	./third_parties/protobuf-interfaces

PROTOSRCDIR	=	$(PROTODIR)/src

PROTOOUTDIR	=	$(PROTODIR)/generated

PROTOSRC	=	common/File.proto \
				Maestro_Santaclaus/Maestro_Santaclaus.proto

PROTOSRC	:= $(addprefix $(PROTOSRCDIR)/, $(PROTOSRC))


gRPC:
	$(PROTOC) \
	--go_out=$(PROTOOUTDIR) \
	--go-grpc_out=$(PROTOOUTDIR) \
	-I $(PROTOSRCDIR) \
	$(PROTOSRC)

#	--go_opt=paths=source_relative \
#	--go-grpc_opt=paths=source_relative \