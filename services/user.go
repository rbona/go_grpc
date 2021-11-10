package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/rbona/fc2-grpc/pb"
)

// type UserServiceClient interface {
// 	AddUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
// 	AddUserVerbose(ctx context.Context, in *User, opts ...grpc.CallOption) (UserService_AddUserVerboseClient, error)
// 	AddUsers(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUsersClient, error)
// }

type UserService struct {
	pb.UnimplementedUserServiceServer
}

var idAtual = 0

func NewUserService() *UserService {

	fmt.Println("Servidor iniciado.")
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	// Insert - Database
	fmt.Println(req.Name)
	idAtual = idAtual + 1

	return &pb.User{
		Id:    strconv.Itoa(idAtual),
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	var userInserted = &pb.User{
		Id:    "1234",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User:   userInserted,
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User:   userInserted,
	})

	time.Sleep(time.Second * 3)

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})

		fmt.Println("Adding", req.GetName())
	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream from the client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("Error sending stream to the client: %v", err)
		}
	}
}
