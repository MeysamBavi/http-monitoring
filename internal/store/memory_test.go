package store_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"github.com/MeysamBavi/http-monitoring/internal/store"
	"go.uber.org/zap"
)

func TestUserGetAndAdd(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	if err := s.User().Add(ctx, &model.User{
		Id:       123,
		Username: "meysam",
		Password: "123456",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user, err := s.User().Get(ctx, 123)

	var notFoundError store.NotFoundError
	if err == nil || !errors.As(err, &notFoundError) {
		t.Fatalf("should throw not found: %v %v", user, err)
	}

	user, err = s.User().GetByUsername(ctx, "meysam")

	if user == nil {
		t.Fatal("user not found")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Username != "meysam" || user.Password != "123456" {
		t.Fatalf("wrong user data: %v", *user)
	}
}

func TestUniqueUsername(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	if err := s.User().Add(ctx, &model.User{
		Id:       123,
		Username: "meysam",
		Password: "123456",
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := s.User().Add(ctx, &model.User{
		Id:       1,
		Username: "meysam",
		Password: "",
	})

	var duplicateError store.DuplicateError
	if err == nil || !errors.As(err, &duplicateError) {
		t.Fatalf("should throw duplicate error: %v", err)
	}
}

func TestGetUsernameById(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	randomId := model.ID(123)
	user := &model.User{
		Id:       randomId,
		Username: "meysam",
		Password: "123456",
	}

	if err := s.User().Add(ctx, user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Id == randomId {
		t.Errorf("id was equal to random Id: %d", user.Id)
	}

	{
		user, err := s.User().Get(ctx, user.Id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if user == nil || user.Username != "meysam" || user.Password != "123456" {
			t.Fatalf("unexpected value of user: %v", *user)
		}
	}
}

func TestGetUserUrls(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	{
		urls, err := s.Url().GetByUserId(ctx, 1)
		if len(urls) != 0 || err != nil {
			t.Fatalf("urls should be empty, err must be nil: %v %v", urls, err)
		}
	}

	{
		randomId := model.ID(123)
		url := &model.URL{
			Id:        randomId,
			UserId:    1,
			Url:       "hello",
			Threshold: 5,
			Interval:  time.Minute,
		}

		if err := s.Url().Add(ctx, url); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		urls, err := s.Url().GetByUserId(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if urls == nil || len(urls) != 1 {
			t.Fatalf("urls list length is not one: %v", urls)
		}

		if urls[0].Id == randomId {
			t.Errorf("id was equal to random Id: %d", urls[0].Id)
		}

		if urls[0].UserId != 1 || urls[0].Url != "hello" || urls[0].Threshold != 5 || urls[0].Interval != time.Minute {
			t.Fatalf("unexpected value of url: %v", *urls[0])
		}
	}

}

func TestUpdateStat(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	var urlId model.ID
	{
		url := &model.URL{
			UserId:    1,
			Url:       "hello",
			Threshold: 5,
			Interval:  time.Minute,
		}

		if err := s.Url().Add(ctx, url); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		urlId = url.Id
	}

	// Add
	{
		url, err := s.Url().UpdateStat(
			ctx,
			1,
			urlId,
			model.DayStat{UrlId: urlId, Date: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), SuccessCount: 5, FailureCount: 6},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if url.Id != urlId {
			t.Fatalf("Ids don't match: %d != %d", url.Id, urlId)
		}

		if len(url.DayStats) != 1 {
			t.Fatalf("stats list length is not one: %v", url.DayStats)
		}

		stat := url.DayStats[0]

		if !(stat.UrlId == urlId &&
			stat.Date == time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC) &&
			stat.SuccessCount == 5 &&
			stat.FailureCount == 6) {
			t.Fatalf("unexpected stat value: %v", stat)
		}
	}

	// Update
	{
		url, err := s.Url().UpdateStat(
			ctx,
			1,
			urlId,
			model.DayStat{UrlId: urlId, Date: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC), SuccessCount: 1, FailureCount: 1},
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if url.Id != urlId {
			t.Fatalf("Ids don't match: %d != %d", url.Id, urlId)
		}

		if len(url.DayStats) != 1 {
			t.Fatalf("stats list length is not one: %v", url.DayStats)
		}

		stat := url.DayStats[0]

		if !(stat.UrlId == urlId &&
			stat.SuccessCount == 6 &&
			stat.FailureCount == 7) {
			t.Fatalf("unexpected stat value: %v", stat)
		}
	}
}

func TestAddAndGetAlert(t *testing.T) {
	s := store.NewInMemoryStore(zap.NewNop())
	ctx := context.Background()

	{
		alerts, err := s.Alert().GetByUrlId(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(alerts) != 0 {
			t.Fatalf("unexpected value of alert list: %v", alerts)
		}
	}

	{
		randomId := model.ID(123)
		alert := &model.Alert{
			Id:       randomId,
			UserId:   1,
			UrlId:    1,
			Url:      "hello",
			IssuedAt: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC),
		}

		if err := s.Alert().Add(ctx, alert); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if alert.Id == randomId {
			t.Errorf("id was equal to random Id: %d", alert.Id)
		}
	}

	{
		alerts, err := s.Alert().GetByUrlId(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(alerts) != 1 {
			t.Fatalf("alerts list length is not one: %v", alerts)
		}

		a := alerts[0]

		if !(a.UserId == 1 && a.UrlId == 1 && a.Url == "hello" && a.IssuedAt == time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("unexpected value of alert: %v", *a)
		}
	}
}
