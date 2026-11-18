package grpc

import (
	"context"
	"fmt"

	"github.com/mushibana/glacier/backend/internal/domain"
	"github.com/mushibana/glacier/backend/internal/service"
	pb "github.com/mushibana/glacier/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	pb.UnimplementedTasksServer
	taskService service.TaskService
}

//	func (h *TaskHandler) GetTasks(context.Context, *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
//		return nil, status.Error(codes.Unimplemented, "method GetTasks not implemented")
//	}
//
//	func (h *TaskHandler) CreateTask(context.Context, *pb.CreateTaskRequest) (*pb.Task, error) {
//		return nil, status.Error(codes.Unimplemented, "method CreateTask not implemented")
//	}

func (h *TaskHandler) UpdateTask(ctx context.Context, r *pb.UpdateTaskRequest) (*pb.Task, error) {
	if r.GetTask() == nil {
		return nil, status.Error(codes.InvalidArgument, "task is required")
	}
	if r.GetTask().GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "task id should be set")
	}
	if r.GetUpdateMask() == nil {
		return nil, status.Error(codes.InvalidArgument, "update mask is required")
	}
	r.GetUpdateMask().Normalize()
	if !r.GetUpdateMask().IsValid(&pb.Task{}) {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid field mask: %s", r.GetUpdateMask().String()))
	}

	params := domain.UpdateTaskParams{ID: r.GetTask().GetId()}

	for _, path := range r.GetUpdateMask().GetPaths() {
		switch path {
		case "name":
			params.Name = domain.SetOptional(r.GetTask().GetName())
		case "description":
			params.Description = domain.SetOptional(r.GetTask().GetDescription())
		case "stream_id":
			params.StreamID = domain.SetOptional(r.GetTask().GetStreamID())
		case "start_at":
			params.StartAt = domain.SetOptional(r.GetTask().GetStartAt())
		case "deadline_at":
			params.DeadlineAt = domain.SetOptional(r.GetTask().GetDeadlineAt())
		case "timelimit":
			params.TimeLimit = domain.SetOptional(r.GetTask().GetTimeLimit())
		case "is_onetime":
			params.IsOnetime = domain.SetOptional(r.GetTask().GetIsOnetime())
		case "repeatable":
			params.Repeatable = domain.SetOptional(r.GetTask().GetRepeatable())
		case "finished":
			params.Finished = domain.SetOptional(r.GetTask().GetFinished())
		default:
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("unsupported path for update in field mask: %s", path))
		}
	}

	task, err := h.taskService.Update(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("partial task update: %v", err))
	}

	return toProtoTask(task), nil
}

// func (h *TaskHandler) DeleteTask(context.Context, *pb.DeleteTaskRequest) (*emptypb.Empty, error) {
// 	return nil, status.Error(codes.Unimplemented, "method DeleteTask not implemented")
// }
// func (h *TaskHandler) MoveTask(context.Context, *pb.MoveTaskRequest) (*pb.MoveTaskResponse, error) {
// 	return nil, status.Error(codes.Unimplemented, "method MoveTask not implemented")
// }
