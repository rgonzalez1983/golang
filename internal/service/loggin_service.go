package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
)

type TLog struct {
	ServiceName string
	Caller      string
	Event       string
	Extra       string
}

type LogsService interface {
	SaveLog(tlog TLog)
}

type logService struct {
	logger log.Logger
}

func NewLogService(logger log.Logger) LogsService {
	return &logService{
		logger: logger,
	}
}

func (s *logService) SaveLog(tlog TLog) {
	logsURL := fmt.Sprintf("%s:%s%s", os.Getenv("LOGS_HOST"), os.Getenv("LOGS_PORT"), os.Getenv("LOGS_PATH"))
	jsonBody, _ := json.Marshal(tlog)

	client := &http.Client{}
	request, err := http.NewRequest("POST", logsURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		s.logger.Log("error", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	_, netErr := client.Do(request)
	if err != nil {
		s.logger.Log("error", netErr)
	}
}
