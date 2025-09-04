package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"digital-id-server/internal/repository"
	"digital-id-server/shared/types"

	"golang.org/x/crypto/bcrypt"
)

// CreateUsers returns the proportional user list with plain passwords.
func CreateUsers() []types.UserRegisterInput {
	users := []types.UserRegisterInput{
		types.UserRegisterInput{
			FirstName:  "Abdulfetah",
			SecondName: "Suudi",
			LastName:   "Hassen",
			Email:      "superadmin1@oict.com",
			Phone:      "0991752985",
			RoleSlug:   "superadmin",
			Password:   "password1",
		}, types.UserRegisterInput{
			FirstName:  "Abdulfetah",
			SecondName: "Jemal",
			LastName:   "Adem",
			Email:      "admin1@oict.com",
			Phone:      "0961219838",
			RoleSlug:   "admin",
			Password:   "password2",
		}, types.UserRegisterInput{
			FirstName:  "Adnan",
			SecondName: "Tahir",
			LastName:   "Abda",
			Email:      "manager1@oict.com",
			Phone:      "989898989",
			RoleSlug:   "manager",
			Password:   "password3",
		}, types.UserRegisterInput{
			FirstName:  "Jemal",
			SecondName: "Gebi",
			LastName:   "",
			Email:      "executive1@oict.com",
			Phone:      "0900110011",
			RoleSlug:   "executive",
			Password:   "password4",
		}, types.UserRegisterInput{
			FirstName:  "Adem",
			SecondName: "Kedir",
			LastName:   "",
			Email:      "cashier1@oict.com",
			Phone:      "0900110011",
			RoleSlug:   "cashier",
			Password:   "password5",
		}, types.UserRegisterInput{
			FirstName:  "Lammessaa",
			SecondName: "",
			LastName:   "",
			Email:      "encoder1@oict.com",
			Phone:      "0911223344",
			RoleSlug:   "encoder",
			Password:   "password6",
		},
	}
	return users
}

func seedUsers(ctx context.Context, queries *repository.Queries) {
	users := CreateUsers()

	start := time.Now()

	for _, user := range users {
		fmt.Println(user.Email, user.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		user, err := queries.CreateUser(ctx, repository.CreateUserParams{
			FirstName:  user.FirstName,
			SecondName: user.SecondName,
			LastName:   user.LastName,
			Email:      user.Email,
			Phone:      user.Phone,
			RoleSlug:   user.RoleSlug,
		})
		if err != nil {
			log.Fatalf("Failed to seed user: %v", err)
		}
		err = queries.CreateAccount(ctx, repository.CreateAccountParams{
			UserID:       user.ID,
			PasswordHash: string(hashedPassword),
		})
		if err != nil {
			log.Fatalf("Failed to seed account: %v", err)
		}
		// log.Printf("✅ User: %s seeded with role: %s.", created.FirstName, created.Role)
	}

	elapsed := time.Since(start)

	fmt.Printf("\n✅ %d users seeded successfully. Took %s\n", len(users), elapsed)
}
