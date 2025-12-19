package filetag

import (
	"context"
	"github.com/saas-mingyang/mingyang-admin-common/utils/convert"

	"mingyang-admin-simple-admin-file/ent/filetag"
	"mingyang-admin-simple-admin-file/internal/svc"
	"mingyang-admin-simple-admin-file/internal/types"
	"mingyang-admin-simple-admin-file/internal/utils/dberrorhandler"

	"github.com/saas-mingyang/mingyang-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileTagLogic {
	return &DeleteFileTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileTagLogic) DeleteFileTag(req *types.IDsReq) (*types.BaseMsgResp, error) {
	_, err := l.svcCtx.DB.FileTag.Delete().Where(filetag.IDIn(convert.StringSliceToUint64Slice(req.Ids)...)).Exec(l.ctx)

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess)}, nil
}
