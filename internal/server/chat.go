package server

import (
	"sync"

	"github.com/starathel/gchat/gen/chat"
	"github.com/starathel/gchat/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatStream = grpc.BidiStreamingServer[chat.JoinChatRequest, chat.MessageIncoming]

type ChatServer struct {
	chat.UnimplementedChatServiceServer

	mu    sync.RWMutex
	rooms map[string]*ChatRoom
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		rooms: make(map[string]*ChatRoom),
	}
}

func (s *ChatServer) JoinChat(stream chatStream) error {
	err := AuthorizationRequired(stream.Context())
	if err != nil {
		return err
	}

	msg, err := stream.Recv()
	if err != nil {
		return utils.NilIfEOF(err)
	}
	chatId := msg.GetChatId()
	if chatId == "" {
		return status.Error(codes.InvalidArgument, "expected room id in first message")
	}

	s.mu.RLock()
	room, exists := s.rooms[chatId]
	s.mu.RUnlock()
	if !exists {
		room = NewChatRoom(chatId)
		s.mu.Lock()
		s.rooms[chatId] = room
		s.mu.Unlock()
	}

	room.AddUser(stream)

	username := stream.Context().Value(usernameKey).(string)
	for {
		req, err := stream.Recv()
		if err != nil {
			room.RemoveUser(stream)
			return utils.NilIfEOF(err)
		}

		text := req.GetMessage().GetText()
		if text != "" {
			room.SendMessage(username, text)
		}
	}
}
