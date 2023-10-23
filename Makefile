##
## Santaclaus - Indexer
##

GO		=	go
PROTOC	=	protoc

NAME	=	santaclaus

SRCDIR	=	src

SRC		=	main.go \
			SantaclausServerCreateDir.go \
			SantaclausServerFindDir.go \
			SantaclausServerGetDirectory.go \
			SantaclausServer.go \
			SantaclausServerInit.go \
			SantaclausServerRemoveDirectory.go \
			SantaclausServerMoveFile.go \
			SantaclausServerStructs.go \
			SantaclausServerUtilsDisks.go


SRC		:= $(addprefix $(SRCDIR)/, $(SRC))

GOFLAGS =	--trimpath --mod=vendor

all: $(NAME)

fclean:
	rm -f  $(NAME)
	rm -f  $(UNIT_TESTS)

$(NAME):	fclean
	$(GO) mod vendor
	$(GO) build $(GOFLAGS) -o $(NAME) $(SRC)


# PROTOBUF - GRPC

PROTODIR	=	./third_parties/protobuf-interfaces

PROTOSRCDIR	=	$(PROTODIR)/src

PROTOOUTDIR	=	$(PROTODIR)/generated

PROTOSRC	=	common/File.proto \
				common/Disk.proto \
				Maestro_Santaclaus/Maestro_Santaclaus.proto \
				Santaclaus_HardwareMalin/Santaclaus_HardwareMalin.proto

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

unit_tests: fclean
	$(GO) test -c -o $(UNIT_TESTS) ./$(SRCDIR)
