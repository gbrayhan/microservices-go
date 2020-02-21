package models

type ExampleElement struct {
	ID          int    `json:"ID,omitempty"`
	FullName    string `json:"fullName,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	JobPosition string `json:"jobPosition,omitempty"`
}

func (element *ExampleElement) CompleteDataID(response *ResponseBase) {
	querySelect := `
	SELECT full_name, email, phone, job_position
		FROM name_database.name_table 
	WHERE id = ?; `

	row := database.DBRead.QueryRow(querySelect, element.ID)
	err := row.Scan(&element.FullName, &element.Email, &element.Phone, &element.JobPosition)

	if err != nil {
		response.AppendFatalError("SQL Error: " + err.Error())
		return
	}

}
