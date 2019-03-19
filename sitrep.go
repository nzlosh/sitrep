package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/jtblin/go-ldap-client"
	"github.com/nzlosh/sitrep/backend"
	"github.com/nzlosh/sitrep/session_manager"
	"gopkg.in/yaml.v2"
	//~ "gopkg.in/gomail.v2"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// include the session manager from github.
const version string = "0.15"

type Range struct {
	inclusive bool     // Inclusive or Exclusive period.
	dow       []int    // 0="all", 1="mon" .. 7="sun"
	month     []int    // 0="all", 1="jan" .. 12="dec"
	woy       []int    // 1..5 Weeks of the year.
	time      []string // "00:00:00" .. "23:59:59" HH:MM:SS
	year      []int    // 0="all", 2018, 2019 etc.
}

type Period struct {
	range_ Range
}
type Periods []Period

type Rule struct {
	id      int
	owner   Person
	name    string
	periods Periods
}

type Person struct {
	id   int
	name string
	/********
	 *
	 * DEFINE REST OF STRUCTURE.
	 *
	 * */
}
type People []Person

type Group struct {
	id      int
	name    string
	members People
}
type Groups []Group

type MySQLConfig struct {
	User         string
	Password     string
	Protocol     string
	Host         string
	Port         int
	DatabaseName string
	Charset      string
	ParseTime    bool
	UseTLS       bool
	Autocommit   bool
}

type LdapConfig struct {
	Attributes         []string
	Base               string
	BindDN             string
	BindPassword       string
	GroupFilter        string
	Host               string
	ServerName         string
	UserFilter         string
	Port               int
	InsecureSkipVerify bool
	UseSSL             bool
	SkipTLS            bool
	ClientCertificates []string
}

type ServerConfig struct {
	Listen string
	Port   int
}

type Config struct {
	Mysql  MySQLConfig
	Ldap   LdapConfig
	Server ServerConfig
}

func main() {

	fmt.Printf("Situation Report Daemon v%s\n", version)
	fmt.Printf("Backend: %s\n", backend_mysql.Name())

	cfg := LoadConfig("sitrep.cfg")

	i := Impl{}
	i.InitDB(cfg.Mysql)
	i.InitSchema()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	// TODO: Implement middleware for user authentication/authorisation
	//~ api.Use(&rest.
	auth_client := &ldap.LDAPClient{
		Base:               cfg.Ldap.Base,
		BindDN:             cfg.Ldap.BindDN,
		BindPassword:       cfg.Ldap.BindPassword,
		GroupFilter:        cfg.Ldap.GroupFilter,
		Host:               cfg.Ldap.Host,
		ServerName:         cfg.Ldap.ServerName, // server_name is used when checking the server certificate.
		UserFilter:         cfg.Ldap.UserFilter,
		Port:               cfg.Ldap.Port,
		InsecureSkipVerify: cfg.Ldap.InsecureSkipVerify,
		UseSSL:             cfg.Ldap.UseSSL,
		SkipTLS:            cfg.Ldap.SkipTLS,
		//~ client_certificates: []
	}
	fmt.Println(AuthenticateUser(auth_client, "carlos", "Wrong"))

	router, err := rest.MakeRouter(
		rest.Get("/alertlog", i.GetAlertLog),
		rest.Get("/alertherolog", i.GetHeroAlertLog),
		rest.Get("/alertcomment", i.GetAlertComment),
		rest.Get("/admins", i.GetAllAdmins),
		//~ rest.Post("/admins", i.PostAdmins),
		rest.Get("/oncallreport", i.GetOncallReport),
		rest.Get("/reportaction", i.GetReportAction),
		rest.Get("/reportimprovement", i.GetReportImprovement),
		rest.Get("/reportseverity", i.GetReportSeverity),
		rest.Get("/version", i.GetVersion),
		rest.Get("/latestevent", i.GetVersion),
		rest.Get("/test_login", session_manager.Login),
		rest.Post("/test_login2", session_manager.Login),
		//~ rest.Put("/auth/$user", AuthenticateUser),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Listen, cfg.Server.Port), api.MakeHandler()))
}

type Impl struct {
	DB *gorm.DB
}

func LoadConfig(cfg_file string) Config {
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
	m := Config{}
	err = yaml.Unmarshal([]byte(str_cfg), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Type assertion required
	return m
}

func (i *Impl) InitDB(mysql MySQLConfig) {
	var err error
	cxn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=%t&tls=%t&autocommit=%t",
		mysql.User,
		mysql.Password,
		mysql.Protocol,
		mysql.Host,
		mysql.Port,
		mysql.DatabaseName,
		mysql.Charset,
		mysql.ParseTime,
		mysql.UseTLS,
		mysql.Autocommit,
	)

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

func LegalConstraints() {
	/*
	 * Source: http://www.ss2ideal.fr/heures-sup-et-astreintes-en-ssiiesn/
	 *
	 * Rest and Recuperation
	 * =====================
	 * 35 consecutive hours at least once per week.
	 * 11 consecutive hours between 2 days of work.
	 *
	 * Remuneration
	 * ============
	 *
	 * Cummulation
	 * ----------
	 * <=8hours per week at 125%
	 * >8hours per week at 150%
	 *
	 * Sunday and Public Holidays
	 * --------------------------
	 * Forbidden to work more than 15 Sundays per year.
	 * Pay at 200%
	 *
	 * Night
	 * -----
	 * 22h to 6h - pay at 150%
	 *
	 */
}

//~ struct User {

//~ }
//~ func (i *Impl)PostUser(w rest.ResponseWriter, r *rest.Request) {
//~ err := r.DecodeJsonPayload()
//~ if err
//~ }

//~ struct Group {

//~ }
//~ func (i *Impl)PostGroup(w rest.ResponseWriter, r *rest.Request) {
//~ err := r.DecodeJsonPayload()
//~ if err
//~ }

/* Reference: https://stackoverflow.com/questions/549/the-definitive-guide-to-form-based-website-authentication */
func AuthenticateUser(auth_client *ldap.LDAPClient, username, password string) bool {
	/* This function should accept an encrypted */
	err := auth_client.Connect()
	if err != nil {
		log.Fatalf("Connection error %v", err)
	}
	defer auth_client.Close()

	success := false
	ok, user, err := auth_client.Authenticate(username, password)
	if err != nil {
		log.Printf("Error authenticating user %s: %+v", username, err)
	}

	if ok {
		log.Printf("Log in succeeded for user '%+v'", user)
		success = true
	} else {
		log.Printf("Authentication failed for user '%s'", username)
	}
	return success
}

func (i *Impl) GetLatestEventId() {
	// This end point is used to poll for the last event id.  It is used by the client
	// to detect when new events have been added to the database.

}

func (i *Impl) InitSchema() {
	//i.DB.AutoMigrate()        // Don't init, the database schema is managed by the DBA.
}

func (i *Impl) GetVersion(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(version)
}

// ============================= ADMINS ======================================
type Admin struct {
	Id       int64  `gorm:"column:admin_id" json:"id"`
	Login    string `gorm:"column:login" sql:"size 50" json:"login"`
	IsActive bool   `gorm:"column:is_active" json:"is_active"`
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
	Id        int64  `gorm:"column:alert_id" json:"id"`
	AlertDate int64  `gorm:"column:alert_date" json:"alert_date"`
	Host      string `gorm:"column:host" sql:"size 50" json:"host"`
	Service   string `gorm:"column:service" sql:"size 100" json:"service"`
	Status    string `gorm:"column:status" sql:"size 50" json:"status"`
	Output    string `gorm:"column:output" sql:"size 500" json:"output"`
}

func (a AlertLog) TableName() string {
	return "alerts_log"
}
func (i *Impl) GetAlertLog(w rest.ResponseWriter, r *rest.Request) {
	alert_log := []AlertLog{}
	i.DB.Find(&alert_log)
	w.WriteJson(&alert_log)
}

func (i *Impl) GetHeroAlertLog(w rest.ResponseWriter, r *rest.Request) {
	alert_log := []AlertLog{}
	i.DB.Where(
		"WEEKDAY(DATE(FROM_UNIXTIME(alert_date))) < 5").Where(
		"TIME(FROM_UNIXTIME(alert_date)) BETWEEN TIME(\"09:00\") AND TIME(\"19:00\")").Find(&alert_log)
	w.WriteJson(&alert_log)
}

// ============================= ONCALL REPORT ======================================
type OncallReport struct {
	Id        int64     `gorm:"column:report_id" json:"id"`
	DateStart time.Time `gorm:"column:date_start" json:"date_start"`
	DateEnd   time.Time `gorm:"column:date_end" json:"date_end"`
	Comment   string    `gorm:"column:comment" sql:"size 500" json:"comment"`
}

func (a OncallReport) TableName() string {
	return "oncall_report"
}
func (i *Impl) GetOncallReport(w rest.ResponseWriter, r *rest.Request) {
	oncall_report := []OncallReport{}
	i.DB.Find(&oncall_report)
	w.WriteJson(&oncall_report)
}

// ============================= REPORT ACTION ======================================
type ReportAction struct {
	Id     int64  `gorm:"column:action_id" json:"id"`
	Action string `gorm:"column:action" sql:"size 30" json:"action"`
}

func (a ReportAction) TableName() string {
	return "report_action"
}
func (i *Impl) GetReportAction(w rest.ResponseWriter, r *rest.Request) {
	report_action := []ReportAction{}
	i.DB.Find(&report_action)
	w.WriteJson(&report_action)
}

// ============================= ALERT COMMENT ======================================
type AlertComment struct {
	Id            int64  `gorm:"column:alert_id" json:"id"`
	ActionId      int64  `gorm:"column:action_id" json:"action_id"`
	ImprovementId int64  `gorm:"column:improvement_id" json:"improvement_id"`
	SeverityId    int64  `gorm:"column:severity_id" json:"severity_id"`
	Note          string `gorm:"column:note" sql:"size 500" json:"note"`
	AdminId       int64  `gorm:"column:admin_id" json:"admin_id"`
	TimeSpent     int64  `gorm:"column:spent" json:"time_spent"`
}

func (a AlertComment) TableName() string {
	return "alerts_comment"
}

func (i *Impl) GetAlertComment(w rest.ResponseWriter, r *rest.Request) {
	report_comment := []AlertComment{}
	i.DB.Find(&report_comment)
	w.WriteJson(&report_comment)
}

// ============================= REPORT IMPROVEMENT ======================================
type ReportImprovement struct {
	Id          int64  `gorm:"column:improvement_id" json:"id"`
	Improvement string `gorm:"column:improvement" sql:"size 30" json:"improvement"`
}

func (a ReportImprovement) TableName() string {
	return "report_improvement"
}
func (i *Impl) GetReportImprovement(w rest.ResponseWriter, r *rest.Request) {
	report_improvement := []ReportImprovement{}
	i.DB.Find(&report_improvement)
	w.WriteJson(&report_improvement)
}

// ============================= REPORT SEVERITY ======================================
type ReportSeverity struct {
	Id       int64  `gorm:"column:severity_id" json:"id"`
	Severity string `gorm:"column:severity" sql:"size 30" json:"severity"`
}

func (a ReportSeverity) TableName() string {
	return "report_severity"
}
func (i *Impl) GetReportSeverity(w rest.ResponseWriter, r *rest.Request) {
	report_severity := []ReportSeverity{}
	i.DB.Find(&report_severity)
	w.WriteJson(&report_severity)
}
