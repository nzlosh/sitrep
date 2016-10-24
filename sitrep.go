package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jinzhu/gorm"
    "gopkg.in/yaml.v2"
    "log"
    "net/http"
    "time"
    "os"
    "fmt"
)

const version string = "0.13"

func main() {

    fmt.Println("Situation Report Daemon v", version)

    cfg := LoadConfig("sitrep.cfg")

    i := Impl{}
    i.InitDB(cfg)
    i.InitSchema()

    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)

    router, err := rest.MakeRouter(
        rest.Get("/alertlog", i.GetAlertLog),
        rest.Get("/alertcomment", i.GetAlertComment),
        rest.Get("/admins", i.GetAllAdmins),
        rest.Get("/oncallreport", i.GetOncallReport),
        rest.Get("/reportaction", i.GetReportAction),
        rest.Get("/reportimprovement", i.GetReportImprovement),
        rest.Get("/reportseverity", i.GetReportSeverity),
        rest.Get("/version", i.GetVersion),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}




type Impl struct {
    DB *gorm.DB
}


func LoadConfig(cfg_file string) string {
    file, err := os.Open(cfg_file)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    bs := make([]byte, stat.Size())
    _, err = file.Read(bs)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    str_cfg := string(bs)
    m := make(map[interface{}]interface{})
    err = yaml.Unmarshal([]byte(str_cfg), &m)
    if err != nil {
            log.Fatalf("error: %v", err)
    }

    // Type assertion required
    return m["mysql"].(string)
}

func (i *Impl) InitDB(cxn string) {
    var err error
    i.DB, err = gorm.Open("mysql", cxn)
    if err != nil {
        log.Fatalf("Got error when connect database, the error is '%v'", err)
    }
    i.DB.LogMode(true)
}

func SyncTables() {
/*
 * PSEUDO CODE
 * set current date/time
 * set last alert to the last record in alert_logs table.
 * fetch all records of type SERVICE NOTIFICATION and HOST NOTIFICATION older than last alert and current date/time.
 * update alert_log table
*/
}


func (i *Impl) InitSchema() {
    //i.DB.AutoMigrate()        // Don't init, the database schema is managed by the DBA.
}


func (i *Impl) GetVersion(w rest.ResponseWriter, r *rest.Request) {
    w.WriteJson(version)
}


// ============================= ADMINS ======================================
type Admin struct {
    Id          int64       `gorm:"column:admin_id" json:"id"`
    Login       string      `gorm:"column:login" sql:"size 50" json:"login"`
    IsActive    bool        `gorm:"column:is_active" json:"is_active"`
}
func (a Admin) TableName() string {
    return "admin"
}

func (i *Impl) GetAllAdmins(w rest.ResponseWriter, r *rest.Request) {
    admins := []Admin{}
    i.DB.Find(&admins)
    w.WriteJson(&admins)
}


// ============================= ALERT LOG ======================================
type AlertLog struct {
    Id          int64       `gorm:"column:alert_id" json:"id"`
    AlertDate   int64       `gorm:"column:alert_date" json:"alert_date"`
    Host        string      `gorm:"column:host" sql:"size 50" json:"host"`
    Service     string      `gorm:"column:service" sql:"size 100" json:"service"`
    Status      string      `gorm:"column:status" sql:"size 50" json:"status"`
    Output      string      `gorm:"column:output" sql:"size 500" json:"output"`
}
func (a AlertLog) TableName() string {
    return "alerts_log"
}
func (i *Impl) GetAlertLog (w rest.ResponseWriter, r *rest.Request) {
    alert_log := []AlertLog{}
    i.DB.Find(&alert_log)
    w.WriteJson(&alert_log)
}


// ============================= ONCALL REPORT ======================================
type OncallReport struct {
    Id          int64       `gorm:"column:report_id" json:"id"`
    DateStart   time.Time   `gorm:"column:date_start" json:"date_start"`
    DateEnd     time.Time   `gorm:"column:date_end" json:"date_end"`
    Comment     string      `gorm:"column:comment" sql:"size 500" json:"comment"`
}
func (a OncallReport) TableName() string {
    return "oncall_report"
}
func (i *Impl) GetOncallReport (w rest.ResponseWriter, r *rest.Request) {
    oncall_report := []OncallReport{}
    i.DB.Find(&oncall_report)
    w.WriteJson(&oncall_report)
}


// ============================= REPORT ACTION ======================================
type ReportAction struct {
    Id          int64       `gorm:"column:action_id" json:"id"`
    Action      string      `gorm:"column:action" sql:"size 30" json:"action"`
}
func (a ReportAction) TableName() string {
    return "report_action"
}
func (i *Impl) GetReportAction (w rest.ResponseWriter, r *rest.Request) {
    report_action := []ReportAction{}
    i.DB.Find(&report_action)
    w.WriteJson(&report_action)
}


// ============================= ALERT COMMENT ======================================
type AlertComment struct {
    Id              int64   `gorm:"column:alert_id" json:"id"`
    ActionId        int64   `gorm:"column:action_id" json:"action_id"`
    ImprovementId   int64   `gorm:"column:improvement_id" json:"improvement_id"`
    SeverityId      int64   `gorm:"column:severity_id" json:"severity_id"`
    Note            string  `gorm:"column:note" sql:"size 500" json:"note"`
    AdminId         int64   `gorm:"column:admin_id" json:"admin_id"`
    TimeSpent       int64   `gorm:"column:spent" json:"time_spent"`
}
func (a AlertComment) TableName() string {
    return "alerts_comment"
}

func (i *Impl) GetAlertComment (w rest.ResponseWriter, r *rest.Request) {
    report_comment := []AlertComment{}
    i.DB.Find(&report_comment)
    w.WriteJson(&report_comment)
}


// ============================= REPORT IMPROVEMENT ======================================
type ReportImprovement struct {
    Id          int64       `gorm:"column:improvement_id" json:"id"`
    Improvement string      `gorm:"column:improvement" sql:"size 30" json:"improvement"`
}
func (a ReportImprovement) TableName() string {
    return "report_improvement"
}
func (i *Impl) GetReportImprovement (w rest.ResponseWriter, r *rest.Request) {
    report_improvement := []ReportImprovement{}
    i.DB.Find(&report_improvement)
    w.WriteJson(&report_improvement)
}


// ============================= REPORT SEVERITY ======================================
type ReportSeverity struct {
    Id          int64       `gorm:"column:severity_id" json:"id"`
    Severity    string      `gorm:"column:severity" sql:"size 30" json:"severity"`
}
func (a ReportSeverity) TableName() string {
    return "report_severity"
}
func (i *Impl) GetReportSeverity (w rest.ResponseWriter, r *rest.Request) {
    report_severity := []ReportSeverity{}
    i.DB.Find(&report_severity)
    w.WriteJson(&report_severity)
}

