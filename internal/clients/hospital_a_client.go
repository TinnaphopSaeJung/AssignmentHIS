package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"his/internal/dto"
)

type HospitalAClient struct {
	baseURL string
	client  *http.Client
}

func NewHospitalAClient() *HospitalAClient {
	baseURL := os.Getenv("HOSPITAL_A_BASE_URL")

	return &HospitalAClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *HospitalAClient) SearchPatient(ctx context.Context, id string) (*dto.HospitalAPatientResponse, error) {
	url := fmt.Sprintf("%s/patient/search/%s", c.baseURL, id)

	/*
		Real external API call example:

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Hospital A API returned status: %d.", resp.StatusCode)
		}

		var result dto.HospitalAPatientResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		return &result, nil
	*/

	_ = url

	// mock data แทน response จากการยิง external api
	mockJSON := map[string]string{
		"3900100445566": `{
			"first_name_th": "ประเสริฐ",
			"middle_name_th": "",
			"last_name_th": "บุญส่ง",
			"first_name_en": "Prasert",
			"middle_name_en": "",
			"last_name_en": "Boonsong",
			"date_of_birth": "1970-03-15",
			"patient_hn": "HN-HY-003",
			"national_id": "3900100445566",
			"passport_id": "",
			"phone_number": "0811112222",
			"email": "prasert.b@example.com",
			"gender": "M"
		}`,
		"P55667788": `{
			"first_name_th": "Maria",
			"middle_name_th": "Garcia",
			"last_name_th": "Lopez",
			"first_name_en": "Maria",
			"middle_name_en": "Garcia",
			"last_name_en": "Lopez",
			"date_of_birth": "1990-06-05",
			"patient_hn": "HN-SK-011",
			"national_id": "",
			"passport_id": "P55667788",
			"phone_number": "0920002222",
			"email": "maria.l@example.com",
			"gender": "F"
		}`,
	}

	mockResponse, ok := mockJSON[id]
	if !ok {
		return nil, errors.New("patient not found from hospital A")
	}

	var result dto.HospitalAPatientResponse
	if err := json.Unmarshal([]byte(mockResponse), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
