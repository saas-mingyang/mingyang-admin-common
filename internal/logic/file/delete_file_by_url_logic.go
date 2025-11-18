package file

import (
	"context"

	"github.com/suyuan32/simple-admin-file-tenant/internal/svc"
	"github.com/suyuan32/simple-admin-file-tenant/internal/types"
	"github.com/suyuan32/simple-admin-file-tenant/internal/utils/filex"

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
	fileId, err := filex.ConvertUrlStringToFileUUID(req.Url)
	if err != nil {
		return nil, err
	}

	logic := NewDeleteFileLogic(l.ctx, l.svcCtx)
	return logic.DeleteFile(&types.IdsReq{Ids: []string{fileId}})
}
