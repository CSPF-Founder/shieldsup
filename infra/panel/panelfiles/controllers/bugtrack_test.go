package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"https://github.com/CSPF-Founder/shieldsup/onpremise/panel/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PuerkitoBio/goquery"
)

// Tests for the BugTrack controller

func (ctx *testContext) mockGetOverviewListByUser(bug models.BugTrack) {
	bugRow := sqlmock.NewRows([]string{"id", "target", "alert_title", "found_at", "status"}).
		AddRow(bug.ID, bug.Target, bug.AlertTitle, bug.FoundDate, bug.Status)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `bugtrack_entries` WHERE user_id = ? ORDER BY `bugtrack_entries`.`id` LIMIT 1")).
		WithArgs(bug.UserID).
		WillReturnRows(bugRow)
}

func (ctx *testContext) mockEmptyGetOverviewListByUser(bug models.BugTrack) {
	bugTrackRow := sqlmock.NewRows([]string{"id", "target", "alert_title", "found_at", "status"})

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `bugtrack_entries` WHERE user_id = ? ORDER BY `bugtrack_entries`.`id` LIMIT 1")).
		WithArgs(bug.ID).
		WillReturnRows(bugTrackRow)
}

func TestDeleteBugTrackWithMissingCSRF(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	testBugTrack := models.BugTrack{
		ID:     1,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)

	deleteURL := fmt.Sprintf("%s/bug-track/delete?id=%d", ctx.server.URL, testBugTrack.ID)
	req, err := http.NewRequest("POST", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login bugtrack: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteInvalidBugTrackID(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	bugJob := models.BugTrack{
		ID:     1,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockEmptyGetOverviewListByUser(bugJob)

	deleteURL := fmt.Sprintf("%s/bug-track/%d", ctx.server.URL, bugJob.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login bugtrack: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteValidBugTrackID(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	testBug := models.BugTrack{
		ID:     1,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockGetOverviewListByUser(testBug)

	deleteURL := fmt.Sprintf("%s/bug-track/delete?id=%d", ctx.server.URL, testBug.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login bugtrack: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteBugTrackWithInvalidSession(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	// ctx, resp := loggedSessionForTest(t, testUser)
	ctx := setupTest(t)

	testBugTrack := models.BugTrack{
		ID:     1,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockGetOverviewListByUser(testBugTrack)

	deleteURL := fmt.Sprintf("%s/bug-track/delete?id=%d", ctx.server.URL, testBugTrack.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login bugtrack: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}
