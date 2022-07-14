package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/zohaibAsif/simple_bank_management_system/db/mock"
	db "github.com/zohaibAsif/simple_bank_management_system/db/sqlc"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

type eqCreateUserParamsMatcher struct {
	args     db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {

	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(args.HashedPassword, e.password)
	if err != nil {
		return false
	}

	e.args.HashedPassword = args.HashedPassword
	return reflect.DeepEqual(e.args, args)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.args, e.password)
}

func EqCreateUserParams(args db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{args, password}
}

func TestCreateUserApi(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(*mockdb.MockStore)
		checkResponse func(*httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(args, password)).Times(1).Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"

			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomName(6)

	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomName(6),
		HashedPassword: hashedPassword,
		FullName:       util.RandomName(6),
		Email:          util.RandomEmail(),
	}

	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
