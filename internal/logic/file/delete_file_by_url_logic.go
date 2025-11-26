package file

import (
	"context"
	"fmt"

	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/filex"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileByUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileByUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileByUrlLogic {
	return &DeleteFileByUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileByUrlLogic) DeleteFileByUrl(req *types.FileDeleteReq) (resp *types.BaseMsgResp, err error) {
	fileId, err := filex.ConvertUrlStringToFileUint64(req.Url)
	if err != nil {
		return nil, err
	}

	logic := NewDeleteFileLogic(l.ctx, l.svcCtx)
	return logic.DeleteFile(&types.IdsReq{Ids: []string{fmt.Sprint(fileId)}})
}
