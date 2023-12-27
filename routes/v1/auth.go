package routes

import (
	"fmt"

	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

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
	s.DB.Raw(`SELECT COUNT(*) FROM users WHERE email = ?`, u.Email).Scan(&count)

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

	newUser := user.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(hashedPassword),
	}

	err = s.DB.Create(&newUser).Error
	if err != nil {
		return err
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

	user := user.User{}

	err = s.DB.Where("email = ?", u.Email).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(404, fmt.Sprintf("User with email `%s` was not found.", u.Email))
	}

	if err != nil {
		return err
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

	*user.RememberToken = sess.ID()

	s.DB.Save(&user)

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

	token := sess.ID()

	err = s.DB.
		Model(&user.User{}).
		Where(&user.User{RememberToken: &token}).
		Update("remember_token", "NULL").Error
	if err != nil {
		return err
	}

	sess.Destroy()
	defer sess.Save()

	return c.SendStatus(fiber.StatusNoContent)
}
