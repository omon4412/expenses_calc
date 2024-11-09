package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"project/config"
	"project/database"
	"project/logging"
	"project/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Hello(c fiber.Ctx) error {
	return c.SendString("Hello world!")
}

func Register(c fiber.Ctx) error {
	logging.Logger.Info("Received a registration request")
	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if data["username"] == "" || data["email"] == "" || data["password"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}
	logging.Logger.Info("User information",
		zap.String("username", data["username"]),
		zap.String("email", data["email"]),
	)

	var existingUser models.User
	if err := database.DB.Where("email = ?", data["email"]).Or("username = ?", data["username"]).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email or username already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	logging.Logger.Info("Creating User")
	user := &models.User{
		Username: data["username"],
		Email:    data["email"],
		Password: string(hashedPassword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	logging.Logger.Info("User registered successfully")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func Login(c fiber.Ctx) error {
	logging.Logger.Info("Received a Login request")

	var data map[string]string
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	if data["email"] == "" || data["password"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	logging.Logger.Info("User email", zap.String("email", data["email"]))

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0 {
		logging.Logger.Warn("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		logging.Logger.Error("Invalid Password:", zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	logging.Logger.Info("Generating JWT token")
	jwtConfig := config.GetConfig().JWT
	expirationTime := time.Now().Add(jwtConfig.Expiration)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(int(user.ID)),
		"exp": expirationTime.Unix(),
	})
	secretKey := config.GetConfig().JWT.Secret
	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		logging.Logger.Error("Error generating token:", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	logging.Logger.Info("Setting cookie")

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	logging.Logger.Info("Authentication successful, returning")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
	})
}

func User(c fiber.Ctx) error {
	logging.Logger.Info("Request to get user...")

	cookie := c.Cookies("jwt")

	secretKey := config.GetConfig().JWT.Secret
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse claims",
		})
	}

	id, _ := strconv.Atoi((*claims)["sub"].(string))
	user := models.User{ID: uint(id)}

	if err := database.DB.Where("id = ?", user.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		log.Println("Database error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user",
		})
	}

	return c.JSON(user)
}

func Logout(c fiber.Ctx) error {
	logging.Logger.Info("Received a logout request")

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}
