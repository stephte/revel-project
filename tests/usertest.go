package tests

import (
	"revel-project/app/services/dtos"
	"encoding/json"
	"fmt"
)

type UserTest struct {
	BaseTest
}

func (t *UserTest) Before() {
	fmt.Println("----- Initializing Auth -----")
	t.InitAuth()
}

func (t *UserTest) After() {
	fmt.Println("----- Cleaning up -----")
	t.CleanupAuth()
}

// ----- User Index -----

func (t *UserTest) TestThatUserIndexRequiresAuth() {
	t.Get("/users")
	t.Assert(t.Response.StatusCode == 401)
}


func (t *UserTest) TestThatUserIndexRequiresAdminAuth() {
	t.SendAsRegularUser("get", "/users", nil)
	t.Assert(t.Response.StatusCode == 401)
}


func (t *UserTest) TestThatUserIndexAcceptsAdminAndSuperAdminAuth() {
	t.SendAsAdmin("get", "/users", nil)
	t.Assert(t.Response.StatusCode == 200)


	t.SendAsSuperAdmin("get", "/users", nil)
	t.Assert(t.Response.StatusCode == 200)
}

// ----- User find -----

func (t *UserTest) TestUserFindWorks() {
	t.SendAsRegularUser("get", fmt.Sprintf("/users/%s", t.regularUser.Key), nil)
	t.Assert(t.Response.StatusCode == 200)
}


func (t *UserTest) TestUserFindDoesntWorkForOtherUser() {
	t.SendAsRegularUser("get", fmt.Sprintf("/users/%s", t.adminUser.Key), nil)
	t.Assert(t.Response.StatusCode == 401)
}


func (t *UserTest) TestUserFindWorksForAdminUsers() {
	t.SendAsSuperAdmin("get", fmt.Sprintf("/users/%s", t.adminUser.Key), nil)
	t.Assert(t.Response.StatusCode == 200)

	t.SendAsAdmin("get", fmt.Sprintf("/users/%s", t.regularUser.Key), nil)
	t.Assert(t.Response.StatusCode == 200)
}

// ----- User create -----

func (t *UserTest) TestUserCreate() {
	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test.com",
		"password": "testpassword7",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsNoOne("post", "/users/create", reqData)

	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	// get key from ResponseBody
	key := body.Key

	t.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", key), nil)
	t.Assert(t.Response.StatusCode == 200)
}


func (t *UserTest) TestInvalidUserCreate() {
	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest2@test.com",
		"password": "testpassword12",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsNoOne("post", "/users/create", reqData)

	t.Assert(t.Response.StatusCode == 400)
}


func (t *UserTest) TestValidAdminUserCreate() {
	data := map[string]interface{}{
		"firstName": "Testy",
		"lastName": "McTest",
		"email": "testymctest@test.com",
		"password": "testpassword7",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsSuperAdmin("post", "/users/create", reqData)

	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	// get key from ResponseBody
	key := body.Key

	t.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", key), nil)
	t.Assert(t.Response.StatusCode == 200)
}


// ----- User update (PATCH) -----

func (t *UserTest) TestUserUpdateWorksForUser() {
	data := map[string]interface{}{
		"firstName": "Testie",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsRegularUser("patch", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Testie")
}

func (t *UserTest) TestUserUpdateFailsForOtherUser() {
	data := map[string]interface{}{
		"firstName": "Testie",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsRegularUser("patch", fmt.Sprintf("/users/%s", t.adminUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 401)
}

func (t *UserTest) TestUserUpdateWorksForSuperAdmin() {
	data := map[string]interface{}{
		"firstName": "Test",
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsSuperAdmin("patch", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Test")
}


func (t *UserTest) TestSuperAdminCanUpdateRole() {
	data := map[string]interface{}{
		"firstName": "Test",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsSuperAdmin("patch", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Test")
	t.Assert(body.Role == 2)
}


func (t *UserTest) TestAdminCanNotUpdateRole() {
	data := map[string]interface{}{
		"firstName": "Test",
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsAdmin("patch", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 401)
}

// ----- User update OG (PUT) -----

func (t *UserTest) TestUserUpdateOGWorksForUser() {
	data := map[string]interface{}{
		"firstName": "Testie",
		"lastName": t.regularUser.LastName,
		"email": t.regularUser.Email,
		"role": t.regularUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsRegularUser("put", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Testie")
	t.Assert(body.LastName == t.regularUser.LastName)
	t.Assert(body.Email == t.regularUser.Email)
	t.Assert(body.Role == t.regularUser.Role)
}

func (t *UserTest) TestUserUpdateOGFailsForOtherUser() {
	data := map[string]interface{}{
		"firstName": "Testie",
		"lastName": t.adminUser.LastName,
		"email": t.adminUser.Email,
		"role": t.adminUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsRegularUser("put", fmt.Sprintf("/users/%s", t.adminUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 401)
}

func (t *UserTest) TestUserUpdateOGWorksForSuperAdmin() {
	data := map[string]interface{}{
		"firstName": "Test",
		"lastName": t.regularUser.LastName,
		"email": t.regularUser.Email,
		"role": t.regularUser.Role,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsSuperAdmin("put", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Test")
	t.Assert(body.LastName == t.regularUser.LastName)
	t.Assert(body.Email == t.regularUser.Email)
	t.Assert(body.Role == t.regularUser.Role)
}


func (t *UserTest) TestSuperAdminCanUpdateRoleWithOG() {
	data := map[string]interface{}{
		"firstName": "Test",
		"lastName": t.regularUser.LastName,
		"email": t.regularUser.Email,
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsSuperAdmin("put", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 200)

	body := dtos.UserDTO{}
	jsonErr = json.Unmarshal(t.ResponseBody, &body)

	t.Assert(body.FirstName == "Test")
	t.Assert(body.Role == 2)
	t.Assert(body.LastName == t.regularUser.LastName)
	t.Assert(body.Email == t.regularUser.Email)
}


func (t *UserTest) TestAdminCanNotUpdateRoleWithOG() {
	data := map[string]interface{}{
		"firstName": "Test",
		"lastName": t.regularUser.LastName,
		"email": t.regularUser.Email,
		"role": 2,
	}

	reqData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr)
	}

	t.SendAsAdmin("put", fmt.Sprintf("/users/%s", t.regularUser.Key), reqData)
	t.Assert(t.Response.StatusCode == 401)
}


// ----- User delete -----

func(t *UserTest) TestDeleteWorksForCurrentUser() {
	t.SendAsRegularUser("delete", fmt.Sprintf("/users/%s", t.regularUser.Key), nil)
	t.Assert(t.Response.StatusCode == 200)
}

func(t *UserTest) TestDeleteDoesNotWorksForOtherUser() {
	t.SendAsRegularUser("delete", fmt.Sprintf("/users/%s", t.adminUser.Key), nil)
	t.Assert(t.Response.StatusCode == 401)
}

func(t *UserTest) TestDeleteWorksForSuperUser() {
	t.SendAsSuperAdmin("delete", fmt.Sprintf("/users/%s", t.adminUser.Key), nil)
	t.Assert(t.Response.StatusCode == 200)
}
