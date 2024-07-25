package greenhouse

import (
	"encoding/json"
	"io"
	"net/http"
)

// EditGreenhouse
type EditGreenhouse struct {
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Zip         string `json:"zip"`
	Address     string `json:"address"`
}

// EditGreenhouseV1Params
type EditGreenhouseV1Params struct {
	ID   string `json:"UUID"`
	GUID string `json:"GUID"`

	DisplayName *string `json:"DisplayName"`
	Status      *string `json:"Status"`
	Zip         *string `json:"Zip"`
	Address     *string `json:"Address"`
}

// EditGreenhouseV1 edits a Greenhouse with given arguments
func (l *Greenhouse) EditGreenhouseV1(p *EditGreenhouseV1Params) error {
	var query string
	var arguments = []any{}
	query += "UPDATE greenhouses SET "
	if *p.DisplayName != "" {
		arguments = append(arguments, *p.DisplayName)
		query += "display_name=?,"
	}
	if *p.Status != "" {
		arguments = append(arguments, *p.Status)
		query += "status=?,"
	}
	if *p.Zip != "" {
		arguments = append(arguments, *p.Zip)
		query += "zip=?,"
	}
	if *p.Address != "" {
		arguments = append(arguments, *p.Address)
		query += "address=?,"
	}

	//delete comma
	query = query[:len(query)-1]
	query += " WHERE guid=?;"

	//uuid as last argument
	arguments = append(arguments, p.GUID)

	_, err := l.dbh.Exec(query, arguments...)
	if err != nil {
		return err
	}
	return nil
}

// EditGreenhouseHandler handles editing one Greenhouse
func (l *Greenhouse) EditGreenhouseHandler(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	req := &EditGreenhouseV1Params{}

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, req)
	err := l.EditGreenhouseV1(req)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal("success")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
