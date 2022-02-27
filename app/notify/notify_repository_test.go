package notify

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
)

func setup() *NotifyRepository {
	tmp_folder, err := os.MkdirTemp(".", "notify_repository_*")
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	return CreateNotifyRepository(tmp_folder)
}

func shutdown(repository *NotifyRepository) {
	os.RemoveAll(repository.DataFolder)
	fmt.Println("TEST SHUTDOWN")
}

func TestComplete(t *testing.T) {
	t.Run("Test Complete", func(t *testing.T) {
		repository := setup()
		n1, err := repository.AddNotification("NOTIFICATION 1", MESSAGE_INFO)
		fmt.Printf("%v\n", n1)
		n2, err := repository.AddNotification("NOTIFICATION 2", MESSAGE_WARNING)
		fmt.Printf("%v\n", n2)
		notifications, err := repository.GetNotifications()
		fmt.Printf("Notifications: %v - Error: %v", notifications, err)
		shutdown(repository)
	})

}
func TestNotifyRepository_GetNotifications(t *testing.T) {
	type fields struct {
		DataFolder string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*Notification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &NotifyRepository{
				DataFolder: tt.fields.DataFolder,
			}
			got, err := repo.GetNotifications()
			if (err != nil) != tt.wantErr {
				t.Errorf("NotifyRepository.GetNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotifyRepository.GetNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}
