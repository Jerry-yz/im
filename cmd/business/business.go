package business

import (
	"errors"
	"learn-im/config"
	"learn-im/internal/business"
	"learn-im/pkg/protocol/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var BusinessCmd = &cobra.Command{
	Use:   "business",
	Short: "this is a cmd for business start",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

func Start() {
	server := grpc.NewServer()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		zap.Any("stop server", sig)
		server.GracefulStop()
	}()
	grpcServer(server)
	listenr, err := net.Listen("tcp", config.Conf.BusinessRPCListenAddr)
	if err != nil {
		zap.Error(errors.New("err"))
	}
	if err := server.Serve(listenr); err != nil {
		zap.Error(errors.New("启动服务失败"))
	}
}

func grpcServer(svc *grpc.Server) {
	pb.RegisterBusinessExtServer(svc, business.NewServer())
}
