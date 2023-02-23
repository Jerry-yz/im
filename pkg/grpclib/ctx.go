package grpclib

import (
	"context"
	"learn-im/pkg/gerrors"
	"strconv"

	"google.golang.org/grpc/metadata"
)

const (
	CtxUserId    = "user_id"
	CtxDeviceId  = "device_id"
	CtxToken     = "token"
	CtxRequestId = "request_id"
)

func ContextWithRequestId(ctx context.Context, requestId int) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(CtxRequestId, strconv.FormatInt(int64(requestId), 10)))
}

func GetCtxRequestId(ctx context.Context) int {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0
	}
	requestId, ok := md[CtxRequestId]
	if !ok && len(requestId) == 0 {
		return 0
	}
	id, err := strconv.ParseInt(requestId[0], 10, 64)
	if err != nil {
		return 0
	}
	return int(id)
}

func GetCtxData(ctx context.Context) (int64, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, gerrors.ErrUserNotFound
	}
	var (
		userId   int64
		deviceId int64
		err      error
	)
	userIds, ok := md[CtxUserId]
	if !ok && len(userIds) == 0 {
		return 0, 0, gerrors.ErrUserNotFound
	}
	userId, err = strconv.ParseInt(userIds[0], 10, 64)
	if err != nil {
		return 0, 0, gerrors.WarpError(err)
	}
	deviceIds, ok := md[CtxDeviceId]
	if !ok && len(deviceIds) == 0 {
		return 0, 0, gerrors.ErrUserNotFound
	}
	deviceId, err = strconv.ParseInt(deviceIds[0], 10, 64)
	if err != nil {
		return 0, 0, gerrors.WarpError(err)
	}
	return userId, deviceId, nil
}

func GetCtxDeviceId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, gerrors.ErrUserNotFound
	}
	deviceIds, ok := md[CtxDeviceId]
	if !ok && len(deviceIds) == 0 {
		return 0, gerrors.ErrNotFound
	}
	deviceId, err := strconv.ParseInt(deviceIds[0], 10, 64)
	if err != nil {
		return 0, gerrors.WarpError(err)
	}
	return deviceId, nil
}

func GetCtxToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", gerrors.ErrNotFound
	}
	tokens, ok := md[CtxToken]
	if !ok && len(tokens) == 0 {
		return "", gerrors.ErrNotFound
	}
	return tokens[0], nil
}

func NewAndCopyRequestId(ctx context.Context) context.Context {
	ctxTmp := context.TODO()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	requestIds, ok := md[CtxRequestId]
	if !ok && len(requestIds) == 0 {
		return nil
	}
	return metadata.NewOutgoingContext(ctxTmp, metadata.Pairs(CtxRequestId, requestIds[0]))
}
