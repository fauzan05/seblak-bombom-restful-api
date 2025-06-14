package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"seblak-bombom-restful-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddApplicationByKey(t *testing.T) {
	ClearAll()

	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")
	partHeader.Set("X-Admin-Key", "rahasia-123#")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications-use-admin-key", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now", responseBody.Data.AppName)
	assert.Equal(t, "07:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "20:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi", responseBody.Data.Description)
	assert.Equal(t, "08133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestAddApplicationSetting(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now", responseBody.Data.AppName)
	assert.Equal(t, "07:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "20:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi", responseBody.Data.Description)
	assert.Equal(t, "08133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestAddApplicationSettingBadRequest(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, responseBody.Error, "invalid request body :")
}

func TestUpdateApplicationSetting(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now", responseBody.Data.AppName)
	assert.Equal(t, "07:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "20:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi", responseBody.Data.Description)
	assert.Equal(t, "08133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// Simulasi multipart body
	var c bytes.Buffer
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now update")
	_ = writer.WriteField("opening_hours", "08:00:00")
	_ = writer.WriteField("closing_hours", "19:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat update")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 update")
	_ = writer.WriteField("description", "Ini adalah deskripsi update")
	_ = writer.WriteField("phone_number", "+628133546789")
	_ = writer.WriteField("email", "seblak-update@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram-update")
	_ = writer.WriteField("instagram_link", "https://www.instagram-update.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter-update")
	_ = writer.WriteField("twitter_link", "https://www.twitter-update.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook-update")
	_ = writer.WriteField("facebook_link", "https://www.facebook-update.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err = GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader = textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err = writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request = httptest.NewRequest(http.MethodPost, "/api/applications", &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytesUpdate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now update", responseBody.Data.AppName)
	assert.Equal(t, "08:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "19:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat update", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 update", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi update", responseBody.Data.Description)
	assert.Equal(t, "+628133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak-update@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram-update", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram-update.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter-update", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter-update.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook-update", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook-update.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestAddApplicationSettingFileLogoExceededThan1Mb(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")
	_ = writer.WriteField("service_fee", "1000")

	filename, content, err := GenerateDummyJPEG(2 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "file size is too large, maxium is 1MB", responseBody.Error)
}

func TestGetApplicationSetting(t *testing.T) {
	TestAddApplicationSetting(t)

	request := httptest.NewRequest(http.MethodGet, "/api/applications", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now", responseBody.Data.AppName)
	assert.Equal(t, "07:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "20:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi", responseBody.Data.Description)
	assert.Equal(t, "08133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}
