package transport

import (
	"context"

	"github.com/starathel/gchat/gen/chat"
	"github.com/starathel/gchat/internal/client/ui/components"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ChatClient struct {
	client chat.ChatServiceClient
}

func (c *ChatClient) GetRoomsList() ([]*components.RoomData, error) {
	md := metadata.New(map[string]string{"Authorization": "biba"})
	ctx := metadata.NewOutgoingContext(context.TODO(), md)
	resp, err := c.client.RoomsList(ctx, nil)
	if err != nil {
		return nil, err
	}

	respRooms := make([]*components.RoomData, 0)
	for _, room := range resp.GetRooms() {
		parsedRoom := &components.RoomData{
			Id:         room.GetId(),
			UsersCount: int(room.GetUserCount()),
		}
		respRooms = append(respRooms, parsedRoom)
	}
	return respRooms, nil
}

func NewClient() (*ChatClient, error) {
	conn, err := grpc.NewClient("127.0.0.1:6969", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := chat.NewChatServiceClient(conn)
	return &ChatClient{client: client}, nil
}
