package models

//TODO - Fix the tests
// func TestFindByID(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	app.App.Load()
// 	crud := CrudTarget{}

// 	targetID := "65eed219e7e71042c5399e3a"

// 	got, err := crud.FindByID(targetID)
// 	if err != nil {
// 		t.Fatalf("FindByID returned an error: %v", err)
// 	}

// 	want := Target{
// 		ID:                got.ID,
// 		CustomerUsername:  got.CustomerUsername,
// 		TargetAddress:     got.TargetAddress,
// 		Flag:              got.Flag,
// 		ScanStatus:        got.ScanStatus,
// 		CreatedAt:         got.CreatedAt,
// 		ScanStartedTime:   got.ScanStartedTime,
// 		ScanCompletedTime: got.ScanCompletedTime,
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("FindByID(%s) = %v, want %v", targetID, got, want)
// 	}
// }

// func TestUpdateScanStatus(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	app.App.Load()
// 	crud := CrudTarget{}
// 	targetID := "65eed219e7e71042c5399e3a"

// 	target, err := crud.FindByID(targetID)

// 	got, updateErr := crud.UpdateScanStatus(&target)

// 	want := true

// 	if got != want {
// 		t.Errorf("UpdateScanStatus(%s) = %v, want %v", targetID, got, want)
// 	}
// }

// func TestMarkAsComplete(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	app.App.Load()
// 	crud := CrudTarget{}
// 	targetID := "65eed219e7e71042c5399e3a"

// 	target, err := crud.FindByID(targetID)

// 	got, updateErr := crud.MarkAsComplete(&target)

// 	want := true

// 	if got != want {
// 		t.Errorf("MarkAsComplete(%s) = %v, want %v", targetID, got, want)
// 	}
// }

// func TestUpdateScanStatusByID(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	app.App.Load()
// 	crud := CrudTarget{}
// 	targetID := "65eed219e7e71042c5399e3a"
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	target, err := crud.FindByID(targetID)

// 	got, updateErr := crud.UpdateScanStatusByID(ctx, target.ID, scanstatus.REPORT_GENERATED)

// 	want := true

// 	if got != want {
// 		t.Errorf("UpdateScanStatus(%s) = %v, want %v", targetID, got, want)
// 	}
// }

// func TestGetAll(t *testing.T) {
// 	err := godotenv.Load("../.env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	app.App.Load()
// 	crud := CrudTarget{}
// 	got, err := crud.GetAll()

// 	id := primitive.NewObjectID()
// 	want := []schemas.TargetResponse{
// 		{
// 			ID:                id,
// 			CustomerUsername:  "test",
// 			TargetAddress:     "192.168.1.1",
// 			Flag:              0,
// 			ScanStatus:        99,
// 			CreatedAt:         time.Now(),
// 			ScanStartedTime:   time.Now(),
// 			ScanCompletedTime: time.Now(),
// 		},
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("GetAll() = %v, want %v", got, want)
// 	}

// }
