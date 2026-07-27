// Copyright (c) Huawei Technologies Co., Ltd. 2025-2025. All rights reserved."
package response

import (
	"GIDS/common/constants/retcode"
	"GIDS/models/resp"
)

func Success(data interface{}) resp.DataResponse {
	return resp.DataResponse{
		BaseResponse: resp.BaseResponse{
			Code:    retcode.Success,
			Message: "success",
		},
		Data: data,
	}
}
