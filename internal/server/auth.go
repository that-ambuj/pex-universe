package server

import (
	"fmt"
	"pex-universe/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"golang.org/x/crypto/bcrypt"
)

func (s *FiberServer) RegisterAuthRoutes() {
	v1 := s.App.Group("/v1")

	v1.Post("/signup", s.signupPost)
	v1.Post("/login", s.loginPost)
	v1.Post("/logout", s.logoutPost)
}

// signupPost godoc
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.UserSignUpDto	true	"Sign Up Data"
//	@Success	201		{object}	model.User
//	@Router		/v1/signup [post]
func (s *FiberServer) signupPost(c *fiber.Ctx) error {
	u := new(model.UserSignUpDto)

	var err error

	err = c.BodyParser(u)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(u)
	if err != nil {
		return err
	}

	count := 0

	// Ignore errors here, only checking existense
	s.db.QueryRow(`SELECT COUNT(*) as count FROM users WHERE email = ?`, u.Email).Scan(&count)

	if count > 0 {
		return &fiber.Error{
			Code:    400,
			Message: fmt.Sprintf("User with email %s already exists.", u.Email),
		}
	}

	hashedPassword := make([]byte, 256)

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO
			users (name, email, password, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?);`,
		u.Name,
		u.Email,
		hashedPassword,
		time.Now(),
		time.Now())
	if err != nil {
		return err
	}

	newUser := new(model.User)

	err = s.db.Get(newUser, `SELECT * FROM users WHERE email = ?`, u.Email)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

// loginPost godoc
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		model.UserLoginDto	true	"Login Data"
//	@Success	200		{object}	model.User
//	@Router		/v1/login [post]
func (s *FiberServer) loginPost(c *fiber.Ctx) error {
	u := new(model.UserLoginDto)

	var err error

	err = c.BodyParser(&u)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(u)
	if err != nil {
		return err
	}

	user := new(model.User)

	err = s.db.Get(user, `SELECT * FROM users WHERE email = ?;`, u.Email)
	if err != nil {
		return fiber.NewError(404, fmt.Sprintf("User with email `%s` was not found.", u.Email))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Wrong Password")
	}

	var sess *session.Session

	sess, err = s.store.Get(c)
	if err != nil {
		return err
	}

	sess.Regenerate()
	defer sess.Save()

	newToken := sess.ID()

	_, err = s.db.Exec(`UPDATE users SET remember_token = ? WHERE id = ?;`, newToken, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// loginPost godoc
//
//	@Tags		auth
//	@Summary	Log out of the current session
//	@Router		/v1/logout [post]
func (s *FiberServer) logoutPost(c *fiber.Ctx) error {
	sess, err := s.store.Get(c)
	if err != nil {
		return err
	}

	s.db.Exec(`UPDATE users SET remember_token = 'NULL' WHERE remember_token = ?;`, sess.ID())

	sess.Destroy()
	// defer sess.Save()

	return c.SendStatus(fiber.StatusNoContent)
}
