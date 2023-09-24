package dto

type (
	GetUploadParamsReq struct {
		SId      int64  `json:"s_id" form:"s_id"`
		UId      int64  `json:"u_id" form:"u_id"`
		FileName string `json:"fn" form:"fn"`
	}

	GetUploadParamsRes struct {
		Id     int64             `json:"id"`
		Url    string            `json:"url"`
		Method string            `json:"method"`
		Params map[string]string `json:"params"`
	}

	GetDownloadUrlReq struct {
		UId int64
		Id  int64
	}
)