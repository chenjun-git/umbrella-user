package sms

import (
	"context"

	"google.golang.org/grpc"

	"github.com/chenjun-git/umbrella-sms/pb"
)

func connGrpcService(target string) (*grpc.ClientConn, error) {
	// // add client side tracing interceptor
	// var interceptors = []grpc.UnaryClientInterceptor{traceconf.GlobalTracingConfig.GrpcUnaryClientInterceptor()}

	// // inject caller name
	// if traceconf.GlobalTracingConfig != nil {
	// 	serviceName := traceconf.GlobalTracingConfig.TracerConfig.RecorderConfig.ServiceName
	// 	if serviceName != "" {
	// 		interceptors = append(interceptors, grpccaller.InjectCallerNameUnary(serviceName))
	// 	}
	// }

	//ui := grpcmiddleware.ChainUnaryClient(interceptors...)
	conn, err := grpc.Dial(target, grpc.WithInsecure())//, grpc.WithUnaryInterceptor(ui))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SendPhoneMsg(addr, phone, verifyCode, uid string) error {
	conn, err := connGrpcService(addr)
	if err != nil {
		return err
	}

	client := pb.NewSmsClient(conn)
	req := pb.SmsSendSingleReq{
		Uid:        uid,
		Mobile:     phone,
		VerifyCode: verifyCode,
	}
	_, err = client.SendSignupVerifyCode(context.Background(), &req)
	if err != nil {
		return err
	}

	return nil
}