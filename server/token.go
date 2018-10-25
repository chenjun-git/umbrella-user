package server

import (
	"context"

	"github.com/chenjun-git/umbrella-user/common"
	"github.com/chenjun-git/umbrella-user/pb"
	"github.com/chenjun-git/umbrella-user/token"
	"github.com/chenjun-git/umbrella-user/utils"
)

var verifyMap = map[token.VerifyTokenResult]pb.VerifyResult{
	token.VerifyTokenAllow:      pb.VerifyResult_Allow,
	token.VerifyTokenDeny:       pb.VerifyResult_Deny,
	token.VerifyTokenExpireSoon: pb.VerifyResult_ExpireSoon,
	token.VerifyTokenExpired:    pb.VerifyResult_Expired,
}

func (s *Server) VerifyToken(ctx context.Context, req *pb.VerifyTokenReq) (*pb.VerifyTokenResp, error){
	verifyTokenResult, err := token.VerifyToken(ctx, req.AccessToken)
	if err != nil {
		return nil, common.ExtendErrorStatus(token.VerifyTokenStatusToCode[verifyTokenResult.VerifyResult], err)
	}

	return &pb.VerifyTokenResp{
		Result:   verifyMap[verifyTokenResult.VerifyResult],
		UserId:   verifyTokenResult.UserID,
		Device:   verifyTokenResult.Device,
		App:      verifyTokenResult.App,
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenResp, error) {
	userIDFromCTX, accessToken := utils.GetCurrentTenantAuth(ctx)
	if userIDFromCTX == "" || accessToken == "" {
		return nil, common.NormalErrorStatus(common.AccountBindNoAuthinfo)
	}

	_, _, userIDFromToken, _, err := token.DecodeRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, common.ExtendErrorStatus(common.TokenDeny, err)
	}

	if exist, err := token.RedisCheckTokenExistByToken(ctx, req.RefreshToken); err != nil {
		return nil, common.ExtendErrorStatus(common.AccountRedisError, err)
	} else if !exist {
		return nil, common.NormalErrorStatus(common.TokenDeny)
	}

	accessToken, refreshToken, err := token.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, common.ExtendErrorStatus(common.AccountServiceToken, err)
	}

	// 如果accessToken有效，那么直接返回
	verifyTokenResult, err := token.VerifyToken(ctx, accessToken)
	if err != nil {
		return nil, common.ExtendErrorStatus(token.VerifyTokenStatusToCode[verifyTokenResult.VerifyResult], err)
	} else if verifyTokenResult.VerifyResult == token.VerifyTokenAllow {
		return &pb.RefreshTokenResp{AccessToken: accessToken, RefreshToken: req.RefreshToken}, nil
	}

	return &pb.RefreshTokenResp{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Server) CreatePersonalToken(context.Context, req *pb.CreatePersonalTokenReq) (*pb.CreatePersonalTokenResp, error) {
	userID, accessToken := utils.GetCurrentTenantAuth(ctx)
	if userID == "" || accessToken == "" {
		return nil, common.NormalErrorStatus(common.AccountBindNoAuthinfo)
	}

	if req.Name == "" {
		return nil, common.NormalErrorStatus(common.AccountBindEmptyPersonalTokenName)
	}

	accessToken, err := token.CreateV1ForeverToken(ctx, req.Name, userID)
	if err != nil {
		return nil, common.ExtendErrorStatus(common.AccountServiceToken, err)
	}

	return &pb.CreatePersonalTokenResp{PersonalToken: accessToken}, nil
}

func (s *Server) DeletePersonalToken(ctx context.Context, req *pb.DeletePersonalTokenReq) (*pb.DeletePersonalTokenResp, error) {
	userID, accessToken := utils.GetCurrentTenantAuth(ctx)
	if userID == "" || accessToken == "" {
		return nil, common.NormalErrorStatus(common.AccountBindNoAuthinfo)
	}

	err := token.RemoveForeverToken(ctx, req.Name, userID)
	if err != nil {
		return nil, common.ExtendErrorStatus(common.AccountServiceToken, err)
	}

	return &pb.DeletePersonalTokenResp{}, nil
}



