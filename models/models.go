package models

import "mime/multipart"

type ResponseStatus struct {
	Status  string            `json:"Status"`
	Message map[string]string `json:"Message"`
}
type ResponseStatusF struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}
type RequestActivity struct {
	ScdactReqID        string                  `form:"scdact_reqid"`
	ScdactBinary       []*multipart.FileHeader `form:"scdact_binary"`
	ScdactCommand      []string                `form:"scdact_command"`
	HtmlDataSlice      [][]byte
	ScdactFilename     []string `form:"scdact_filename"`
	ScdactFiletype     []string `form:"scdact_filetype"`
	ScdactFilehash     []string `form:"scdact_filehash"`
	ScdactFilesize     []string `form:"scdact_filesize"`
	ScdactFilecreated  []string `form:"scdact_filecreated"`
	ScdactFilemodified []string `form:"scdact_filemodified"`
	ScdactFilelocation []string `form:"scdact_filelocation"`
	ScdactAction       string   `form:"scdact_action"`

	ScdactStatus     string `form:"scdact_status"`
	ScdactName       string `form:"scdact_name"`
	ScdactType       string `form:"scdact_type"`
	ScdactReciepient string `form:"scdact_reciepient"`
	ScdactStartTime  string `form:"scdact_starttime"`
	ScdactEndTime    string `form:"scdact_endtime"`

	ScdactPeriodDay  string `form:"scdact_periodday"`
	ScdactPeriodHour string `form:"scdact_periodhour"`
	ScdactNumberOpen string `form:"scdact_numberopen"`

	ScdactNoLimit               bool   `form:"scdact_nolimit"`
	ScdactCvtOriginal           bool   `form:"scdact_cvtoriginal"`
	ScdactEdit                  bool   `form:"scdact_edit"`
	ScdactPrint                 bool   `form:"scdact_print"`
	ScdactCopy                  bool   `form:"scdact_copy"`
	ScdactScrWatermark          bool   `form:"scdact_scrwatermark"`
	ScdactWatermark             bool   `form:"scdact_watermark"`
	ScdactCvtHtml               bool   `form:"scdact_cvthtml"`
	ScdactCvtFcl                bool   `form:"scdact_cvtfcl"`
	ScdactMarcro                bool   `form:"scdact_marcro"`
	ScdactEnableConvertOriginal bool   `form:"scdact_enableconvertoriginal"`
	ScdactMsgText               string `form:"scdact_msgtext"`
	ScdactSubject               string `form:"scdact_subject"`
	ScdactSender                string `form:"scdact_sender"`
	ScdactCreateLocation        string `form:"scdact_createlocation"`
	ScdactUpdateLocation        string `form:"scdact_updatelocation"`
	ScdactActionTime            string `form:"scdact_actiontime"`
	ScdactDepartmentID          string `form:"scdact_departmentID"`
	ScdactTimestamp             string `form:"scdact_timestamp"`

	OrgmbatOrgmbid string `form:"uuid_member"`
}

type ApproverStatus struct {
	Approver *string `json:"approver"`
	Status   *string `json:"status"`
	// Level    *int    `json:"level"`
}

type StatusResponseByRequestID struct {
	RequestID string           `json:"requestID"`
	Approvals []ApproverStatus `json:"approvals"`
}

type ApproverStatusArray struct {
	Approver string `json:"approver"`
	Status   string `json:"status"`
}

type Approvers []ApproverStatusArray
