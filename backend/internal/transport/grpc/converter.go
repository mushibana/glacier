package grpc

import (
	"github.com/mushibana/glacier/backend/internal/domain"
	pb "github.com/mushibana/glacier/backend/proto"
)

// Not used anywhere
// func fromProtoTask(pbTask *pb.Task) (domain.Task, error) {
// 	if pbTask == nil {
// 		return domain.Task{}, fmt.Errorf("protobuf task isn't provided")
// 	}
// 	if pbTask.GetId() <= 0 {
// 		return domain.Task{}, fmt.Errorf("task id should be set")
// 	}
// 	if pbTask.GetName() == "" {
// 		return domain.Task{}, fmt.Errorf("task name should be set")
// 	}

// 	t := domain.Task{
// 		ID:          pbTask.GetId(),
// 		ParentID:    pbTask.ParentId,
// 		Name:        pbTask.GetName(),
// 		Description: pbTask.GetDescription(),
// 		StreamID:    pbTask.StreamId,
// 		StartAt:     pbTask.GetStartAt().AsTime(),
// 		DeadlineAt:  pbTask.GetDeadlineAt().AsTime(),
// 		Timelimit:   pbTask.GetTimelimit(),
// 		IsOnetime:   pbTask.GetIsOnetime(),
// 		Repeatable:  pbTask.GetRepeatable(),
// 		Finished:    pbTask.GetFinished(),
// 		CreatedAt:   pbTask.GetCreatedAt().AsTime(),
// 		UpdatedAt:   pbTask.GetUpdatedAt().AsTime(),
// 	}

// 	return t, nil
// }

func toProtoTask(task domain.Task) *pb.Task {
	return &pb.Task{
		Id:          task.ID,
		ParentId:    task.ParentID,
		Name:        task.Name,
		Description: domain.StringOrNil(task.Description),
		StreamId:    task.StreamID,
		StartAt:     timestamppbOrNil(task.StartAt),
		DeadlineAt:  timestamppbOrNil(task.DeadlineAt),
		Timelimit:   domain.Int64OrNil(task.Timelimit),
		IsOnetime:   domain.BoolOrNil(task.IsOnetime),
		Repeatable:  domain.BoolOrNil(task.Repeatable),
		Finished:    domain.BoolOrNil(task.Finished),
		CreatedAt:   timestamppbOrNil(task.CreatedAt),
		UpdatedAt:   timestamppbOrNil(task.UpdatedAt),
	}
}

func toProtoListEntry(tasks []domain.Task) []*pb.ListTaskEntry {
	l := make([]*pb.ListTaskEntry, 0, len(tasks))
	for _, task := range tasks {
		l = append(l, &pb.ListTaskEntry{
			Id:       task.ID,
			ParentId: task.ParentID,
			Name:     task.Name,
			OrderNum: task.OrderNum,
		})
	}
	return l
}
