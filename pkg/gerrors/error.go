package gerrors

import (
	"fmt"
	"runtime"
	"strings"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func WarpError(err error) error {
	if err == nil {
		return nil
	}
	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: "",
		Details: []*anypb.Any{
			{
				TypeUrl: "type_url_stack",
				Value:   []byte(stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func WarpRPCError(err error) error {
	if err == nil {
		return nil
	}
	e, _ := status.FromError(err)
	s := &spb.Status{
		Code:    int32(codes.Unknown),
		Message: "",
		Details: []*anypb.Any{
			{
				TypeUrl: "type_url_stack",
				Value:   []byte(GetErrorStack(e) + "--grpc-- \n" + stack()),
			},
		},
	}
	return status.FromProto(s).Err()
}

func GetErrorStack(s *status.Status) string {
	pbs := s.Proto()
	for i, v := range pbs.Details {
		if v.TypeUrl == "type_url_stack" {
			return string(pbs.Details[i].Value)
		}
	}
	return ""
}

func stack() string {
	pc := make([]uintptr, 20)
	n := runtime.Callers(3, pc)
	var build strings.Builder
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i] - 1)
		file, line := f.FileLine(pc[i] - 1)
		strings.Index(file, "im")
		if n != -1 {
			build.WriteString(fmt.Sprintf("%s:%d", file[n:], line))
		}
	}
	return build.String()
}
