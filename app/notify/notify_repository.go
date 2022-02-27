package notify

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/danwakefield/fnmatch"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	"gopkg.in/yaml.v2"
)

const (
	MESSAGE_INFO = 1 << iota
	MESSAGE_WARNING
	MESSAGE_ERROR
)

type (
	Notification struct {
		filename     string    `yaml:"-"`
		CreationTime time.Time `yaml:"ct"`
		Message      string    `yaml:"msg"`
		Readen       bool      `yaml:"rd"`
		Type         int       `yaml:"tp"`
	}

	NotifyRepository struct {
		DataFolder string
	}
)

func (notification *Notification) String() string {
	return fmt.Sprintf("%s %s", notification.CreationTime.Format(time.ANSIC), notification.Message)
}

func CreateNotifyRepository(datafolder string) *NotifyRepository {
	log.Printf("Creating notification repository - %s", datafolder)
	var err error
	if _, err = os.Stat(datafolder); err != nil {
		if err = pathtools.CreatePath(datafolder); err != nil {
			log.Printf(" ! Failed: %v", err)
			return nil
		}
		log.Printf(" OK")

	}
	repo := &NotifyRepository{
		DataFolder: datafolder,
	}
	return repo
}

func createNotificationFile(datafolder, message string) *Notification {
	cd := time.Now()
	notification := &Notification{
		filename:     path.Join(datafolder, NotificationFileName(cd)),
		CreationTime: cd,
		Message:      message,
		Readen:       false,
	}
	if body, err := yaml.Marshal(notification); err == nil {
		if err := os.WriteFile(notification.filename, body, 0664); err == nil {
			return notification
		}
	}
	return nil
}

func (notification *Notification) Update() (err error) {
	var body []byte
	if body, err = yaml.Marshal(notification); err == nil {
		err = os.WriteFile(notification.filename, body, 0664)
	}
	return err
}

func (notification *Notification) Delete() error {
	return os.Remove(notification.filename)
}

func (repo *NotifyRepository) validateRepository() error {
	if _, err := os.Stat(repo.DataFolder); err == nil {
		return nil
	}
	log.Printf("Creating notification repository - %s", repo.DataFolder)
	if err := pathtools.CreatePath(repo.DataFolder); err != nil {
		log.Printf(" ! Failed: %v", err)
		return err
	}
	log.Printf(" OK")
	return nil
}
func NotificationFileName(creationTime time.Time) string {
	return fmt.Sprintf("%010d.notify", creationTime.Unix())
}

func (notification *Notification) GetFileName(dataFolder string) string {
	return path.Join(dataFolder, NotificationFileName(notification.CreationTime))
}

func (repo *NotifyRepository) AddNotification(message string, messageType int) (*Notification, error) {
	ct := time.Now()
	notification := &Notification{
		CreationTime: ct,
		Message:      message,
		Readen:       false,
		Type:         messageType,
	}
	body, err := yaml.Marshal(notification)
	if err != nil {
		return nil, err
	}

	return notification, os.WriteFile(notification.GetFileName(repo.DataFolder), body, 0664)
}

func (repo *NotifyRepository) DeleteNotification(creationTime time.Time) error {
	if err := repo.validateRepository(); err != nil {
		return err
	}
	notiFileName := path.Join(repo.DataFolder, NotificationFileName(creationTime))
	return os.Remove(notiFileName)
}

func (repo *NotifyRepository) ReadNotification(filename string) (*Notification, error) {
	if err := repo.validateRepository(); err != nil {
		return nil, err
	}
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var notification *Notification
	if err := yaml.Unmarshal(body, &notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (repo *NotifyRepository) GetNotifications() ([]*Notification, error) {
	if err := repo.validateRepository(); err != nil {
		return nil, err
	}
	f, err := os.Open(repo.DataFolder)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	files, err := f.ReadDir(0)
	if err != nil {
		return nil, err
	}
	notifications := make([]*Notification, len(files))
	index := 0
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := path.Join(repo.DataFolder, file.Name())
		if fnmatch.Match("*.notify", path.Base(filename), fnmatch.FNM_FILE_NAME) {
			if notification, err := repo.ReadNotification(filename); err == nil {
				notifications[index] = notification
				index++
			}
		}
	}
	return notifications[0:index], nil
}
