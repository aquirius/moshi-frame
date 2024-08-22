package greenhouse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	ID   string `json:"uuid"`
	GUID uint64 `json:"guid"`

	DisplayName *string  `json:"display_name"`
	Status      *string  `json:"status"`
	Zip         *string  `json:"zip"`
	Address     *string  `json:"address"`
	Destination *string  `json:"destination"`
	TempIn      *float64 `json:"tempIn"`
	TempOut     *float64 `json:"tempOut"`
	Humidity    *float64 `json:"humidity"`
	Brightness  *float64 `json:"brightness"`
	Co2         *float64 `json:"co2"`
}

// EditGreenhouseV1 edits a Greenhouse with given arguments
func (l *Greenhouse) EditGreenhouseV1(p *EditGreenhouseV1Params) error {
	var query string
	var arguments = []any{}

	fmt.Println(*p.DisplayName)
	fmt.Println(*p.Address)
	fmt.Println(*p.Zip)
	fmt.Println(*p.TempIn)
	fmt.Println(*p.TempOut)
	fmt.Println(*p.Humidity)
	fmt.Println(p.GUID)

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
	if *p.Destination != "" {
		arguments = append(arguments, *p.Destination)
		query += "destination=?,"
	}
	if *p.TempIn != 0 {
		arguments = append(arguments, *p.TempIn)
		query += "tempIn=?,"
	}
	if *p.TempOut != 0 {
		arguments = append(arguments, *p.TempOut)
		query += "tempOut=?,"
	}
	if *p.Humidity != 0 {
		arguments = append(arguments, *p.Humidity)
		query += "humidity=?,"
	}
	if *p.Brightness != 0 {
		arguments = append(arguments, *p.Brightness)
		query += "brightness=?,"
	}
	if *p.Co2 != 0 {
		arguments = append(arguments, *p.Co2)
		query += "co2=?,"
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
	var vars = mux.Vars(r)
	guid, _ := strconv.ParseUint(vars["guid"], 0, 32)
	req := &EditGreenhouseV1Params{
		GUID: guid,
	}

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
