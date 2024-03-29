package logic

import (
	"context"

	"github.com/satyarth42/chatter/dto"
)

func Login(ctx context.Context, req *dto.LoginReq) (dto.LoginResp, *dto.CommonError) {
	return dto.LoginResp{}, nil
}
