package app

import (
	_ "github.com/revel/modules"
	"github.com/joho/godotenv"
	"github.com/revel/revel"
	"strconv"
	"os"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func LoadEnv() {
	var err error
	if (revel.RunMode == "Prod") {
		revel.AppLog.Error("prod not supported yet")
		err = godotenv.Load(".env/prod.env")
	} else {
		err = godotenv.Load(".env/dev.env")
	}

	if err != nil {
		revel.AppLog.Error(err.Error())
	} else {
		revel.AppLog.Info("Successfully Loaded Env")
	}
}

var DBConnection DBConn

func InitDB() {
	DBConnection = DBConn{}

	port, err := strconv.Atoi(os.Getenv("REVEL_DBPORT"))
	if err != nil {
		panic(err)
	}

	DBConnection.SetHost(os.Getenv("REVEL_DBHOST"))
	DBConnection.SetUser(os.Getenv("REVEL_DBUSER"))
	DBConnection.SetPassword(os.Getenv("REVEL_DBPASSWORD"))
	DBConnection.SetName(os.Getenv("REVEL_DBNAME"))
	DBConnection.SetPort(port)

	err = DBConnection.FireUp()
	
	if err != nil {
		panic(err)
	}

	revel.OnAppStop(func() {DBConnection.CoolDown()})
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	revel.OnAppStart(LoadEnv, 0)
	revel.OnAppStart(InitDB)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
