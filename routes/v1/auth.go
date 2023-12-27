package routes

import (
	"fmt"
	"time"

	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"golang.org/x/crypto/bcrypt"
)

func (c *Controller) RegisterAuthRoutes() {
	v1 := c.App.Group("/v1")

	v1.Post("/signup", c.signupPost)
	v1.Post("/login", c.loginPost)
	v1.Post("/logout", c.logoutPost)
}

// signupPost godoc
//
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		user.UserSignUpDto	true	"Sign Up Data"
//	@Success	201		{object}	user.User
//	@Router		/v1/signup [post]
func (s *Controller) signupPost(c *fiber.Ctx) error {
	u := new(user.UserSignUpDto)

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

	// Ignore errors here, only checking existence
	s.OldDB.QueryRow(`SELECT COUNT(*) as count FROM users WHERE email = ?`, u.Email).Scan(&count)

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

	_, err = s.OldDB.Exec(`
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

	newUser := new(user.User)

	err = s.OldDB.Get(newUser, `SELECT * FROM users WHERE email = ?`, u.Email)
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
//	@Param		request	body		user.UserLoginDto	true	"Login Data"
//	@Success	200		{object}	user.User
//	@Router		/v1/login [post]
func (s *Controller) loginPost(c *fiber.Ctx) error {
	u := new(user.UserLoginDto)

	var err error

	err = c.BodyParser(&u)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(u)
	if err != nil {
		return err
	}

	user := new(user.User)

	err = s.OldDB.Get(user, `SELECT * FROM users WHERE email = ?;`, u.Email)
	if err != nil {
		return fiber.NewError(404, fmt.Sprintf("User with email `%s` was not found.", u.Email))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Wrong Password")
	}

	var sess *session.Session

	sess, err = s.Store.Get(c)
	if err != nil {
		return err
	}

	sess.Regenerate()
	defer sess.Save()

	newToken := sess.ID()

	_, err = s.OldDB.Exec(`UPDATE users SET remember_token = ? WHERE id = ?;`, newToken, user.Id)
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
func (s *Controller) logoutPost(c *fiber.Ctx) error {
	sess, err := s.Store.Get(c)
	if err != nil {
		return err
	}

	s.OldDB.Exec(`UPDATE users SET remember_token = 'NULL' WHERE remember_token = ?;`, sess.ID())

	sess.Destroy()
	defer sess.Save()

	return c.SendStatus(fiber.StatusNoContent)
}
