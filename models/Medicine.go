package models

type Medicine struct {
  Id          string `json:"id"`
  Name        string `json:"name"`
  Description string `json:"description"`
  EANCode     string `json:"ean_code"`
  Laboratory  string `json:"laboratory"`
}

func (medicine *Medicine) SaveMedicine() (err error) {

  return
}

func GetAllMedicines() (allMedicines []Medicine, err error) {

  return
}

func (medicine *Medicine) GetMedicineById() (err error) {
  querySelect := `
	SELECT  name, description, ean_code, laboratory
		FROM medicines 
	WHERE id = ?; `

  row := dbBoilerplateGo.Read.QueryRow(querySelect, medicine.Id)
  err = row.Scan(&medicine.Name, &medicine.Description, &medicine.EANCode, &medicine.Laboratory)

  return
}
