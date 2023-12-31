package rest

import (
	"fmt"
	"net/http"

	"test/logger"
	"test/models"
	"test/service"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type userController struct {
	userService service.UserService
	logger      logger.Logger
}

type UserController interface {
	CreateUserAccount(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
}

func NewUserController(userService service.UserService, logger logger.Logger) UserController {
	return &userController{
		userService: userService,
		logger:      logger,
	}
}

// CreateUserAccount handles the HTTP request to create a new user account.
// It expects a JSON request body containing user signup details.
// Parameters:
//   - ctx: The Gin context representing the HTTP request and response.
func (u *userController) CreateUserAccount(ctx *gin.Context) {
	// Declare a variable to hold the user signup request
	var request models.UserSignUpRequest

	// Bind the JSON request body to the request variable
	err := ctx.BindJSON(&request)
	if err != nil {
		u.logger.Errorf(fmt.Sprintf("error in binding JSON in CreateUserAccount API err: %v", err.Error()))

		// Return a JSON response for a bad request due to invalid input
		ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_INPUT, models.INVALID_INPUT_MESSAGE, nil))
		return
	}

	// Log the request details
	u.logger.Info(fmt.Sprintf("request to CreateUserAccount with body: %v", request))

	// Validate the request body
	err = request.Validate()
	if err != nil {
		u.logger.Errorf(fmt.Sprintf("error in request body JSON in CreateUserAccount API err: %v", err.Error()))

		// Return a JSON response for bad request with appropriate error message
		if err.Error() == models.CHOOSE_BETTER_PASSCODE_MESSAGE {
			ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.CHOOSE_BETTER_PASSCODE, models.CHOOSE_BETTER_PASSCODE_MESSAGE, nil))
			return
		}
		ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_INPUT, models.INVALID_INPUT_MESSAGE, nil))
		return
	}

	// Log additional details about the signup request
	logger.Instance().Info("===== Signup Request ==========")
	logger.Instance().Info("body ", request.FirstName, request.LastName, request.Email)

	// Call the user service to create the user
	err = u.userService.CreateUser(request)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error in user registration API err: %v", err.Error()))

		// Return a JSON response for internal server error with appropriate error message
		if standardError, ok := err.(*service.StandardError); ok {
			ctx.JSON(http.StatusInternalServerError, NewStandardResponse(false, standardError.Code, standardError.Message, nil))
			err = errors.Wrapf(err, "%+v", errors.New(standardError.Message))
			ctx.Errors = append(ctx.Errors, ctx.Error(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, NewStandardResponse(false, models.USER_EXIST, models.USER_ALREADY_EXIST, err))
			return
		}
	}

	// Return a JSON response for successful user creation
	ctx.JSON(http.StatusOK, NewStandardResponse(true, models.SUCCESS, models.PERSONAL_USER_CREATED, nil))
}

// Login handles the HTTP request for user login.
// It expects a JSON request body containing user sign-in details.
// Upon successful login, it generates a JWT token and returns it in the response.
// Parameters:
//   - ctx: The Gin context representing the HTTP request and response.
func (u *userController) Login(ctx *gin.Context) {
	// Declare a variable to hold the user sign-in request
	var request models.UserSignInRequest

	// Bind the JSON request body to the request variable
	err := ctx.BindJSON(&request)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error in binding JSON in user SignIn API err: %v", err.Error()))

		// Return a JSON response for a bad request due to invalid input
		ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_INPUT, models.INVALID_INPUT_MESSAGE, nil))
		return
	}

	// Log the request details
	u.logger.Info(fmt.Sprintf("request for user SignIn with body: %v", request))

	// Validate the request body
	err = request.Validate()
	if err != nil {
		u.logger.Error(fmt.Sprintf("error in request body JSON in user SignIn API err: %v", err.Error()))

		// Return a JSON response for unauthorized due to invalid input
		ctx.JSON(http.StatusUnauthorized, NewStandardResponse(false, models.INVALID_INPUT, models.INVALID_INPUT_MESSAGE, nil))
		return
	}

	// Call the user service to perform login
	token, err := u.userService.Login(request)
	if err != nil {
		// Return a JSON response for bad request with appropriate error message
		ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_INPUT, models.INVALID_CREDENTIALS, nil))
		return
	}

	// Return a JSON response for successful login with the generated token
	ctx.JSON(http.StatusOK, NewStandardResponse(true, models.SUCCESS, models.LOGIN_SUCCESSFUL, token))
}

// GetUserByEmail handles the HTTP request to retrieve user details based on the user's email.
// It extracts the user's email from the context and calls the user service to get the user details.
// Parameters:
//   - ctx: The Gin context representing the HTTP request and response.
func (u *userController) GetUserByEmail(ctx *gin.Context) {
	// Get the current user's email from the context
	email := getCurrentUserEmail(ctx)

	// Call the user service to get user details by email
	userDetail, err := u.userService.GetUserByEmail(email)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error in getting user details API, err: %v", err.Error()))

		// Return a JSON response for bad request with appropriate error message
		ctx.JSON(http.StatusBadRequest, NewStandardResponse(false, models.INVALID_INPUT, models.INTERNAL_SERVER_ERROR_MESSAGE, nil))
		return
	}

	// If the user is not found, return a JSON response with status not found
	if userDetail.Id == 0 {
		u.logger.Error(fmt.Sprintf("error in getting user details API, err: %v", err.Error()))
		ctx.JSON(http.StatusNotFound, NewStandardResponse(false, http.StatusNotFound, models.USER_NOT_FOUND, nil))
		return
	}

	// Return a JSON response for successful retrieval of user details
	ctx.JSON(http.StatusOK, NewStandardResponse(true, models.SUCCESS, models.USER_FETCHED_SUCCESSFULLY, userDetail))
}
