package services

import (
	"errors"
	"finalCode/models"
	"finalCode/repositories"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kballard/go-shellquote"
)

type ServicesPort interface {
	AddActivityServices(req models.RequestActivity, c *gin.Context, ScdactBinary []*multipart.FileHeader, ScdactCommands []string) (int, error)
	DeleteByScdactReqIDServices(req models.RequestActivity) error
	GETuserFileByscdactIDServices(scdact_id string) ([][]byte, error)
}

type serviceAdapter struct {
	r repositories.RepositoryPort
}

func NewServiceAdapter(r repositories.RepositoryPort) ServicesPort {
	return &serviceAdapter{r: r}
}

var (
	BaseDirectory     = "D:\\Developer\\Back-end\\Project\\securedoc"
	Datapath          = filepath.Join(BaseDirectory, "API3.40R02.0001\\64bit\\data")
	Cmdfinal_codepath = filepath.Join(BaseDirectory, "API3.40R02.0001\\64bit\\bin_win")
)

func (s *serviceAdapter) AddActivityServices(req models.RequestActivity, c *gin.Context, ScdactBinary []*multipart.FileHeader, ScdactCommand []string) (int, error) {
	// สร้างโฟลเดอร์หลัก
	mainFolderPath := filepath.Join(Datapath, req.ScdactReqID)
	if err := os.MkdirAll(mainFolderPath, os.ModePerm); err != nil {
		return 0, err
	}

	// บันทึกไฟล์ที่อัปโหลด
	for i, file := range ScdactBinary {
		if i >= len(ScdactCommand) {
			return 0, errors.New("insufficient data for index " + strconv.Itoa(i))
		}
		filePath := filepath.Join(mainFolderPath, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			return 0, err
		}
	}

	for i, cmd := range ScdactCommand {
		if i >= len(req.ScdactFilename) ||
			i >= len(req.ScdactFiletype) ||
			i >= len(req.ScdactFilehash) ||
			i >= len(req.ScdactFilesize) ||
			i >= len(req.ScdactFilecreated) ||
			i >= len(req.ScdactFilemodified) ||
			i >= len(req.ScdactFilelocation) {
			return 0, errors.New("insufficient data for index " + strconv.Itoa(i))
		}

		// command := exec.Command("bash", "-c", cmd) // for linux
		// command.Dir = Cmdfinal_codepath

		cmdArgs, err := shellquote.Split(cmd)
		if err != nil {
			return 0, err
		}

		command := exec.Command("cmd", append([]string{"/C"}, cmdArgs...)...)
		command.Dir = Cmdfinal_codepath

		_, err = command.CombinedOutput()

		exitCode := 0
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		if err != nil {
			return exitCode, err
		}
	}

	htmlFiles, err := filepath.Glob(filepath.Join(Datapath, req.ScdactReqID, "*.html"))
	if err != nil {
		return 0, err
	}

	for _, htmlFilePath := range htmlFiles {
		htmlData, err := os.ReadFile(htmlFilePath)
		if err != nil {
			return 0, err
		}
		req.HtmlDataSlice = append(req.HtmlDataSlice, htmlData)
		// fmt.Printf("ไฟล์: %s\n", htmlFilePath)
		// fmt.Printf("ข้อมูล: %s\n", htmlData)
	}

	// Insert data into the database
	err = s.r.AddActivityRepositories(req)
	if err != nil {
		return 0, err

	}

	return 0, nil
}

func (s *serviceAdapter) DeleteByScdactReqIDServices(req models.RequestActivity) error {
	folderPath := filepath.Join(Datapath, req.ScdactReqID)
	return os.RemoveAll(folderPath)
}

func (s *serviceAdapter) GETuserFileByscdactIDServices(scdact_id string) ([][]byte, error) {
	// เรียกใช้ GETUser_File ใน Repository
	scdactBinary, err := s.r.GETuserFileByscdactIDRepositories(scdact_id)
	if err != nil {
		return nil, err
	}

	return scdactBinary, nil
}
