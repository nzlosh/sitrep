package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "log"
    "net/http"
    "time"
)

func main() {

    i := Impl{}
    i.InitDB("")
    //i.InitSchema()   // Don't do this, the database is already in place.

    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    router, err := rest.MakeRouter(
        rest.Get("/reminders", i.GetAllReminders),
        rest.Post("/reminders", i.PostReminder),
        rest.Get("/reminders/:id", i.GetReminder),
        rest.Put("/reminders/:id", i.PutReminder),
        rest.Delete("/reminders/:id", i.DeleteReminder),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}

type Admin struct {
    Id          int64           `json:"id"`
    Login       string          `sql:"size 50" json:"login"`
    IsActive    bool            `json:"is_active"`
}

/* This is the table that needs to be denormalised to data to be returned to client.
type AlertsComment struct {
| alert_id       | int(10) unsigned | NO   | PRI | NULL    |       |
| action_id      | int(10) unsigned | NO   |     | 1       |       |
| improvement_id | int(10) unsigned | NO   |     | 1       |       |
| severity_id    | int(10) unsigned | NO   |     | 1       |       |
| note           | varchar(500)     | YES  |     |         |       |
| admin_id       | int(10) unsigned | YES  |     | NULL    |       |
| spent          | int(10) unsigned | YES  |     | NULL    |       |
}
*/

type AlertsLog struct {
    Id          int64       `json:"id"`
    AlertDate   int64       `json:"alert_date"`
    Host        string      `sql:"size 50" json:"host"`
    Service     string      `sql:"size 100" json:"service"`
    Status      string      `sql:"size 50" json:"status"`
    Output      string      `sql:"size 500" json:"output"`
}

type OncallReport struct {
    Id          int64       `json:"id"`
    DateStart   time.Time   `json:"date_start"`
    DateEnd     time.Time   `json:"date_end"`
    Comment     string      `sql:"size 500" json:"comment"`
}

type ReportAction struct {
    Id          int64       `json:"id"`
    Action      string      `sql:"size 30" json:"action"`
}

type ReportImprovement struct {
    Id          int64       `json:"id"`
    Improvement string      `sql:"size 30" json:"improvement"`
}

type ReportSeverity struct {
    Id          int64       `json:"id"`
    Severity    string      `sql:"size 30" json:"severity"`
}


type Impl struct {
    DB *gorm.DB
}

func (i *Impl) InitDB(cxn string) {
    var err error
    i.DB, err = gorm.Open("mysql", cxn)
    if err != nil {
        log.Fatalf("Got error when connect database, the error is '%v'", err)
    }
    i.DB.LogMode(true)
}


func (i *Impl) InitSchema() {
    i.DB.AutoMigrate(&Reminder{})
}

func (i *Impl) GetAllReminders(w rest.ResponseWriter, r *rest.Request) {
    reminders := []Reminder{}
    i.DB.Find(&reminders)
    w.WriteJson(&reminders)
}

func (i *Impl) GetReminder(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    reminder := Reminder{}
    if i.DB.First(&reminder, id).Error != nil {
        rest.NotFound(w, r)
        return
    }
    w.WriteJson(&reminder)
}

func (i *Impl) PostReminder(w rest.ResponseWriter, r *rest.Request) {
    reminder := Reminder{}
    if err := r.DecodeJsonPayload(&reminder); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := i.DB.Save(&reminder).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteJson(&reminder)
}

func (i *Impl) PutReminder(w rest.ResponseWriter, r *rest.Request) {

    id := r.PathParam("id")
    reminder := Reminder{}
    if i.DB.First(&reminder, id).Error != nil {
        rest.NotFound(w, r)
        return
    }

    updated := Reminder{}
    if err := r.DecodeJsonPayload(&updated); err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    reminder.Message = updated.Message

    if err := i.DB.Save(&reminder).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteJson(&reminder)
}

func (i *Impl) DeleteReminder(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    reminder := Reminder{}
    if i.DB.First(&reminder, id).Error != nil {
        rest.NotFound(w, r)
        return
    }
    if err := i.DB.Delete(&reminder).Error; err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
}
