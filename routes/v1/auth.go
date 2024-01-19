package routes

import (
	"fmt"
	"time"

	"pex-universe/internal/errors"
	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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
		return errors.BadRequestErr(err)
	}

	err = s.ValidateStruct(u)
	if err != nil {
		log.Error(err)
		return err
	}

	count := 0

	// Ignore errors here, only checking existence
	s.DB.Raw(`SELECT COUNT(*) FROM site_users WHERE email = ?`, u.Email).Scan(&count)

	if count > 0 {
		return errors.BadRequestMsg(
			fmt.Sprintf("User with email %s already exists.", u.Email))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		log.Error(err)
		return err
	}

	newUser := user.User{
		Name:          u.Name,
		Email:         u.Email,
		Username:      u.Username,
		Password:      string(hashedPassword),
		RetAdLastSent: time.Now(),
	}

	err = s.DB.Create(&newUser).Error
	if err != nil {
		log.Error(err)
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
	u := user.UserLoginDto{}

	err := c.BodyParser(&u)
	if err != nil {
		return errors.BadRequestErr(err)
	}

	err = s.ValidateStruct(&u)
	if err != nil {
		log.Error(err)
		return err
	}

	user := user.User{}

	err = s.DB.Where("email = ?", u.Email).First(&user).Error

	switch err {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		return errors.NotFoundEntity("User", "email", u.Email)
	default:
		log.Error(err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		log.Error(err)
		return errors.UnauthorizedMsg("Wrong Password")
	}

	var sess *session.Session

	sess, err = s.Store.Get(c)
	if err != nil {
		log.Error(err)
		return err
	}

	err = sess.Regenerate()
	if err != nil {
		log.Error(err)
		return err
	}

	var (
		token = sess.ID()
		now   = time.Now()
	)

	user.RememberToken = &token
	user.LastLoggedIn = &now

	err = s.DB.Save(&user).Error
	if err != nil {
		log.Error(err)
		return err
	}

	err = sess.Save()
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&user)
}

// loginPost godoc
//
//	@Tags		auth
//	@Summary	Log out of the current session
//	@Router		/v1/logout [post]
func (s *Controller) logoutPost(c *fiber.Ctx) error {
	sess, err := s.Store.Get(c)
	if err != nil {
		log.Error(err)
		return err
	}

	token := sess.ID()

	err = s.DB.
		Model(&user.User{}).
		Where(&user.User{RememberToken: &token}).
		Update("remember_token", nil).Error
	if err != nil {
		log.Error(err)
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "cart_id",
		Value:   "",
		Expires: time.Now(),
	})

	//nolint
	sess.Destroy()
	//nolint
	defer sess.Save()

	return c.SendStatus(fiber.StatusNoContent)
}
