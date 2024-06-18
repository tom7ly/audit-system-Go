package main

import (
	"audit-system/internal/model"
	"audit-system/internal/service"
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var container *service.Container

func initializeTestContainer() {
	container = service.GetTestContainer("host=localhost port=5433 user=test_pq password=test_pq dbname=test_audit sslmode=disable")
}

func dropTestDatabase() {
	dsn := "host=localhost port=5433 user=pq password=pq dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = 'test_audit'
		  AND pid <> pg_backend_pid();`)
	if err != nil {
		log.Fatalf("failed to terminate existing connections: %v", err)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS test_audit")
	if err != nil {
		log.Fatalf("failed to drop test database: %v", err)
	}
}

func cleanupTestDatabase() {
	if container != nil {
		container.Shutdown()
	}
	dropTestDatabase()
}

func createUser(ctx context.Context, email, name string, age int) (*model.User, error) {
	log.Printf("Creating user %s...\n", name)
	user := model.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	createdUser, err := container.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s created: %+v\n", name, createdUser)
	return createdUser, nil
}

func getUsers(ctx context.Context) ([]model.User, error) {
	log.Println("Getting all users...")
	users, err := container.UserService.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("All users: %+v\n", users)
	return users, nil
}

func getUserByEmail(ctx context.Context, email string) (*model.User, error) {
	log.Printf("Getting user by email %s...\n", email)
	user, err := container.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	log.Printf("User with email %s: %+v\n", email, user)
	return user, nil
}

func updateUser(ctx context.Context, email, name string, age int) (*model.User, error) {
	log.Printf("Updating user %s...\n", name)
	user := model.User{
		Email: email,
		Name:  name,
		Age:   age,
	}
	updatedUser, err := container.UserService.UpdateUser(ctx, email, user)
	if err != nil {
		return nil, err
	}
	log.Printf("User %s updated: %+v\n", name, updatedUser)
	return updatedUser, nil
}

func createAccount(ctx context.Context, email string, balance float64) (*model.Account, error) {
	log.Printf("Creating account for user %s...\n", email)
	account := model.Account{
		Balance:          balance,
		LastTransferTime: time.Now(),
	}
	createdAccount, err := container.AccountService.CreateAccount(ctx, account, email)
	if err != nil {
		return nil, err
	}
	log.Printf("Account for user %s created: %+v\n", email, createdAccount)
	return createdAccount, nil
}

func createTransaction(ctx context.Context, email string, fromAccountID, toAccountID int, amount float64) error {
	log.Printf("Creating transaction from account %d to account %d...\n", fromAccountID, toAccountID)
	err := container.TransactionService.CreateTransaction(ctx, email, fromAccountID, toAccountID, amount)
	if err != nil {
		return err
	}
	log.Println("Transaction created successfully.")
	return nil
}

func verifyBalances(ctx context.Context, t *testing.T, email string, accountID int, expectedBalance float64) {
	updatedAccount, err := container.AccountService.GetAccountByID(ctx, email, accountID)
	assert.NoError(t, err, "failed to get updated account for user")
	assert.Equal(t, expectedBalance, updatedAccount.Balance, "balance of user's account should be updated")
	log.Printf("Updated Account for user %s: %+v\n", email, updatedAccount)
}

func deleteAccount(ctx context.Context, t *testing.T, email string, accountID int) {
	log.Printf("Deleting account %d for user %s...\n", accountID, email)
	err := container.AccountService.DeleteAccount(ctx, accountID)
	assert.NoError(t, err, "failed to delete account")
	log.Printf("Account %d deleted successfully.\n", accountID)

	deletedAccount, err := container.AccountService.GetAccountByID(ctx, email, accountID)
	assert.Error(t, err, "should not be able to fetch deleted account")
	assert.Nil(t, deletedAccount, "deleted account should be nil")
}

func deleteUser(ctx context.Context, t *testing.T, email string) {
	log.Printf("Deleting user %s...\n", email)
	err := container.UserService.DeleteUser(ctx, email)
	assert.NoError(t, err, "failed to delete user")
	log.Printf("User %s deleted successfully.\n", email)

	_, err = container.UserService.GetUserByEmail(ctx, email)
	assert.Error(t, err, "should not be able to fetch deleted user")
}

func getAuditLogs(ctx context.Context) ([]*model.AuditLog, error) {
	log.Println("Getting all audit logs...")
	auditLogs, err := container.AuditLogService.GetAllAuditLogs(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("All audit logs: %+v\n", auditLogs)
	return auditLogs, nil
}

func TestBankingSystem(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user1, err := createUser(ctx, "user1@example.com", "User One", 30)
	assert.NoError(t, err, "failed to create user1")
	account1, err := createAccount(ctx, user1.Email, 1000.0)
	assert.NoError(t, err, "failed to create account for user1")

	user2, err := createUser(ctx, "user2@example.com", "User Two", 25)
	assert.NoError(t, err, "failed to create user2")
	account2, err := createAccount(ctx, user2.Email, 500.0)
	assert.NoError(t, err, "failed to create account for user2")

	err = createTransaction(ctx, user1.Email, account1.ID, account2.ID, 200.0)
	assert.NoError(t, err, "failed to create transaction from user1 to user2")

	verifyBalances(ctx, t, user1.Email, account1.ID, 800.0)
	verifyBalances(ctx, t, user2.Email, account2.ID, 700.0)

	deleteAccount(ctx, t, user1.Email, account1.ID)
	deleteAccount(ctx, t, user2.Email, account2.ID)

	deleteUser(ctx, t, user1.Email)
	deleteUser(ctx, t, user2.Email)
}

func TestCreateAndDeleteUser(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user, err := createUser(ctx, "testuser@example.com", "Test User", 20)
	assert.NoError(t, err, "failed to create test user")

	deleteUser(ctx, t, user.Email)
}

func TestCreateAndDeleteAccount(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user, err := createUser(ctx, "testaccountuser@example.com", "Test Account User", 40)
	assert.NoError(t, err, "failed to create test user for account")

	account, err := createAccount(ctx, user.Email, 1500.0)
	assert.NoError(t, err, "failed to create test account")

	deleteAccount(ctx, t, user.Email, account.ID)

	deleteUser(ctx, t, user.Email)
}

func TestTransactionBetweenAccounts(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user1, err := createUser(ctx, "transactionuser1@example.com", "Transaction User One", 35)
	assert.NoError(t, err, "failed to create transaction user1")
	account1, err := createAccount(ctx, user1.Email, 2000.0)
	assert.NoError(t, err, "failed to create account for transaction user1")

	user2, err := createUser(ctx, "transactionuser2@example.com", "Transaction User Two", 28)
	assert.NoError(t, err, "failed to create transaction user2")
	account2, err := createAccount(ctx, user2.Email, 1000.0)
	assert.NoError(t, err, "failed to create account for transaction user2")

	err = createTransaction(ctx, user1.Email, account1.ID, account2.ID, 500.0)
	assert.NoError(t, err, "failed to create transaction between accounts")

	verifyBalances(ctx, t, user1.Email, account1.ID, 1500.0)
	verifyBalances(ctx, t, user2.Email, account2.ID, 1500.0)

	deleteAccount(ctx, t, user1.Email, account1.ID)
	deleteAccount(ctx, t, user2.Email, account2.ID)

	deleteUser(ctx, t, user1.Email)
	deleteUser(ctx, t, user2.Email)
}

func TestCreateMultipleUsersAndAccounts(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user1, err := createUser(ctx, "multiuser1@example.com", "Multi User One", 30)
	assert.NoError(t, err, "failed to create multi user1")
	account1, err := createAccount(ctx, user1.Email, 1000.0)
	assert.NoError(t, err, "failed to create account for multi user1")

	user2, err := createUser(ctx, "multiuser2@example.com", "Multi User Two", 35)
	assert.NoError(t, err, "failed to create multi user2")
	account2, err := createAccount(ctx, user2.Email, 1500.0)
	assert.NoError(t, err, "failed to create account for multi user2")

	user3, err := createUser(ctx, "multiuser3@example.com", "Multi User Three", 40)
	assert.NoError(t, err, "failed to create multi user3")
	account3, err := createAccount(ctx, user3.Email, 2000.0)
	assert.NoError(t, err, "failed to create account for multi user3")

	verifyBalances(ctx, t, user1.Email, account1.ID, 1000.0)
	verifyBalances(ctx, t, user2.Email, account2.ID, 1500.0)
	verifyBalances(ctx, t, user3.Email, account3.ID, 2000.0)

	deleteAccount(ctx, t, user1.Email, account1.ID)
	deleteAccount(ctx, t, user2.Email, account2.ID)
	deleteAccount(ctx, t, user3.Email, account3.ID)

	deleteUser(ctx, t, user1.Email)
	deleteUser(ctx, t, user2.Email)
	deleteUser(ctx, t, user3.Email)
}

func TestUpdateUserDetails(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user, err := createUser(ctx, "updateuser@example.com", "Update User", 30)
	assert.NoError(t, err, "failed to create update user")

	user.Name = "Updated User"
	user.Age = 31
	_, err = updateUser(ctx, user.Email, user.Name, user.Age)
	assert.NoError(t, err, "failed to update user")

	updatedUser, err := getUserByEmail(ctx, user.Email)
	assert.NoError(t, err, "failed to get updated user")
	assert.Equal(t, "Updated User", updatedUser.Name, "user name should be updated")
	assert.Equal(t, 31, updatedUser.Age, "user age should be updated")

	deleteUser(ctx, t, user.Email)
}

func TestGetAllUsers(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	user1, err := createUser(ctx, "getalluser1@example.com", "Get All User One", 25)
	assert.NoError(t, err, "failed to create get all user1")

	user2, err := createUser(ctx, "getalluser2@example.com", "Get All User Two", 35)
	assert.NoError(t, err, "failed to create get all user2")

	users, err := getUsers(ctx)
	assert.NoError(t, err, "failed to get all users")
	assert.Len(t, users, 2, "there should be two users")

	deleteUser(ctx, t, user1.Email)
	deleteUser(ctx, t, user2.Email)
}

func TestAuditLogs(t *testing.T) {
	initializeTestContainer()
	ctx := context.Background()

	// Create a user to generate an audit log
	user, err := createUser(ctx, "audituser@example.com", "Audit User", 22)
	assert.NoError(t, err, "failed to create audit user")

	// Fetch the audit logs and verify they exist
	auditLogs, err := getAuditLogs(ctx)
	assert.NoError(t, err, "failed to get audit logs")
	assert.NotEmpty(t, auditLogs, "audit logs should not be empty")

	// Wait for TTL to expire and cleanup goroutine to run
	log.Println("Waiting for TTL to expire...")
	time.Sleep(90 * time.Second)

	// Fetch the audit logs again and verify they have been deleted
	auditLogs, err = getAuditLogs(ctx)
	assert.NoError(t, err, "failed to get audit logs")
	assert.Empty(t, auditLogs, "audit logs should be empty after TTL")

	deleteUser(ctx, t, user.Email)
}

func TestMain(m *testing.M) {
	initializeTestContainer()
	code := m.Run()
	cleanupTestDatabase()
	os.Exit(code)
}
