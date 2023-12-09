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
	}

	result := s.Store.Db.Create(newIncome)
	if result.Error != nil {
		return nil, result.Error
	}

	return newIncome, nil
}

func (s *DbRepository) GetIncome() ([]*model.Income, error) {
	var incomes []*model.Income
	result := s.Store.Db.Order("id").Find(incomes)
	if result.Error != nil {
		return nil, result.Error
	}
	return incomes, nil
}

func (s *DbRepository) GetIncomeById(id string) (*model.Income, error) {
	var income *model.Income
	query := "SELECT * FROM incomes WHERE id = ?"

	result := s.Store.Db.Raw(query, id).Scan(&income)

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

	s.Store.Db.Save(incomeById)

	return incomeById, nil
}

func (s *DbRepository) DeleteIncome(id string) error {
	result := s.Store.Db.Delete(&model.Income{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *DbRepository) CreateNewIncome(price float64, description string, when time.Time, payment model.Payment) *model.Income {
	return &model.Income{
		CustomModel: model.CustomModel{},
		Price:       price,
		Description: description,
		When:        when,
		Payment:     payment,
	}
}
