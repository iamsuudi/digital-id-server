package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/iamsuudi/digital-id-server/internal/repository"
	"github.com/iamsuudi/digital-id-server/shared/types"
	"golang.org/x/crypto/bcrypt"
)

// CreateUsers returns the proportional user list with plain passwords.
func CreateUsers() []types.UserRegisterInput {
	var users []types.UserRegisterInput
	id := 0
	next := func(role, city, sc, k string) types.UserRegisterInput {
		id++
		return types.UserRegisterInput{
			FirstName:  fmt.Sprintf("%s-%s", city, sc),
			SecondName: k,
			LastName:   strconv.Itoa(id),
			Email:      strings.ToLower(fmt.Sprintf("%s%d@oict.com", role, id)),
			Phone:      fmt.Sprintf("+251911%06d", id),
			Password:   fmt.Sprintf("password%d", id), // simple, unique password
			Role:       role,
		}
	}

	// 1 superadmin
	users = append(users, next("SUPERADMIN", "system", "", ""))

	// 1 admin per city
	for _, city := range Data {
		users = append(users, next("ADMIN", city.Name, "", ""))
	}

	// 1 manager per sub-city
	for _, city := range Data {
		for _, sc := range city.SubCities {
			users = append(users, next("MANAGER", city.Name, sc.Name, ""))
		}
	}

	// 3 executives/encoders/cashiers per kebele
	for _, city := range Data {
		for _, sc := range city.SubCities {
			for _, kb := range sc.Kebeles {
				for _, role := range []string{"EXECUTIVE", "ENCODER", "CASHIER"} {
					users = append(users, next(role, city.Name, sc.Name, kb.Name))
				}
			}
		}
	}
	return users
}

func seedUsers(ctx context.Context, queries *repository.Queries) {
	users := CreateUsers()

	start := time.Now()

	for _, user := range users {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		_, err = queries.CreateUser(ctx, repository.CreateUserParams{
			FirstName:    user.FirstName,
			SecondName:   user.SecondName,
			LastName:     user.LastName,
			Email:        user.Email,
			Phone:        user.Phone,
			PasswordHash: string(hashedPassword),
			Role:         user.Role,
		})
		if err != nil {
			log.Fatalf("Failed to seed user: %v", err)
		}
		// log.Printf("✅ User: %s seeded with role: %s.", created.FirstName, created.Role)
	}

	elapsed := time.Since(start)

	fmt.Printf("\n✅ %d users seeded successfully. Took %s\n", len(users), elapsed)
}
