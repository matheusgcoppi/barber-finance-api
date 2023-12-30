package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"time"
)

func (s *DbRepository) CreateExpense(expense *model.Expense) (*model.Expense, error) {
	var newExpense model.Expense

	insertQuery :=
		`INSERT INTO expenses (created_at, updated_at,account_id, price, description, "when", payment)
		VALUES (?, ?, ?, ?, ?, ?, ?) 
		RETURNING id, created_at, updated_at, deleted_at, account_id, price, description, "when", payment`

	result := s.Store.Db.Raw(
		insertQuery,
		time.Now(),
		time.Now(),
		expense.AccountID,
		expense.Price,
		expense.Description,
		expense.When,
		expense.Payment).Row()

	err := result.Scan(
		&newExpense.ID,
		&newExpense.CreatedAt,
		&newExpense.UpdatedAt,
		&newExpense.DeletedAt,
		&newExpense.AccountID,
		&newExpense.Price,
		&newExpense.Description,
		&newExpense.When,
		&newExpense.Payment,
	)
	if err != nil {
		return nil, err
	}

	return &newExpense, nil
}

func (s *DbRepository) GetExpense() ([]*model.Expense, error) {
	var expenses []*model.Expense
	s.Store.Db.Raw("SELECT * FROM expenses ORDER BY id").Scan(&expenses)

	return expenses, nil
}

func (s *DbRepository) GetExpenseById(id string) (*model.Expense, error) {
	var expense *model.Expense
	query := "SELECT * FROM expenses WHERE id = ?"

	result := s.Store.Db.Raw(query, id).Scan(&expense)

	if expense == nil {
		return nil, fmt.Errorf("user with id %s not found", id)
	}

	if result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, result.Error
	}

	return expense, nil
}

func (s *DbRepository) UpdateExpense(expense *model.ExpenseDTO, id string) (*model.Expense, error) {
	fmt.Println("antes de getByID")
	expenseById, err := s.GetExpenseById(id)
	fmt.Println("depois de getByID")
	if err != nil {
		return nil, err
	}

	fmt.Println(expenseById)

	if expense.Price != 0 {
		expenseById.Price = expense.Price
	}

	fmt.Println("entrou esse erro")

	if expense.Description != "" {
		expenseById.Description = expense.Description
	}

	if !expense.When.IsZero() {
		expenseById.When = expense.When
	}

	if expense.Payment > 0 && expense.Payment < 7 {
		expenseById.Payment = expense.Payment
	}

	fmt.Println("antes de salvar")

	s.Store.Db.Save(expenseById)

	return expenseById, nil
}

func (s *DbRepository) DeleteExpense(id string) error {
	result := s.Store.Db.Exec("DELETE FROM expenses WHERE id = ?", id).RowsAffected
	if result == 0 {
		return fmt.Errorf("Income id = " + id + " not found")
	}

	return nil
}

func (s *DbRepository) CreateNewExpense(accountId uint, price float64, description string, when time.Time, payment model.Payment) *model.Expense {
	return &model.Expense{
		CustomModel: model.CustomModel{},
		AccountID:   accountId,
		Price:       price,
		Description: description,
		When:        when,
		Payment:     payment,
	}
}
