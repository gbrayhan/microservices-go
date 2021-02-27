package models

import (
  "database/sql"
  "errors"
)

type Medicine struct {
  ID          int    `json:"id" example:"123"`
  Name        string `json:"name" example:"Paracetamol"`
  Description string `json:"description" example:"Some Description"`
  EANCode     string `json:"ean_code" example:"9900000124"`
  Laboratory  string `json:"laboratory" example:"Roche"`
  CreatedAt   string `json:"created_at,omitempty" example:"2021-02-24 20:19:39"`
  UpdatedAt   string `json:"updated_at,omitempty" example:"2021-02-24 20:19:39"`
}

func (medicine *Medicine) SaveMedicine() (err error) {
  query := `
	INSERT INTO medicines
	(name, description, ean_code, laboratory)
	VALUES (?, ?, ?, ?); `

  res, err := dbBoilerplateGo.Write.Exec(query, medicine.Name, medicine.Description, medicine.EANCode, medicine.Laboratory)
  if err != nil {
    return
  }

  id, _ := res.LastInsertId()
  medicine.ID = int(id)
  return
}



func GetAllMedicines() (allMedicines []Medicine, err error) {
  var medicine Medicine
  querySelect := `
	SELECT id, name, description, ean_code, laboratory, created_at, updated_at
		FROM medicines `

  rows, err := dbBoilerplateGo.Read.Query(querySelect)
  defer func() {
    err = rows.Close()
  }()
  if err != nil {
    return
  }

  for rows.Next() {
    err = rows.Scan(&medicine.ID, &medicine.Name, &medicine.Description, &medicine.EANCode, &medicine.Laboratory, &medicine.CreatedAt, &medicine.UpdatedAt)
    if err != nil {
      return
    }

    allMedicines = append(allMedicines, medicine)
  }
  return
}

func (medicine *Medicine) GetMedicineByID() (err error) {
  querySelect := `
	SELECT  name, description, ean_code, laboratory, created_at, updated_at
		FROM medicines 
	WHERE id = ?; `

  row := dbBoilerplateGo.Read.QueryRow(querySelect, medicine.ID)
  err = row.Scan(&medicine.Name, &medicine.Description, &medicine.EANCode, &medicine.Laboratory, &medicine.CreatedAt, &medicine.UpdatedAt)

  if errors.Is(err, sql.ErrNoRows) {
    err = nil
  }

  return
}
