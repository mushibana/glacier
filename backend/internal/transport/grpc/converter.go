package grpc

import (
	"github.com/mushibana/glacier/backend/internal/domain"
	pb "github.com/mushibana/glacier/backend/proto"
)

func fromProtoTask(pbTask *pb.Task) domain.Task {
	// TODO:
	return domain.Task{}
}

func toProtoTask(task domain.Task) *pb.Task {
	// TODO:
	return nil
}
