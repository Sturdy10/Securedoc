package repositories

import (
	"database/sql"
	"encoding/json"
	"finalCode/models"
	"fmt"
	"log"

	"github.com/lib/pq"
)

type RepositoryPort interface {
	AddActivityRepositories(req models.RequestActivity) error
	GETuserFileByscdactIDRepositories(scdact_id string) ([][]byte, error)
}

type repositoryAdapter struct {
	db *sql.DB
}

func NewRepositoryAdapter(db *sql.DB) RepositoryPort {
	return &repositoryAdapter{db: db}
}

func (r *repositoryAdapter)  AddActivityRepositories(req models.RequestActivity) error {
	var organizeMemberUUID, organizeMemberRole, organizeMemberDepartment, holderMember string
	err := r.db.QueryRow(`
		SELECT orgmb_id, orgmb_role, orgmb_department, orgmb_holder
		FROM organize_member 
		WHERE orgmb_id = $1`, req.OrgmbatOrgmbid).Scan(&organizeMemberUUID, &organizeMemberRole, &organizeMemberDepartment, &holderMember)
	if err != nil {
		return fmt.Errorf("unable to fetch orgmbat_id: %v", err)
	}

	var orgmbatID string
	err = r.db.QueryRow("SELECT orgmbat_id FROM organize_member_authorization WHERE orgmbat_orgmbid = $1", organizeMemberUUID).Scan(&orgmbatID)
	if err != nil {
		log.Printf("Failed to retrieve data from organize_member_authorization. Error:%v", err)
		return err
	}

	rows, err := r.db.Query(`
		SELECT temjt_orgmb_id
		FROM teamlead_junction
		WHERE temjt_orgdp_id = $1 AND temjt_status = 'Approved';`, organizeMemberDepartment)
	if err != nil {
		log.Printf("Failed to query teamlead from table teamlead_junction. Error:%v", err)
		return err
	}
	defer rows.Close()

	var temjtOrgmbIds []string
	for rows.Next() {
		var temjtOrgmbId string
		if err := rows.Scan(&temjtOrgmbId); err != nil {
			log.Printf("Failed to scan orgmbat_id. Error: %v", err)
			return err
		}
		temjtOrgmbIds = append(temjtOrgmbIds, temjtOrgmbId)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Failed during rows iteration. Error: %v", err)
		return err
	}

	approvers := make(models.Approvers, 0)
	for _, id := range temjtOrgmbIds {
		approvers = append(approvers, models.ApproverStatusArray{Approver: id, Status: "Pending"})
	}
	approversJSON, err := json.Marshal(approvers)
	if err != nil {
		log.Printf("Failed to marshal approvers to JSON. Error: %v", err)
		return err
	}

	var securedocWorkflowUUID string
	err = r.db.QueryRow("INSERT INTO securedoc_workflow (scdwfl_actor, scdwfl_org) VALUES ($1, $2) RETURNING scdwfl_id", orgmbatID, holderMember).Scan(&securedocWorkflowUUID)
	if err != nil {
		log.Printf("Failed to insert data and retrieving scdwfl_id from table securedoc_workflow. Error: %v", err)
		return err
	}

	var securedocWorkflowApproversUUID string
	err = r.db.QueryRow("INSERT INTO securedoc_workflow_approvers (scdwflap_sequence, scdwflap_approver, scdwflap_reqid) VALUES ($1, $2, $3) RETURNING scdwflap_id", securedocWorkflowUUID, approversJSON, req.ScdactReqID).Scan(&securedocWorkflowApproversUUID)
	if err != nil {
		log.Printf("Failed to insert data and retrieving scdwflap_id from table securedoc_workflow_approvers. Error:%v", err)
		return err
	}

	_, err = r.db.Exec("UPDATE securedoc_workflow SET scdwfl_approver = $1 WHERE scdwfl_id = $2", securedocWorkflowApproversUUID, securedocWorkflowUUID)
	if err != nil {
		log.Printf("Failed to update scdwfl_approver in securedoc_workflow. Error: %v", err)
		return err
	}

	for i := range req.HtmlDataSlice {
		_, err := r.db.Exec(`
            INSERT INTO securedoc_activity (
              
            scdact_actor, scdact_reqid, scdact_action, scdact_actiontime, scdact_binary,  scdact_command, scdact_filename, scdact_filetype, scdact_filehash,  scdact_filesize,
		    scdact_filecreated,scdact_filemodified, scdact_filelocation, scdact_status, scdact_name, scdact_type, scdact_reciepient, scdact_starttime, scdact_endtime, scdact_periodday,
            scdact_periodhour, scdact_numberopen, scdact_nolimit, scdact_cvtoriginal, scdact_edit, scdact_print, scdact_copy, scdact_scrwatermark, scdact_watermark, scdact_cvthtml,
            scdact_cvtfcl,  scdact_marcro,  scdact_msgtext,  scdact_subject, scdact_createlocation, scdact_updatelocation, scdact_sender,  scdact_enableconvertoriginal, scdact_approverid ,scdact_orgdp_id, scdact_timestamp
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37 ,$38 ,$39, $40, $41 )`,
			orgmbatID,
			req.ScdactReqID, req.ScdactAction, req.ScdactActionTime, pq.Array([][]byte{req.HtmlDataSlice[i]}), req.ScdactCommand[i], req.ScdactFilename[i], req.ScdactFiletype[i], req.ScdactFilehash[i], req.ScdactFilesize[i], req.ScdactFilecreated[i],
			req.ScdactFilemodified[i], req.ScdactFilelocation[i], req.ScdactStatus, req.ScdactName, req.ScdactType, req.ScdactReciepient, req.ScdactStartTime, req.ScdactEndTime, req.ScdactPeriodDay, req.ScdactPeriodHour,
			req.ScdactNumberOpen, req.ScdactNoLimit || false, req.ScdactCvtOriginal || false, req.ScdactEdit || false, req.ScdactPrint || false, req.ScdactCopy || false, req.ScdactScrWatermark || false, req.ScdactWatermark || false, req.ScdactCvtHtml || false, req.ScdactCvtFcl || false,
			req.ScdactMarcro || false, req.ScdactMsgText, req.ScdactSubject, req.ScdactCreateLocation, req.ScdactUpdateLocation, req.ScdactSender, req.ScdactEnableConvertOriginal || false, securedocWorkflowApproversUUID, req.ScdactDepartmentID, req.ScdactTimestamp,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repositoryAdapter) GETuserFileByscdactIDRepositories(scdact_id string) ([][]byte, error) {
	// SQL query ที่ใช้ parameterized query
	var scdactBinary [][]byte
	err := r.db.QueryRow("SELECT scdact_binary FROM securedoc_activity WHERE scdact_id = $1",
		scdact_id).Scan(pq.Array(&scdactBinary))
	if err != nil {
		return nil, err
	}

	return scdactBinary, nil
}
