package repository

import "C"
import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"time"
)

func (s *DbRepository) CreateIncome(income *model.Income) (*model.Income, error) {
	newIncome := &model.Income{
		AccountID:   income.AccountID,
		Price:       income.Price,
		Description: income.Description,
		When:        income.When,
		Payment:     income.Payment,
	}

	result := s.Store.Db.Create(newIncome)
	if result.Error != nil {
		return nil, result.Error
	}

	return newIncome, nil
}

func (s *DbRepository) GetIncome() ([]*model.Income, error) {
	var incomes []*model.Income
	s.Store.Db.Raw("SELECT * FROM incomes").Scan(&incomes)

	return incomes, nil
}

func (s *DbRepository) GetIncomeById(id string) (*model.Income, error) {
	var income *model.Income
	query := "SELECT * FROM incomes WHERE id = ?"

	result := s.Store.Db.Raw(query, id).Scan(&income)

	if income == nil {
		return nil, fmt.Errorf("user with id %s not found", id)
	}

	if result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, result.Error
	}

	return income, nil
}

func (s *DbRepository) UpdateIncome(income *model.IncomeDTO, id string) (*model.Income, error) {
	incomeById, err := s.GetIncomeById(id)
	if err != nil {
		return nil, err
	}

	if income.Price != 0 {
		incomeById.Price = income.Price
	}

	if income.Description != "" {
		incomeById.Description = income.Description
	}

	if !income.When.IsZero() {
		incomeById.When = income.When
	}

	if income.Payment > 0 && income.Payment < 7 {
		incomeById.Payment = income.Payment
	}

	s.Store.Db.Save(incomeById)

	return incomeById, nil
}

func (s *DbRepository) DeleteIncome(id string) error {
	result := s.Store.Db.Exec("DELETE FROM incomes WHERE id = ?", id).RowsAffected
	if result == 0 {
		return fmt.Errorf("Income id = " + id + " not found")
	}

	return nil
}

func (s *DbRepository) CreateNewIncome(accountId uint, price float64, description string, when time.Time, payment model.Payment) *model.Income {
	return &model.Income{
		CustomModel: model.CustomModel{},
		AccountID:   accountId,
		Price:       price,
		Description: description,
		When:        when,
		Payment:     payment,
	}
}
